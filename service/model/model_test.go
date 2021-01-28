package model

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/gofrs/uuid"
	fs "github.com/micro/micro/v3/service/store/file"
)

type User struct {
	ID      string `json:"id"`
	Age     int    `json:"age"`
	HasPet  bool   `json:"hasPet"`
	Created int64  `json:"created"`
	Tag     string `json:"tag"`
	Updated int64  `json:"updated"`
}

func TestQueryEqualsByID(t *testing.T) {
	table := New(User{}, &Options{
		Store:     fs.NewStore(),
		Namespace: uuid.Must(uuid.NewV4()).String(),
	})

	err := table.Create(User{
		ID:  "1",
		Age: 12,
	})
	if err != nil {
		t.Fatal(err)
	}
	err = table.Create(User{
		ID:  "2",
		Age: 25,
	})
	if err != nil {
		t.Fatal(err)
	}
	users := []User{}
	q := QueryEquals("ID", "1")
	q.Order.Type = OrderTypeUnordered
	err = table.Read(q, &users)
	if err != nil {
		t.Fatal(err)
	}
	if len(users) != 1 {
		t.Fatal(users)
	}
}

type User1 struct {
	Id      string `json:"id"`
	Age     int    `json:"age"`
	HasPet  bool   `json:"hasPet"`
	Created int64  `json:"created"`
	Tag     string `json:"tag"`
	Updated int64  `json:"updated"`
}

func TestQueryEqualsLowerCaseID(t *testing.T) {
	table := New(&User1{}, &Options{
		Store:     fs.NewStore(),
		Namespace: uuid.Must(uuid.NewV4()).String(),
	})

	// pointer insert
	err := table.Create(&User1{
		Id:  "1",
		Age: 12,
	})
	if err != nil {
		t.Fatal(err)
	}
	err = table.Create(User1{
		Id:  "2",
		Age: 25,
	})
	if err != nil {
		t.Fatal(err)
	}
	users := []User1{}
	q := QueryEquals("Id", "1")
	q.Order.Type = OrderTypeUnordered
	err = table.Read(q, &users)
	if err != nil {
		t.Fatal(err)
	}
	if len(users) != 1 {
		t.Fatal(users)
	}
}

func TestQueryEqualsMismatchIDCapitalization(t *testing.T) {
	table := New(&User1{}, &Options{
		Store:     fs.NewStore(),
		Namespace: uuid.Must(uuid.NewV4()).String(),
	})

	// pointer insert
	err := table.Create(&User1{
		Id:  "1",
		Age: 12,
	})
	if err != nil {
		t.Fatal(err)
	}
	err = table.Create(User1{
		Id:  "2",
		Age: 25,
	})
	if err != nil {
		t.Fatal(err)
	}
	users := []User1{}
	q := QueryEquals("id", "1")
	q.Order.Type = OrderTypeUnordered
	err = table.Read(q, &users)
	if err != nil {
		t.Fatal(err)
	}
	if len(users) != 1 {
		t.Fatal(users)
	}
}

func TestQueryEqualsByIDMap(t *testing.T) {
	m := map[string]interface{}{
		"ID":      "id",
		"age":     1,
		"hasPet":  true,
		"created": 1,
		"tag":     "tag",
		"updated": 1,
	}
	table := New(m, &Options{
		Store:     fs.NewStore(),
		Namespace: uuid.Must(uuid.NewV4()).String(),
	})

	err := table.Create(map[string]interface{}{
		"ID":  "1",
		"Age": 12,
	})
	if err != nil {
		t.Fatal(err)
	}
	err = table.Create(map[string]interface{}{
		"ID":  "2",
		"Age": 25,
	})
	if err != nil {
		t.Fatal(err)
	}
	users := []map[string]interface{}{}
	q := QueryEquals("ID", "1")
	q.Order.Type = OrderTypeUnordered
	err = table.Read(q, &users)
	if err != nil {
		t.Fatal(err)
	}
	if len(users) != 1 {
		t.Fatal(users)
	}
}

