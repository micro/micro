---
title: Store
keywords: micro
tags: [micro]
sidebar: home_sidebar
permalink: /reference/store
summary: Using the store, using key-value stores efficiently
---

For a good beginner level doc on the Store, please see the [helloworld tutorial](/helloworld).

# The key-value store

Micro's store interface is a key value store with support for odering of keys.

## Key-value stores in general

Key-value stores that support ordering of keys can be used to build complex applications.
Due to their very limited feature set, key-value stores generally scale easily and reliably, often linearly with the number of nodes added.

This scalability comes at the expense of inconvenience and mental overhead when writing business logic. For usecases where linear scalability is important, this tradeoff is preferred.

## Query by ID

Reading by ID is the archetypical job for key value stores. Storing data to enable this ID works just like in any other database:

```sh
# entries designed for querying "users by id"
KEY         VALUE
id1         {"id":"id1", "name":"Jane", "class":"firstGrade",   "avgScore": 98}
id2         {"id":"id2", "name":"Alice","class":"secondGrade",  "avgScore": 92}
id3         {"id":"id3", "name":"Joe",  "class":"secondGrade"   "avgScore": 89}
id4         {"id":"id4", "name":"Betty","class":"thirdGrade"    "avgScore": 94}
```

```go
import "github.com/micro/micro/v3/service/store"

records, err := store.Read("id1")
if err != nil {
	fmt.Println("Error reading from store: ", err)
}
fmt.Println(records[0].Value)
// Will output {"id":"id1", "name":"Jane", "class":"firstGrade",   "avgScore": 98}
```

Given this data structure, we can do two queries:

- reading a given key (get "id1", get "id2")
- if the keys are ordered, we can ask for X number of entries after a key (get 3 entries after "id2")

Finding values in an ordered set is possibly the simplest task we can ask a database.
The problem with the above data structure is that it's not very useful to ask "find me keys coming in the order after "id2". To enable other kind of queries, the data must be saved with different keys.

In the case of the schoold students, let's say we wan't to list by class. To do this, having the query in mind, we can copy the data over to an other table named after the query we want to do:

## Query by field value equality

```sh
# entries designed for querying "users by class"
KEY             VALUE
firstGrade/id1  {"id":"id1", "name":"Jane", "class":"firstGrade",   "avgScore": 98}
secondGrade/id2 {"id":"id2", "name":"Alice","class":"secondGrade",  "avgScore": 92}
secondGrade/id3 {"id":"id3", "name":"Joe",  "class":"secondGrade"   "avgScore": 89}
thirdGrade/id4  {"id":"id4", "name":"Betty","class":"thirdGrade"    "avgScore": 94}
```


```go
import "github.com/micro/micro/v3/service/store"

records, err := store.Read("", store.Prefix("secondGrade"))
if err != nil {
	fmt.Println("Error reading from store: ", err)
}
fmt.Println(records[0].Value)
// Will output
// secondGrade/id2 {"id":"id2", "name":"Alice","class":"secondGrade",  "avgScore": 92}
// secondGrade/id3 {"id":"id3", "name":"Joe",  "class":"secondGrade"   "avgScore": 89}
```

Since the keys are ordered it is very trivial to get back let's say "all second graders".
Key value stores which have their keys ordered support something similar to "key starts with/key has prefix" query. In the case of second graders, listing all records where the "keys start with `secondGrade`" will give us back all the second graders.

This query is basically a `field equals to` as we essentially did a `field class == secondGrade`. But we could also exploit the ordered nature of the keys to do a value comparison query, ie `field avgScores is less than 90` or `field AvgScores is between 90 and 95` etc., if we model our data appropriately:

## Query for field value ranges

```sh
# entries designed for querying "users by avgScore"
KEY         VALUE
089/id3     {"id":"id3", "name":"Joe",  "class":"secondGrade"   "avgScore": 89}
092/id2     {"id":"id2", "name":"Alice","class":"secondGrade",  "avgScore": 92}
094/id4     {"id":"id4", "name":"Betty","class":"thirdGrade"    "avgScore": 94}
098/id1     {"id":"id1", "name":"Jane", "class":"firstGrade",   "avgScore": 98}
```

It's worth remembering that the keys are strings, and that they are ordered lexicographically. For this reason when dealing with numbering values, we must make sure that they are prepended to the same length appropriately.

At the moment Micro's store does not support this kind of query, this example is only here to hint at future possibilities with the store.

## Tables and avoiding key collisions

Micro services only have access to one Store table. This means all keys take live in the same namespace and can collide. A very useful pattern is to separate the entries by their intended query pattern, ie taking the "users by id" and users by class records above:

```sh
KEY         VALUE
# entries designed for querying "users by id"
usersById/id1         		{"id":"id1", "name":"Jane", "class":"firstGrade",   "avgScore": 98}
usersById/id2         		{"id":"id2", "name":"Alice","class":"secondGrade",  "avgScore": 92}
usersById/id3         		{"id":"id3", "name":"Joe",  "class":"secondGrade"   "avgScore": 89}
usersById/id4         		{"id":"id4", "name":"Betty","class":"thirdGrade"    "avgScore": 94}
# entries designed for querying "users by class"
usersByClass/firstGrade/id1  {"id":"id1", "name":"Jane", "class":"firstGrade",   "avgScore": 98}
usersByClass/secondGrade/id2 {"id":"id2", "name":"Alice","class":"secondGrade",  "avgScore": 92}
usersByClass/secondGrade/id3 {"id":"id3", "name":"Joe",  "class":"secondGrade"   "avgScore": 89}
usersByClass/thirdGrade/id4  {"id":"id4", "name":"Betty","class":"thirdGrade"    "avgScore": 94}
```

Respective go examples this way become:

```go
import "github.com/micro/micro/v3/service/store"

const idPrefix = "usersById/"

records, err := store.Read(idPrefix + "id1")
if err != nil {
	fmt.Println("Error reading from store: ", err)
}
fmt.Println(records[0].Value)
// Will output {"id":"id1", "name":"Jane", "class":"firstGrade",   "avgScore": 98}
```

```go
import "github.com/micro/micro/v3/service/store"

const classPrefix = "usersByClass/"

records, err := store.Read("", store.Prefix(classPrefix + "secondGrade"))
if err != nil {
	fmt.Println("Error reading from store: ", err)
}
fmt.Println(records[0].Value)
// Will output
// secondGrade/id2 {"id":"id2", "name":"Alice","class":"secondGrade",  "avgScore": 92}
// secondGrade/id3 {"id":"id3", "name":"Joe",  "class":"secondGrade"   "avgScore": 89}
```