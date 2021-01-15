// Package model implements convenience methods for
// managing indexes on top of the Store.
// See this doc for the general idea https://github.com/m3o/dev/blob/feature/storeindex/design/auto-indexes.md
// Prior art/Inspirations from github.com/gocassa/gocassa, which
// is a similar package on top an other KV store (Cassandra/gocql)
package model

import (
	"encoding/base32"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"reflect"
	"strings"
	"unicode/utf8"

	"github.com/micro/micro/v3/service/store"
)

var (
	ErrorNotFound             = errors.New("not found")
	ErrorMultipleRecordsFound = errors.New("multiple records found")
	ErrorStoreIsNil           = errors.New("the store is nil, something is not initialized properly")
)

type OrderType string

const (
	OrderTypeUnordered = OrderType("unordered")
	OrderTypeAsc       = OrderType("ascending")
	OrderTypeDesc      = OrderType("descending")
)

const (
	queryTypeEq = "eq"
	indexTypeEq = "eq"
)

func defaultIndex() Index {
	idIndex := ByEquality("ID")
	idIndex.Order.Type = OrderTypeUnordered
	return idIndex
}

type model struct {
	store store.Store
	// helps logically separate keys in a model where
	// multiple `Model`s share the same underlying
	// physical database.
	namespace string
	indexes   []Index
	options   ModelOptions
	instance  interface{}
}

// Model represents a place where data can be saved to and
// queried from.
type Model interface {
	// Create a new object. (Maintains indexes set up)
	Create(v interface{}) error
	// Update will take an existing object and update it.
	// TODO: Make use of "sync" interface to lock, read, write, unlock
	Update(v interface{}) error
	// Read accepts a pointer to a value and expects to fine one or more
	// elements. Read throws an error if a value is not found or we can't
	// find a matching index for a slice based query.
	Read(query Query, resultPointer interface{}) error
	// Deletes a record. Delete only support Equals("id", value) for now.
	// @todo Delete only supports string keys for now.
	Delete(query Query) error
}

type ModelOptions struct {
	Debug     bool
	IdIndex   Index
	Namespace string
}

func New(store store.Store, instance interface{}, indexes []Index, options *ModelOptions) Model {
	debug := false
	var idIndex Index
	namespace := reflect.TypeOf(instance).String()
	if options != nil {
		debug = options.Debug
		if options.IdIndex.Type != "" {
			idIndex = options.IdIndex
		}

		if len(options.Namespace) > 0 {
			namespace = options.Namespace
		}
	}
	if idIndex.Type == "" {
		idIndex = defaultIndex()
	}
	return &model{
		store, namespace, indexes, ModelOptions{
			Debug:   debug,
			IdIndex: idIndex,
		}, instance}
}

type Index struct {
	FieldName string
	// Type of index, eg. equality
	Type  string
	Order Order
	// Do not allow duplicate values of this field in the index.
	// Useful for emails, usernames, post slugs etc.
	Unique bool
	// Strings for ordering will be padded to a fix length
	// Not a useful property for Querying, please ignore this at query time.
	// Number is in bytes, not string characters. Choose a sufficiently big one.
	// Consider that each character might take 4 bytes given the
	// internals of reverse ordering. So a good rule of thumbs is expected
	// characters in a string X 4
	StringOrderPadLength int
	// True = base32 encode ordered strings for easier management
	// or false = keep 4 bytes long runes that might dispaly weirdly
	Base32Encode bool

	FloatFormat string
	Float64Max  float64
	Float32Max  float32
}

type Order struct {
	FieldName string
	// Ordered or unordered keys. Ordered keys are padded.
	// Default is true. This option only exists for strings, where ordering
	// comes at the cost of having rather long padded keys.
	Type OrderType
}

func (i Index) ToQuery(value interface{}) Query {
	return Query{
		Index: i,
		Value: value,
		Order: i.Order,
	}
}

func Indexes(indexes ...Index) []Index {
	return indexes
}

// ByEquality constructs an equiality index on `fieldName`
func ByEquality(fieldName string) Index {
	return Index{
		FieldName: fieldName,
		Type:      indexTypeEq,
		Order: Order{
			Type:      OrderTypeAsc,
			FieldName: fieldName,
		},
		StringOrderPadLength: 16,
		Base32Encode:         false,
		FloatFormat:          "%019.5f",
		Float64Max:           92233720368547,
		Float32Max:           922337,
	}
}

type Query struct {
	Index
	Order  Order
	Value  interface{}
	Offset int64
	Limit  int64
}