func TestQueryEqualsByIDMapNoSchemaWithIndexes(t *testing.T) {
	m := map[string]interface{}{}
	table := New(m, &Options{
		Store:     fs.NewStore(),
		Namespace: uuid.Must(uuid.NewV4()).String(),
		Indexes:   []Index{ByEquality("Age")},
	})

	err := table.Create(map[string]interface{}{
		"ID":  "1",
		"Age": 12,
	})
	if err != nil {
		t.Fatal(err)
	}
	err = table.Create(map[string]interface{}{
		"ID":  "2",
		"Age": 25,
	})
	if err != nil {
		t.Fatal(err)
	}
	users := []map[string]interface{}{}
	q := QueryEquals("ID", "1")
	q.Order.Type = OrderTypeUnordered
	err = table.Read(q, &users)
	if err != nil {
		t.Fatal(err)
	}
	if len(users) != 1 {
		t.Fatal(users)
	}
}

func TestListAllMap(t *testing.T) {
	m := map[string]interface{}{
		"ID":      "id",
		"age":     1,
		"hasPet":  true,
		"created": 1,
		"tag":     "tag",
		"updated": 1,
	}
	table := New(m, &Options{
		Store:     fs.NewStore(),
		Namespace: uuid.Must(uuid.NewV4()).String(),
	})

	err := table.Create(map[string]interface{}{
		"ID":  "1",
		"Age": 12,
	})
	if err != nil {
		t.Fatal(err)
	}
	err = table.Create(map[string]interface{}{
		"ID":  "2",
		"Age": 25,
	})
	if err != nil {
		t.Fatal(err)
	}
	users := []map[string]interface{}{}
	q := QueryAll()
	err = table.Read(q, &users)
	if err != nil {
		t.Fatal(err)
	}
	if len(users) != 2 {
		t.Fatal(users)
	}
}

// TestNewModel tests the creation using NewModel and Register
func TestNewModel(t *testing.T) {
	// create a new model
	table := NewModel(
		WithStore(fs.NewStore()),
		WithNamespace(uuid.Must(uuid.NewV4()).String()),
	)

	// register the instance
	table.Register(User{})

	err := table.Create(User{
		ID:  "1",
		Age: 12,
	})
	if err != nil {
		t.Fatal(err)
	}
	err = table.Create(User{
		ID:  "2",
		Age: 25,
	})
	if err != nil {
		t.Fatal(err)
	}
	users := []User{}
	q := QueryEquals("ID", "1")
	q.Order.Type = OrderTypeUnordered
	err = table.Read(q, &users)
	if err != nil {
		t.Fatal(err)
	}
	if len(users) != 1 {
		t.Fatal(users)
	}
}

func TestRead(t *testing.T) {
	table := New(User{}, &Options{
		Store:     fs.NewStore(),
		Indexes:   []Index{ByEquality("age")},
		Namespace: uuid.Must(uuid.NewV4()).String(),
	})
	user := User{}
	// intentionally querying Age to test case tolerance
	err := table.Read(QueryEquals("Age", 25), &user)
	if err != ErrorNotFound {
		t.Fatal(err)
	}

	// test pointer create
	err = table.Create(&User{
		ID:  "1",
		Age: 25,
	})
	if err != nil {
		t.Fatal(err)
	}

	err = table.Read(QueryEquals("age", 25), &user)
	if err != nil {
		t.Fatal(err)
	}
	if user.ID != "1" {
		t.Fatal(user)
	}

	err = table.Create(User{
		ID:  "2",
		Age: 25,
	})
	if err != nil {
		t.Fatal(err)
	}

	err = table.Read(QueryEquals("age", 25), &user)
	if err != ErrorMultipleRecordsFound {
		t.Fatal(err)
	}
}

func TestQueryEquals(t *testing.T) {
	table := New(User{}, &Options{
		Store:     fs.NewStore(),
		Indexes:   []Index{ByEquality("age")},
		Namespace: uuid.Must(uuid.NewV4()).String(),
	})

	err := table.Create(User{
		ID:  "1",
		Age: 12,
	})
	if err != nil {
		t.Fatal(err)
	}
	err = table.Create(User{
		ID:  "2",
		Age: 25,
	})
	if err != nil {
		t.Fatal(err)
	}
	err = table.Create(User{
		ID:  "3",
		Age: 12,
	})
	if err != nil {
		t.Fatal(err)
	}
	users := []User{}
	err = table.Read(QueryEquals("age", 12), &users)
	if err != nil {
		t.Fatal(err)
	}
	if len(users) != 2 {
		t.Fatal(users)
	}
}

