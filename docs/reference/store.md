---
title: Store
keywords: micro
tags: [micro]
sidebar: home_sidebar
permalink: /reference/store
summary: Using the store, using key-value stores efficiently
---

Micro's store interface is a key value store with support for odering of keys.

## Key-value stores in general

### Tradeoffs

Key-value stores that support ordering of keys can be used to build complex applications.
Due to their very limited feature set, key-value stores generally scale easily and reliably, often linearly with the number of nodes added.

This scalability comes at the expense of inconvenience and mental overhead when writing business logic. For usecases where linear scalability is important, this tradeoff is preferred.

### How KV stores can enable complex applications

As it was mentioned, KV stores ofsten support the "ordering of keys" feature which is crucial for building nontrivial applications. How would that work?

Let's say we have the following data in our table:

```sh
# contents of table "users"
KEY         VALUE
id1         {"id":"id1", "name":"Jane", "class":"firstGrade",   "avgScore": 98}
id2         {"id":"id2", "name":"Alice","class":"secondGrade",  "avgScore": 92}
id3         {"id":"id3", "name":"Joe",  "class":"secondGrade"   "avgScore": 89}
id4         {"id":"id4", "name":"Betty","class":"thirdGrade"    "avgScore": 94}
```

Given this data structure, we can do two queries:

- reading a given key (get "id1", get "id2")
- if the keys are ordered, we can ask for X number of entries after a key (get 3 entries after "id2")

Finding values in an ordered set is possibly the simplest task we can ask a database.
The problem with our example is that it's not very useful to ask "find me keys coming in the order after "id2".

We usually don't care about ids (apart from it enabling read by id queries), as in most cases they are uuids in the form of `something-very-long-and-random`. So what we need to do is save the data with the query in mind.

In the case of our schoold students, let's say we wan't to list by class. To do this, having the query in mind, we can copy the data over to an other table named after the query we want to do:

```sh
# contents of table "usersByClass"
KEY             VALUE
firstGrade/id1  {"id":"id1", "name":"Jane", "class":"firstGrade",   "avgScore": 98}
secondGrade/id2 {"id":"id2", "name":"Alice","class":"secondGrade",  "avgScore": 92}
secondGrade/id3 {"id":"id3", "name":"Joe",  "class":"secondGrade"   "avgScore": 89}
thirdGrade/id4  {"id":"id4", "name":"Betty","class":"thirdGrade"    "avgScore": 94}
```

Since the keys are ordered it is very trivial to get back let's say "all second graders".
Key value stores which have their keys ordered support something similar to "key starts with/key has prefix" query. In the case of second graders, listing all records where the "keys start with `secondGrade`" will give us back all the second graders.

This query is basically a `field equals to` as we essentially did a `field class == secondGrade`. But we could also exploit the ordered nature of the keys to do a value comparison query, ie `field avgScores is less than 90` or `field AvgScores is between 90 and 95` etc., if we model our data appropriately:

```sh
# contents of table "usersByAvgScore"
KEY         VALUE
089/id3     {"id":"id3", "name":"Joe",  "class":"secondGrade"   "avgScore": 89}
092/id2     {"id":"id2", "name":"Alice","class":"secondGrade",  "avgScore": 92}
094/id4     {"id":"id4", "name":"Betty","class":"thirdGrade"    "avgScore": 94}
098/id1     {"id":"id1", "name":"Jane", "class":"firstGrade",   "avgScore": 98}
```

It's worth remembering that the keys are strings, and that they are ordered lexicographically. For this reason when dealing with numbering values, we must make sure that they are prepended to the same length appropriately.

## Micro's key-value store