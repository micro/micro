---
layout: post
title: Functions with Micro
date:   2017-06-06 09:00:00
---
<br>
As technology evolves so do our programming models. We've gone from monoliths to microservices 
and more recently started to push this separation even further towards functions.

Micro looks to simplify distributed systems development, with [go-micro](https://github.com/micro/go-micro) 
providing a pluggable framework for microservices. Go-micro has historically included a high level [Service](https://godoc.org/github.com/micro/go-micro#Service) 
interface, encapsulating the lower level requirements for microservices. 

Today we're introducing the [Function](https://godoc.org/github.com/micro/go-micro#Function) 
interface, a one time executing Service.

<script src="https://gist.github.com/asim/bfbaf036c90761879dbf6e939e5172e4.js"></script>

### The Inspiration

Ben Firshman open sourced a project last year called [Funker](https://github.com/bfirsh/funker), functions as docker containers. The concept is very 
simple but also extremely clever. 

Functions could quite simply be programs with one method, listening on the network for a request and exiting after 
executing once, leveraging docker swarm services for lifecycle management, discovery, etc.

This sparked the inspiration for including functions as part of go-micro.

### Why Functions?

The function programming model is the evolution of microservices. As our scale requirements increase both technically and organisationally 
there's a need to decouple systems and teams so they can operate independently.

In the past 5 years we've seen the emergence of microservices as a way of dealing with those scaling requirements. The microservices 
architecture pattern is of course nothing new but we've now started to define best practices which help us build better software. 

Functions push us into a new realm of possibility in terms of simplifying distributed systems development and solving software problems. 
Going back to the unix philosophy, "do one thing and do it well", functions truly embody that philosophy even more so than microservices.

While infrastructure helps us build scalable systems, remember that microservices and functions are software architecture patterns 
and programming models, so with that we need tools which help us to write software using those patterns.

### Example Function

Here's a straight forward example of writing a function with go-micro. 

As you can tell it looks almost identical to a service definition. That's because underneath the covers they are exactly the 
same except for one small detail, functions exit after one execution of a handler or subscriber.

Functions give you the same functionality as services, letting you leverage all the existing micro ecosystem tooling.

<script src="https://gist.github.com/asim/7d70cf1160ad1279597f12985fe3fbd5.js"></script>

### Running Functions

As previously stated, functions in micro are one time executing services, the function will exit after completing a request. This then 
poses the question, how do we keep functions running?

There is an abundance of existing tooling out there for process lifecycle management, so feel free to use any of your favourite 
process managers.

However the micro toolkit now includes a convenience tool called [**micro run**](https://m3o.com/docs/run.html).

Here's how to run a function:

```
micro run -r github.com/micro/examples/function
```

The **micro run** command fetches, builds and executes from source. The `-r` flag tells it restart the function on exit. 
It's currently a simple and experimental tool for running micro based microservices and functions. From source to running in one command.

There will be a separate post for the run command once it's more stable.

### Summary

Functions are a natural extension of microservices as the next programming model to help simplify distributed systems development. 
Micro treats functions as a first class citizen.

While functions have been added to go-micro, it does not mean 100% of your software needs to be written with them. It's important 
to understand when monoliths, microservices or functions are appropriate.

Look to see more on integrating micro functions with existing systems and serverless tooling in the near future.

<center><p>...</p></center>
If you want to learn more about the services we offer or microservices, checkout the [website](https://m3o.com) or 
visit [GitHub](https://github.com/micro/micro).

Follow us on [Twitter](https://twitter.com/m3ocloud) or join the [Slack](http://slack.m3o.com) community.