func reverse(ss []string) {
	last := len(ss) - 1
	for i := 0; i < len(ss)/2; i++ {
		ss[i], ss[last-i] = ss[last-i], ss[i]
	}
}

func TestOrderingStrings(t *testing.T) {
	type caze struct {
		tags    []string
		reverse bool
	}
	cazes := []caze{
		{
			tags:    []string{"2", "1"},
			reverse: false,
		},
		{
			tags:    []string{"2", "1"},
			reverse: true,
		},
		{

			tags:    []string{"abcd", "abcde", "abcdf"},
			reverse: false,
		},
		{
			tags:    []string{"abcd", "abcde", "abcdf"},
			reverse: true,
		},
		{
			tags:    []string{"2", "abcd", "abcde", "abcdf", "1"},
			reverse: false,
		},
		{
			tags:    []string{"2", "abcd", "abcde", "abcdf", "1"},
			reverse: true,
		},
	}
	for _, c := range cazes {
		tagIndex := ByEquality("tag")
		if c.reverse {
			tagIndex.Order.Type = OrderTypeDesc
		}
		tagIndex.StringOrderPadLength = 12
		table := New(User{}, &Options{
			Store:     fs.NewStore(),
			Indexes:   []Index{tagIndex},
			Namespace: uuid.Must(uuid.NewV4()).String(),
		})
		for _, key := range c.tags {
			err := table.Create(User{
				ID:  uuid.Must(uuid.NewV4()).String(),
				Tag: key,
			})
			if err != nil {
				t.Fatal(err)
			}
		}
		users := []User{}
		q := QueryEquals("tag", nil)
		if c.reverse {
			q.Order.Type = OrderTypeDesc
		}
		err := table.Read(q, &users)
		if err != nil {
			t.Fatal(err)
		}

		tags := sort.StringSlice(c.tags)
		sort.Sort(tags)
		if c.reverse {
			reverse(tags)
		}
		if len(tags) != len(users) {
			t.Fatal(tags, users)
		}
		for i, key := range tags {
			if users[i].Tag != key {
				userTags := []string{}
				for _, v := range users {
					userTags = append(userTags, v.Tag)
				}
				t.Fatalf("Should be %v, got %v, is reverse: %v", tags, userTags, c.reverse)
			}
		}
	}

}

func reverseInt(is []int) {
	last := len(is) - 1
	for i := 0; i < len(is)/2; i++ {
		is[i], is[last-i] = is[last-i], is[i]
	}
}

func TestOrderingNumbers(t *testing.T) {
	type caze struct {
		dates   []int
		reverse bool
	}
	cazes := []caze{
		{
			dates:   []int{20, 30},
			reverse: false,
		},
		{
			dates:   []int{20, 30},
			reverse: true,
		},
	}
	for _, c := range cazes {
		createdIndex := ByEquality("created")
		if c.reverse {
			createdIndex.Order.Type = OrderTypeDesc
		}
		table := New(User{}, &Options{
			Store:     fs.NewStore(),
			Indexes:   []Index{createdIndex},
			Namespace: uuid.Must(uuid.NewV4()).String(),
		})
		for _, key := range c.dates {
			err := table.Create(User{
				ID:      uuid.Must(uuid.NewV4()).String(),
				Created: int64(key),
			})
			if err != nil {
				t.Fatal(err)
			}
		}
		users := []User{}
		q := QueryEquals("created", nil)
		if c.reverse {
			q.Order.Type = OrderTypeDesc
		}
		err := table.Read(q, &users)
		if err != nil {
			t.Fatal(err)
		}

		dates := sort.IntSlice(c.dates)
		sort.Sort(dates)
		if c.reverse {
			reverseInt([]int(dates))
		}
		if len(users) != len(dates) {
			t.Fatalf("Expected %v, got %v", len(dates), len(users))
		}
		for i, date := range dates {
			if users[i].Created != int64(date) {
				userDates := []int{}
				for _, v := range users {
					userDates = append(userDates, int(v.Created))
				}
				t.Fatalf("Should be %v, got %v, is reverse: %v", dates, userDates, c.reverse)
			}
		}
	}

}

