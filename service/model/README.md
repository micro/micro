# Model

Package model is a convenience wrapper around what the `Store` provides.
It's main responsibility is to maintain indexes that would otherwise be maintaned by the users to enable different queries on the same data.

## Usage

The following snippets will this piece of code prepends them.

```go
import(
    model "github.com/micro/micro/v3/service/model"
    fs "github.com/micro/micro/v3/service/store/file"
)

type User struct {
	ID      string `json:"id"`
 	Name string    `json:"name"`
	Age     int    `json:"age"`
	HasPet  bool   `json:"hasPet"`
	Created int64  `json:"created"`
	Tag     string `json:"tag"`
	Updated int64  `json:"updated"`
}
```

## Query by field equality

For each field we want to query on we have to create an index. Index by `id` is provided by default to each `DB`, there is no need to specify it.

```go
ageIndex := model.ByEquality("age")

db := model.New(fs.NewStore(), User{}, []model.Index{(ageIndex})

err := db.Create(User{
    ID: "1",
    Name: "Alice",
    Age: 20,
})
if err != nil {
    // handle save error
}
err := db.Create(User{
    ID: "2",
    Name: "Jane",
    Age: 22
})
if err != nil {
    // handle save error
}

err = db.Read(model.Equals("age", 22), &users)
if err != nil {
	// handle list error
}
fmt.Println(users)

// will print
// [{"id":"2","name":"Jane","age":22}]
```

## Reading all records in an index

Reading can be done without specifying a value:

```go
db.Read(Equals("age", nil), &users)
```

Readings will be unordered, ascending ordered or descending ordered depending on the ordering settings of the index.

## Ordering

Indexes by default are ordered. If we want to turn this behaviour off:

```go
ageIndex.Order.Type = OrderTypeUnordered

ageQuery := model.Equals("age", 22)
ageQuery.Order.Type = OrderTypeUnordered
```

### Filtering by one field, ordering by other

```go
typeIndex := ByEquality("type")
typeIndex.Order = Order{
	Type:      OrderTypeDesc,
	FieldName: "age",
}

// Results will be ordered by age
db.Read(typeIndex.ToQuery("a-certain-type-value"))
```

By default the ordering field is the same as the filtering field.

### Reverse order

```go
ageQuery.Desc = true
```

### Queries must match indexes

It is important to note that queries must match indexes. The following index-query pairs match (separated by an empty line)

```go
// Ascending ordered index by age
index := model.Equality("age")
// Read ascending ordered by age
query := model.Equals("age", nil)
// Read ascending ordered by age where age = 20
query2 := model.Equals("age", 20) 

// Descending ordered index by age
index := model.Equality("age")
index.Order.Type = OrderTypeDesc
// Read descending ordered by age
query := model.Equals("age", nil)
query.Order.Type = OrderTypeDesc
// Read descending ordered by age where age = 20
query2 := model.Equals("age", 20)
query2.Order.Type = OrderTypeDesc

// Unordered index by age
index := model.Equality("age")
index.Order.Type = OrderTypeUnordered
// Read unordered by age
query := model.Equals("age", nil)
query.Order.Type = OrderTypeUnordered
// Read unordered by age where age = 20
query2 := model.Equals("age", 20)
query2.Order.Type = OrderTypeUnordered
```

Of course, maintaining this might be inconvenient, for this reason the `ToQuery` method was introduced, see below.

#### Creating a query out of an Index

```go

index := model.Equality("age")
index.Order.Type = OrderTypeUnordered

db.Read(index.ToQuery(25))
```

### Unordered listing without value

It's easy to see how listing things by unordered indexes on different fields should result in the same output: a randomly ordered list, ie:

```go
ageIndex := model.Equality("age")
ageIndex.Order.Type = OrderTypeUnordered

emailIndex := model.Equality("email")
emailIndex.Order.Type = OrderTypeUnordered

result1 := []User{}
result2 := []User{}

db.Read(model.Equals("age"), &result1)
db.Read(model.Equals("email"), &result2)

// Both result1 and result2 will be an unordered listing without
// filtering on either the age or email fields.
// Could be thought of as a noop query despite not having an explicit "no query" listing.
```

### Ordering by string fields

Ordering comes for "free" when dealing with numeric or boolean fields, but it involves  in padding, inversing and order preserving base32 encoding of values to work for strings.

This can sometimes result in large keys saved, as the inverse of a small 1 byte character in a string is a 4 byte rune. Optionally adding base32 encoding on top to prevent exotic runes appearing in keys, strings blow up in size even more. If saving space is a requirement and ordering is not, ordering for strings should be turned off.

The matter is further complicated by the fact that the padding size must be specified ahead of time.

```go
nameIndex := model.ByEquality("name")
nameIndex.StringOrderPadLength = 10

nameQuery := model.Equals("age", 22)
// `StringOrderPadLength` is not needed to be specified for the query
```

To turn off base32 encoding and keep the runes:

```go
nameIndex.Base32Encode = false
```

## Unique indexes

```go
emailIndex := model.ByEquality("email")
emailIndex.Unique = true
```

## Design

### Restrictions

To maintain all indexes properly, all fields must be filled out when saving.
This sometimes requires a `Read, Modify, Write` pattern. In other words, partial updates will break indexes.

This could be avoided later if model does the loading itself.

## TODO

- Implement deletes
- Implement counters, for pattern inspiration see the [tags service](https://github.com/micro/services/tree/master/blog/tags)
- Test boolean indexes and its ordering
- There is a stuttering in the way `id` fields are being saved twice. ID fields since they are unique do not need `id` appended after them in the record keys.
