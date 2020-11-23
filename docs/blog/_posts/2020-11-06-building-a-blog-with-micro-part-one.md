---
layout: post
author: Janos Dobronszki
title: Building a Blog with Micro - Part One
keywords: tutorials, blog
tags: [blog]
---

This series will cover how to build a blog service using Micro. We'll decompose a monolithic Blog into multiple services. 
In part one we'll focus on building a Post service. It will be good way to learn how to build nontrivial applications with 
the [store](https://micro.mu/reference#store) and the [model](https://github.com/micro/dev/tree/master/model).

The most important takeaway from this post will likely be the the usage of the key-value store for non-trivial usecases 
(such as querying blog posts by slug and listing them by reverse creation order).

## The Basics

Head to the [Getting Started Guide](/getting-started) if you haven't used Micro before.

If you have let's use that knowledge! As a reminder, we have to make sure `micro server` is running in an other terminal, 
and we are connected to it, ie


Running the micro server

```sh
micro server
```

Looking up our local environment

```sh
$ micro env
* local      127.0.0.1:8081         Local running micro server
  dev        proxy.m3o.dev          Cloud hosted development environment
  platform   proxy.m3o.com          Cloud hosted production environment
```

We can see the local environment picked. If not, we can issue `micro env set local` to remedy.   

Now back to the `micro new` command:

```sh
$ micro new posts
$ ls posts
Dockerfile	Makefile	README.md	generate.go	go.mod		handler		main.go		proto
```

Great! The best way to start a service is to define the proto. The generated default should be something similar to this:

In our post service, we want 3 methods:
- `Save` for blog insert and update
- `Query` for reading and listing
- `Delete` for deletion

Let's start with the post method.

<script src="https://gist.github.com/asim/31628c6faf61b392b31bcefc79daa289.js"></script>

Astute readers might notice that although we have defined a `Post` message type, we still redefine some of the fields as top level fields for the `SaveRequest` message type.
The main reason for this is that we don't want our [dynamic commands](https://micro.mu/reference#dynamic-commands).

Ie. if we would embed a `Post post = 1` inside `SaveRequest`, we would call the posts service the following way:

```sh
micro posts save --post_title=Title --post_content=Content
```

but we don't want to keep repeating `post`, our preferred way is:

```sh
micro posts save --title=Title --content=Content
```

To regenerate the proto, we have to issue the `make proto` command in the project root.

Now, the `main.go`:

<script src="https://gist.github.com/asim/99d3bea6fdd234b373473309cab055a1.js"></script>

After that's done, let's adjust the handler to match our proto! This snippet is a bit longer, so cover it piece by piece:

<script src="https://gist.github.com/asim/10b96ed42e300959ec96e963ed4cc6e2.js"></script>

The above piece of code uses the [model package](https://github.com/micro/dev/tree/master/model). It sets up the indexes which will enable us to query the data and also tells model to maintain these indexes.

- The id index is needed to read by id
- The created index is needed so when we list posts the order of the posts will be descending based on the created field
- The slug index is needed to we can read posts by slug (ie. `myblog.com/post/awesome-post-url`)

At this point `micro run .` in project root should deploy our post service. Let's verify with `micro logs posts`:

```
$ micro logs posts
Starting [service] posts
Server [grpc] Listening on [::]:53031
Registry [service] Registering node: posts-b36361ae-f2ae-48b0-add5-a8d4797508be
```

(The exact output might depend on the actual config format configuraton.)

## Saving posts

Let's make our service do something useful now: save a post.

<script src="https://gist.github.com/asim/a1870dccc4ad1eb137a602ae9e884af6.js"></script>

After a `micro update .` in project root, we can start saving posts!

```
micro posts save --id=1 --title="Post one" --content="First saved post"
micro posts save --id=2 --title="Post two" --content="Second saved post"
```

## Querying posts

Again, implementation starts with defining the protos:

<script src="https://gist.github.com/asim/d3dea1832609b32538838529e693f9bf.js"></script>

A `make proto` issued in the command root should regenerate the Go proto files and we should be ready to define our new handler:

We want our query handler to enable querying by id, slug and also enable listing of posts:
<script src="https://gist.github.com/asim/4cc45b6fceb459332a61bed0f88a99d6.js"></script>

As mentioned, the existing indexes can be used for querying too with the `ToQuery` method.

After doing a `micro update .` in the project root, we can now query the posts:

```
$ micro posts query
{
	"posts": [
		{
			"id": "2",
			"title": "Post two",
			"slug": "post-two",
			"content": "Second saved post",
			"created": "1604423363"
		},
		{
			"id": "1",
			"title": "Post one",
			"slug": "post-one",
			"content": "First saved post",
			"created": "1604423297"
		}
	]
}

```

Stellar! Now only `Delete` remains to be implemented to have a basic post service.

## Deleting posts

Since we have already defined `Delete` in our proto, we only have to implement the handler. It is rather simple:

<script src="https://gist.github.com/asim/a3c7875a5fc47439b43740b274ab3ac3.js"></script>

## Conclusions

This brings us to the end of the first post in the blogs tutorial series.
There are many more features we will add later, like saving and querying by tags, but this post already taught us enough to digest.
We will cover those aspect in later parts of this series.

The source code for this can be found [here](https://github.com/micro/dev/tree/master/blog/v1-posts).
Further versions will be in the same `blog` folder with different versions, ie `v2-posts` and once we have more services, `v2-tags`, `v2-comments`.
Folders with the same prefix will be meant to be deployed together, but more on this later.