func TestStaleIndexRemoval(t *testing.T) {
	tagIndex := ByEquality("tag")
	table := New(User{}, &Options{
		Store:     fs.NewStore(),
		Indexes:   []Index{tagIndex},
		Namespace: uuid.Must(uuid.NewV4()).String(),
	})
	err := table.Create(User{
		ID:  "1",
		Tag: "hi-there",
	})
	if err != nil {
		t.Fatal(err)
	}
	err = table.Create(User{
		ID:  "1",
		Tag: "hello-there",
	})
	if err != nil {
		t.Fatal(err)
	}
	res := []User{}
	err = table.Read(QueryEquals("tag", nil), &res)
	if err != nil {
		t.Fatal(err)
	}
	if len(res) > 1 {
		t.Fatal(res)
	}
}

func TestUniqueIndex(t *testing.T) {
	tagIndex := ByEquality("tag")
	tagIndex.Unique = true
	table := New(User{}, &Options{
		Store:     fs.NewStore(),
		Indexes:   []Index{tagIndex},
		Namespace: uuid.Must(uuid.NewV4()).String(),
	})
	err := table.Create(User{
		ID:  "1",
		Tag: "hi-there",
	})
	if err != nil {
		t.Fatal(err)
	}
	err = table.Create(User{
		ID:  "2",
		Tag: "hello-there",
	})
	if err != nil {
		t.Fatal(err)
	}
	err = table.Create(User{
		ID:  "3",
		Tag: "hi-there",
	})
	if err == nil {
		t.Fatal("Create shoud fail with duplicate tag error because the index is unique")
	}
}

type Tag struct {
	Slug string `json:"slug"`
	Age  int    `json:"age"`
	Type string `json:"type"`
}

func TestNonIDKeys(t *testing.T) {
	table := New(Tag{}, &Options{
		Store:     fs.NewStore(),
		Key:       "slug",
		Namespace: uuid.Must(uuid.NewV4()).String(),
	})

	err := table.Create(Tag{
		Slug: "1",
		Age:  12,
	})
	if err != nil {
		t.Fatal(err)
	}
	err = table.Create(Tag{
		Slug: "2",
		Age:  25,
	})
	if err != nil {
		t.Fatal(err)
	}
	users := []User{}
	q := QueryEquals("slug", "1")
	q.Order.Type = OrderTypeUnordered
	err = table.Read(q, &users)
	if err != nil {
		t.Fatal(err)
	}
	if len(users) != 1 {
		t.Fatal(users)
	}
}

// This might be an almost duplicate test, I used it to try reproduce an issue
// Leaving this here for now as we dont have enough tests anyway.
func TestReadByString(t *testing.T) {
	typeIndex := ByEquality("type")

	table := New(Tag{}, &Options{
		Key:       "slug",
		Store:     fs.NewStore(),
		Indexes:   []Index{typeIndex},
		Debug:     false,
		Namespace: uuid.Must(uuid.NewV4()).String(),
	})

	err := table.Create(Tag{
		Slug: "1",
		Type: "post-tag",
	})
	if err != nil {
		t.Fatal(err)
	}
	err = table.Create(Tag{
		Slug: "2",
		Type: "post-tag",
	})
	if err != nil {
		t.Fatal(err)
	}
	tags := []Tag{}
	q := QueryEquals("type", "post-tag")
	err = table.Read(q, &tags)
	if err != nil {
		t.Fatal(err)
	}
	if len(tags) != 2 {
		t.Fatal(tags)
	}
}