// Equals is an equality query by `fieldName`
// It filters records where `fieldName` equals to a value.
func Equals(fieldName string, value interface{}) Query {
	return Query{
		Index: Index{
			Type:      queryTypeEq,
			FieldName: fieldName,
			Order: Order{
				FieldName: fieldName,
				Type:      OrderTypeAsc,
			},
		},
		Value: value,
		Order: Order{
			FieldName: fieldName,
			Type:      OrderTypeAsc,
		},
	}
}

func (d *model) Create(instance interface{}) error {
	// @todo replace this hack with reflection
	js, err := json.Marshal(instance)
	if err != nil {
		return err
	}

	// get the old entries so we can compare values
	// @todo consider some kind of locking (even if it's not distributed) by key here
	// to avoid 2 read-writes happening at the same time
	idQuery := d.options.IdIndex.ToQuery(getFieldValue(instance, d.options.IdIndex.FieldName))

	oldEntry := reflect.New(reflect.ValueOf(instance).Type()).Interface()

	err = d.Read(idQuery, &oldEntry)
	if err != nil && err != ErrorNotFound {
		return err
	}

	// Do uniqueness checks before saving any data
	for _, index := range d.indexes {
		if !index.Unique {
			continue
		}
		potentialClash := reflect.New(reflect.ValueOf(instance).Type()).Interface()
		err = d.Read(index.ToQuery(getFieldValue(instance, index.FieldName)), &potentialClash)
		if err != nil && err != ErrorNotFound {
			return err
		}

		if err == nil {
			return errors.New("Unique index violation")
		}
	}

	id := getFieldValue(instance, d.options.IdIndex.FieldName)
	for _, index := range append(d.indexes, d.options.IdIndex) {
		// delete non id index keys to prevent stale index values
		// ie.
		//
		//  # prefix  slug     id
		//  postByTag/hi-there/1
		//  # if slug gets changed to "hello-there" we will have two records
		//  # without removing the old stale index:
		//  postByTag/hi-there/1
		//  postByTag/hello-there/1`
		//
		// @todo this check will only work for POD types, ie no slices or maps
		// but it's not an issue as right now indexes are only supported on POD
		// types anyway
		if !indexesMatch(d.options.IdIndex, index) &&
			oldEntry != nil &&
			!reflect.DeepEqual(getFieldValue(oldEntry, index.FieldName), getFieldValue(instance, index.FieldName)) {
			k := d.indexToKey(index, id, oldEntry, true)
			if d.store == nil {
				return ErrorStoreIsNil
			}
			err = d.store.Delete(k)
			if err != nil {
				return err
			}
		}
		k := d.indexToKey(index, id, instance, true)
		if d.options.Debug {
			fmt.Printf("Saving key '%v', value: '%v'\n", k, string(js))
		}
		if d.store == nil {
			return ErrorStoreIsNil
		}
		err = d.store.Write(&store.Record{
			Key:   k,
			Value: js,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func getFieldValue(struc interface{}, field string) interface{} {
	r := reflect.ValueOf(struc)
	f := reflect.Indirect(r).FieldByName(strings.Title(field))

	if !f.IsValid() {
		return reflect.Zero(f.Type())
	}
	return f.Interface()
}

// TODO: implement the full functionality. Currently offloads to create.
func (d *model) Update(v interface{}) error {
	return d.Create(v)
}

func setFieldValue(struc interface{}, field string, value interface{}) {
	r := reflect.ValueOf(struc)

	f := reflect.Indirect(r).FieldByName(strings.Title(field))
	f.Set(reflect.ValueOf(value))
}

func (d *model) Read(query Query, resultPointer interface{}) error {
	t := reflect.TypeOf(resultPointer)

	// check if it's a pointer
	if v := t.Kind(); v != reflect.Ptr {
		return fmt.Errorf("Require pointer type. Got %v", v)
	}

	// retrieve the non pointer type
	t = t.Elem()

	// if its a slice then use the list query method
	if t.Kind() == reflect.Slice {
		return d.List(query, resultPointer)
	}

	// otherwise continue on as normal
	for _, index := range append(d.indexes, d.options.IdIndex) {
		if indexMatchesQuery(index, query) {
			k := d.queryToListKey(index, query)
			if d.options.Debug {
				fmt.Printf("Listing key '%v'\n", k)
			}
			if d.store == nil {
				return ErrorStoreIsNil
			}
			recs, err := d.store.Read(k, store.ReadPrefix())
			if err != nil {
				return err
			}
			if len(recs) == 0 {
				return ErrorNotFound
			}
			if len(recs) > 1 {
				return ErrorMultipleRecordsFound
			}
			if d.options.Debug {
				fmt.Printf("Found value '%v'\n", string(recs[0].Value))
			}
			return json.Unmarshal(recs[0].Value, resultPointer)
		}
	}
	return fmt.Errorf("For query type '%v', field '%v' does not match any indexes", query.Type, query.FieldName)
}

func (d *model) List(query Query, resultSlicePointer interface{}) error {
	for _, index := range append(d.indexes, d.options.IdIndex) {
		if indexMatchesQuery(index, query) {
			k := d.queryToListKey(index, query)
			if d.options.Debug {
				fmt.Printf("Listing key '%v'\n", k)
			}
			if d.store == nil {
				return ErrorStoreIsNil
			}
			recs, err := d.store.Read(k, store.ReadPrefix())
			if err != nil {
				return err
			}
			// @todo speed this up with an actual buffer
			jsBuffer := []byte("[")
			for i, rec := range recs {
				jsBuffer = append(jsBuffer, rec.Value...)
				if i < len(recs)-1 {
					jsBuffer = append(jsBuffer, []byte(",")...)
				}
			}
			jsBuffer = append(jsBuffer, []byte("]")...)
			if d.options.Debug {
				fmt.Printf("Found values '%v'\n", string(jsBuffer))
			}
			return json.Unmarshal(jsBuffer, resultSlicePointer)
		}
	}
	return fmt.Errorf("For query type '%v', field '%v' does not match any indexes", query.Type, query.FieldName)
}

func indexMatchesQuery(i Index, q Query) bool {
	if i.FieldName == q.FieldName &&
		i.Type == q.Type &&
		i.Order.Type == q.Order.Type {
		return true
	}
	return false
}

func indexesMatch(i, j Index) bool {
	if i.FieldName == j.FieldName &&
		i.Type == j.Type &&
		i.Order.Type == j.Order.Type {
		return true
	}
	return false
}

func (d *model) queryToListKey(i Index, q Query) string {
	if q.Value == nil {
		return fmt.Sprintf("%v:%v", d.namespace, indexPrefix(i))
	}
	if i.FieldName != i.Order.FieldName && i.Order.FieldName != "" {
		return fmt.Sprintf("%v:%v:%v", d.namespace, indexPrefix(i), q.Value)
	}

	val := reflect.New(reflect.ValueOf(d.instance).Type()).Interface()
	if q.Value != nil {
		setFieldValue(val, i.FieldName, q.Value)
	}
	return d.indexToKey(i, "", val, false)
}

// appendID true should be used when saving, false when querying
// appendID false should also be used for 'id' indexes since they already have the unique
// id. The reason id gets appended is make duplicated index keys unique.
// ie.
// # index # age # id
// users/30/1
// users/30/2
// without ids we could only have one 30 year old user in the index
func (d *model) indexToKey(i Index, id interface{}, entry interface{}, appendID bool) string {
	format := "%v:%v"
	values := []interface{}{d.namespace, indexPrefix(i)}
	filterFieldValue := getFieldValue(entry, i.FieldName)
	orderFieldValue := getFieldValue(entry, i.FieldName)
	orderFieldKey := i.FieldName

	if i.FieldName != i.Order.FieldName && i.Order.FieldName != "" {
		orderFieldValue = getFieldValue(entry, i.Order.FieldName)
		orderFieldKey = i.Order.FieldName
	}

	switch i.Type {
	case indexTypeEq:
		// If the filtering field is different than the ordering field,
		// append the filter key to the key.
		if i.FieldName != i.Order.FieldName && i.Order.FieldName != "" {
			format += ":%v"
			values = append(values, filterFieldValue)
		}
	}

	// Handle the ordering part of the key.
	// The filter and the ordering field might be the same
	typ := reflect.TypeOf(orderFieldValue)
	typName := "nil"
	if typ != nil {
		typName = typ.String()
	}
	format += ":%v"

	switch v := orderFieldValue.(type) {
	case string:
		if i.Order.Type != OrderTypeUnordered {
			values = append(values, d.getOrderedStringFieldKey(i, v))
			break
		}
		values = append(values, v)
	case int64:
		// int64 gets padded to 19 characters as the maximum value of an int64
		// is 9223372036854775807
		// @todo handle negative numbers
		if i.Order.Type == OrderTypeDesc {
			values = append(values, fmt.Sprintf("%019d", math.MaxInt64-v))
			break
		}
		values = append(values, fmt.Sprintf("%019d", v))
	case float32:
		// @todo fix display and padding of floats
		if i.Order.Type == OrderTypeDesc {
			values = append(values, fmt.Sprintf(i.FloatFormat, i.Float32Max-v))
			break
		}
		values = append(values, fmt.Sprintf(i.FloatFormat, v))
	case float64:
		// @todo fix display and padding of floats
		if i.Order.Type == OrderTypeDesc {
			values = append(values, fmt.Sprintf(i.FloatFormat, i.Float64Max-v))
			break
		}
		values = append(values, fmt.Sprintf(i.FloatFormat, v))
	case int:
		// int gets padded to the same length as int64 to gain
		// resiliency in case of model type changes.
		// This could be removed once migrations are implemented
		// so savings in space for a type reflect in savings in space in the index too.
		if i.Order.Type == OrderTypeDesc {
			values = append(values, fmt.Sprintf("%019d", math.MaxInt32-v))
			break
		}
		values = append(values, fmt.Sprintf("%019d", v))
	case int32:
		// int gets padded to the same length as int64 to gain
		// resiliency in case of model type changes.
		// This could be removed once migrations are implemented
		// so savings in space for a type reflect in savings in space in the index too.
		if i.Order.Type == OrderTypeDesc {
			values = append(values, fmt.Sprintf("%019d", math.MaxInt32-v))
			break
		}
		values = append(values, fmt.Sprintf("%019d", v))
	case bool:
		if i.Order.Type == OrderTypeDesc {
			v = !v
		}
		values = append(values, v)
	default:
		panic("bug in code, unhandled type: " + typName + " for field " + orderFieldKey)
	}

	if appendID {
		format += ":%v"
		values = append(values, id)
	}
	return fmt.Sprintf(format, values...)
}

// indexPrefix returns the first part of the keys, the namespace + index name
func indexPrefix(i Index) string {
	var ordering string
	switch i.Order.Type {
	case OrderTypeUnordered:
		ordering = "Unord"
	case OrderTypeAsc:
		ordering = "Asc"
	case OrderTypeDesc:
		ordering = "Desc"
	}
	typ := i.Type
	orderingField := i.Order.FieldName
	if len(orderingField) == 0 {
		orderingField = i.FieldName
	}
	filterField := i.FieldName
	return fmt.Sprintf("%vBy%v%vBy%v", typ, strings.Title(filterField), ordering, strings.Title(orderingField))
}

// pad, reverse and optionally base32 encode string keys
func (d *model) getOrderedStringFieldKey(i Index, fieldValue string) string {
	runes := []rune{}
	if i.Order.Type == OrderTypeDesc {
		for _, char := range fieldValue {
			runes = append(runes, utf8.MaxRune-char)
		}
	} else {
		for _, char := range fieldValue {
			runes = append(runes, char)
		}
	}

	// padding the string to a fixed length
	if len(runes) < i.StringOrderPadLength {
		pad := []rune{}
		for j := 0; j < i.StringOrderPadLength-len(runes); j++ {
			if i.Order.Type == OrderTypeDesc {
				pad = append(pad, utf8.MaxRune)
			} else {
				// space is the first non control operator char in ASCII
				// consequently in Utf8 too so we use it as the minimal character here
				// https://en.wikipedia.org/wiki/ASCII
				//
				// Displays somewhat unfortunately
				// @todo think about a better min rune value to use here.
				pad = append(pad, rune(32))
			}
		}
		runes = append(runes, pad...)
	}

	var keyPart string
	bs := []byte(string(runes))
	if i.Order.Type == OrderTypeDesc {
		if i.Base32Encode {
			// base32 hex should be order preserving
			// https://stackoverflow.com/questions/53301280/does-base64-encoding-preserve-alphabetical-ordering
			dst := make([]byte, base32.HexEncoding.EncodedLen(len(bs)))
			base32.HexEncoding.Encode(dst, bs)
			// The `=` must be replaced with a lower value than the
			// normal alphabet of the encoding since we want reverse order.
			keyPart = strings.ReplaceAll(string(dst), "=", "0")
		} else {
			keyPart = string(bs)
		}
	} else {
		keyPart = string(bs)

	}
	return keyPart
}

func (d *model) Delete(query Query) error {
	if !indexMatchesQuery(d.options.IdIndex, query) {
		return errors.New("Delete query does not match default index")
	}
	oldEntry := reflect.New(reflect.ValueOf(d.instance).Type()).Interface()
	err := d.Read(query, &oldEntry)
	if err != nil {
		return err
	}

	// first delete maintained indexes then id index
	// if we delete id index first then the entry wont
	// be deletable by id again but the maintained indexes
	// will be stuck in limbo
	for _, index := range append(d.indexes, d.options.IdIndex) {
		key := d.indexToKey(index, getFieldValue(oldEntry, d.options.IdIndex.FieldName), oldEntry, true)
		if d.options.Debug {
			fmt.Printf("Deleting key '%v'\n", key)
		}
		if d.store == nil {
			return ErrorStoreIsNil
		}
		err = d.store.Delete(key)
		if err != nil {
			return err
		}
	}
	return nil
}
