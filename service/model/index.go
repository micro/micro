package model

import (
	"fmt"
	"strings"
)

// Index represents a data model index for fast access
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

// Order is the order of the index
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

func indexMatchesQuery(i Index, q Query) bool {
	if strings.ToLower(i.FieldName) == strings.ToLower(q.FieldName) &&
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
	// hack for all listing where we use the eq ID index
	// without a value to list all
	if i.Type == indexTypeAll {
		typ = indexTypeEq
	}
	orderingField := i.Order.FieldName
	if len(orderingField) == 0 {
		orderingField = i.FieldName
	}
	filterField := i.FieldName
	return fmt.Sprintf("%vBy%v%vBy%v", typ, strings.Title(filterField), ordering, strings.Title(orderingField))
}

func newIndex(v string) Index {
	idIndex := ByEquality(v)
	idIndex.Order.Type = OrderTypeUnordered
	return idIndex
}