func TestOderByDifferentFieldThanFilterField(t *testing.T) {
	typeIndex := ByEquality("type")
	typeIndex.Order = Order{
		Type:      OrderTypeDesc,
		FieldName: "age",
	}

	table := New(Tag{}, &Options{
		Key:       "slug",
		Store:     fs.NewStore(),
		Indexes:   []Index{typeIndex},
		Namespace: uuid.Must(uuid.NewV4()).String(),
		Debug:     false,
	})

	err := table.Create(Tag{
		Slug: "1",
		Type: "post-tag",
		Age:  15,
	})
	if err != nil {
		t.Fatal(err)
	}
	err = table.Create(Tag{
		Slug: "2",
		Type: "post-tag",
		Age:  25,
	})
	if err != nil {
		t.Fatal(err)
	}
	err = table.Create(Tag{
		Slug: "3",
		Type: "other-tag",
		Age:  30,
	})
	if err != nil {
		t.Fatal(err)
	}
	tags := []Tag{}
	err = table.Read(typeIndex.ToQuery("post-tag"), &tags)
	if err != nil {
		t.Fatal(err)
	}
	if len(tags) != 2 {
		t.Fatal(tags)
	}
	if tags[0].Age != 25 {
		t.Fatal(tags)
	}
	if tags[1].Age != 15 {
		t.Fatal(tags)
	}

	err = table.Read(typeIndex.ToQuery(nil), &tags)
	if err != nil {
		t.Fatal(err)
	}
	if len(tags) != 3 {
		t.Fatal(tags)
	}
}

func TestDeleteIndexCleanup(t *testing.T) {
	slugIndex := ByEquality("slug")
	slugIndex.Order.Type = OrderTypeUnordered

	typeIndex := ByEquality("type")

	table := New(Tag{}, &Options{
		Key:       "slug",
		Store:     fs.NewStore(),
		Indexes:   []Index{typeIndex},
		Namespace: uuid.Must(uuid.NewV4()).String(),
		Debug:     false,
	})

	err := table.Create(Tag{
		Slug: "1",
		Type: "post-tag",
	})
	if err != nil {
		t.Fatal(err)
	}
	err = table.Create(Tag{
		Slug: "2",
		Type: "post-tag",
	})
	if err != nil {
		t.Fatal(err)
	}
	tags := []Tag{}
	q := QueryEquals("type", "post-tag")
	err = table.Read(q, &tags)
	if err != nil {
		t.Fatal(err)
	}
	if len(tags) != 2 {
		t.Fatal(tags)
	}

	err = table.Delete(slugIndex.ToQuery("1"))
	if err != nil {
		t.Fatal(err)
	}

	q = QueryEquals("type", "post-tag")
	err = table.Read(q, &tags)
	if err != nil {
		t.Fatal(err)
	}
	if len(tags) != 1 {
		t.Fatal(tags)
	}
}

func TestDeleteByUnmatchingIndex(t *testing.T) {
	table := New(User{}, &Options{
		Store:     fs.NewStore(),
		Namespace: uuid.Must(uuid.NewV4()).String(),
		Debug:     false,
	})

	err := table.Create(User{
		ID:  "1",
		Age: 20,
	})
	if err != nil {
		t.Fatal(err)
	}
	err = table.Create(User{
		ID:  "2",
		Age: 30,
	})
	if err != nil {
		t.Fatal(err)
	}

	err = table.Delete(QueryEquals("ID", "1"))
	if err != nil {
		t.Fatal(err)
	}
	t.Run("Test read by unspecified index", func(t *testing.T) {
		users := []User{}
		err = table.Read(QueryEquals("ID", "1"), &users)
		if err != nil {
			t.Fatal(err)
		}
		if len(users) != 0 {
			t.Fatal(users)
		}
		err = table.Read(QueryEquals("ID", "2"), &users)
		if err != nil {
			t.Fatal(err)
		}
		if len(users) != 1 {
			t.Fatal(users)
		}
	})
}

func TestDeleteByUnmatchingIndexMap(t *testing.T) {
	table := New(map[string]interface{}{}, &Options{
		Store:     fs.NewStore(),
		Namespace: uuid.Must(uuid.NewV4()).String(),
		Debug:     false,
	})

	err := table.Create(map[string]interface{}{
		"ID":  "1",
		"Age": 20,
	})
	if err != nil {
		t.Fatal(err)
	}
	err = table.Create(map[string]interface{}{
		"ID":  "2",
		"Age": 30,
	})
	if err != nil {
		t.Fatal(err)
	}

	err = table.Delete(QueryEquals("ID", "1"))
	if err != nil {
		t.Fatal(err)
	}
	t.Run("Test read by unspecified index", func(t *testing.T) {
		users := []map[string]interface{}{}
		err = table.Read(QueryEquals("ID", "1"), &users)
		if err != nil {
			t.Fatal(err)
		}
		if len(users) != 0 {
			t.Fatal(users)
		}
		err = table.Read(QueryEquals("ID", "2"), &users)
		if err != nil {
			t.Fatal(err)
		}
		if len(users) != 1 {
			t.Fatal(users)
		}
	})
}

func TestUpdateDeleteIndexMaintenance(t *testing.T) {
	updIndex := ByEquality("updated")
	updIndex.Order.Type = OrderTypeDesc

	table := New(User{}, &Options{
		Store:     fs.NewStore(),
		Indexes:   []Index{updIndex},
		Namespace: uuid.Must(uuid.NewV4()).String(),
		Debug:     false,
	})

	err := table.Create(User{
		ID:      "1",
		Age:     12,
		Updated: 5000,
	})
	if err != nil {
		t.Fatal(err)
	}

	err = table.Create(User{
		ID:      "2",
		Age:     25,
		Updated: 5001,
	})
	if err != nil {
		t.Fatal(err)
	}
	users := []User{}
	q := updIndex.ToQuery(nil)
	err = table.Read(q, &users)
	if err != nil {
		t.Fatal(err)
	}
	if len(users) != 2 {
		t.Fatal(users)
	}
	if users[0].ID != "2" || users[1].ID != "1" {
		t.Fatal(users)
	}

	err = table.Create(User{
		ID:      "1",
		Age:     12,
		Updated: 5002,
	})
	if err != nil {
		t.Fatal(err)
	}

	err = table.Read(q, &users)
	if err != nil {
		t.Fatal(err)
	}
	if len(users) != 2 {
		t.Fatal(users)
	}
	if users[0].ID != "1" || users[1].ID != "2" {
		t.Fatal(users)
	}
}

func TestUpdateDeleteIndexMaintenanceMap(t *testing.T) {
	updIndex := ByEquality("Updated")
	updIndex.Order.Type = OrderTypeDesc

	table := New(map[string]interface{}{}, &Options{
		Store:     fs.NewStore(),
		Indexes:   []Index{updIndex},
		Namespace: uuid.Must(uuid.NewV4()).String(),
		Debug:     false,
	})

	err := table.Create(map[string]interface{}{
		"ID":      "1",
		"Age":     "12",
		"Updated": "5000",
	})
	if err != nil {
		t.Fatal(err)
	}

	err = table.Create(map[string]interface{}{
		"ID":      "2",
		"Age":     "25",
		"Updated": "5001",
	})
	if err != nil {
		t.Fatal(err)
	}
	users := []map[string]interface{}{}
	q := updIndex.ToQuery(nil)
	err = table.Read(q, &users)
	if err != nil {
		t.Fatal(err)
	}
	if len(users) != 2 {
		t.Fatal(users)
	}
	if users[0]["ID"] != "2" || users[1]["ID"] != "1" {
		t.Fatal(users)
	}

	err = table.Create(map[string]interface{}{
		"ID":      "1",
		"Age":     "12",
		"Updated": "5002",
	})
	if err != nil {
		t.Fatal(err)
	}

	err = table.Read(q, &users)
	if err != nil {
		t.Fatal(err)
	}
	if len(users) != 2 {
		t.Fatal(users)
	}
	if users[0]["ID"] != "1" || users[1]["ID"] != "2" {
		t.Fatal(users)
	}
}

type TypeTest struct {
	ID  string `json:"ID"`
	F32 float32
	F64 float64
	I   int
	I32 int32
	I64 int64
	S   string
	B   bool
}

// Test aimed specifically to test all types
func TestAllCombos(t *testing.T) {
	// go over all filter + order combos
	// for equality indexing

	v := reflect.ValueOf(TypeTest{})
	for filterFieldI := 0; filterFieldI < v.NumField(); filterFieldI++ {
		filterField := v.Field(filterFieldI)
		for orderFieldI := 0; orderFieldI < v.NumField(); orderFieldI++ {
			orderField := v.Field(orderFieldI)

			filterFieldName := v.Type().Field(filterFieldI).Name
			orderFieldName := v.Type().Field(orderFieldI).Name
			if filterFieldName == "ID" || orderFieldName == "ID" {
				continue
			}
			if filterFieldName == orderFieldName {
				continue
			}

			t.Run(fmt.Sprintf("Filter by %v, order by %v ASC", filterField.Type().Name(), orderField.Type().Name()), func(t *testing.T) {
				index := ByEquality(filterFieldName)
				index.Order.Type = OrderTypeAsc
				index.Order.FieldName = orderFieldName

				table := New(TypeTest{}, &Options{
					Store:     fs.NewStore(),
					Indexes:   []Index{index},
					Namespace: uuid.Must(uuid.NewV4()).String(),
					Debug:     false,
				})

				small := TypeTest{
					ID: "1",
				}
				v1 := getExampleValue(table.(*model).getFieldValue(small, orderFieldName), 1)
				table.(*model).setFieldValue(&small, orderFieldName, v1)

				large := TypeTest{
					ID: "2",
				}
				v2 := getExampleValue(table.(*model).getFieldValue(large, orderFieldName), 2)
				table.(*model).setFieldValue(&large, orderFieldName, v2)

				err := table.Create(small)
				if err != nil {
					t.Fatal(err)
				}
				err = table.Create(large)
				if err != nil {
					t.Fatal(err)
				}
				results := []TypeTest{}
				err = table.Read(QueryEquals(filterFieldName, nil), &results)
				if err != nil {
					t.Fatal(err)
				}
				if len(results) < 2 {
					t.Fatal(results)
				}
				if results[0].ID != "1" || results[1].ID != "2" {
					t.Fatal("Results:", results, results[0].ID, results[1].ID)
				}
			})
			t.Run(fmt.Sprintf("Filter by %v, order by %v DESC", filterField.Type().Name(), orderField.Type().Name()), func(t *testing.T) {
				index := ByEquality(filterFieldName)
				index.Order.Type = OrderTypeDesc
				index.Order.FieldName = orderFieldName

				table := New(TypeTest{}, &Options{
					Store:     fs.NewStore(),
					Indexes:   []Index{index},
					Namespace: uuid.Must(uuid.NewV4()).String(),
					Debug:     false,
				})

				small := TypeTest{
					ID: "1",
				}
				v1 := getExampleValue(table.(*model).getFieldValue(small, orderFieldName), 1)
				table.(*model).setFieldValue(&small, orderFieldName, v1)

				large := TypeTest{
					ID: "2",
				}
				v2 := getExampleValue(table.(*model).getFieldValue(large, orderFieldName), 2)
				table.(*model).setFieldValue(&large, orderFieldName, v2)

				err := table.Create(small)
				if err != nil {
					t.Fatal(err)
				}
				err = table.Create(large)
				if err != nil {
					t.Fatal(err)
				}
				results := []TypeTest{}
				err = table.Read(index.ToQuery(nil), &results)
				if err != nil {
					t.Fatal(err)
				}
				if len(results) < 2 {
					t.Fatal(results)
				}
				if results[0].ID != "2" || results[1].ID != "1" {
					t.Fatal("Results:", results, results[0].ID, results[1].ID)
				}
			})
		}
	}
}

// returns an example value generated
// nth = each successive number should cause this
// function to return a "bigger value" for each type
func getExampleValue(i interface{}, nth int) interface{} {
	switch v := i.(type) {
	case string:
		return strings.Repeat("a", nth)
	case bool:
		if nth == 1 {
			return false
		}
		return true
	case float32:
		return v + float32(nth) + .1
	case float64:
		return v + float64(nth) + 0.1
	case int:
		return v + nth
	case int32:
		return v + int32(nth)
	case int64:
		return v + int64(nth)
	}
	return nil
}
