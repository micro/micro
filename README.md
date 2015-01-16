# Micro
Home for a microservices protocol and suite of language agnostic tooling

## Premise
Microservices is an architecture pattern used to decompose a single large application in to a smaller suite of services. Generally the goal is to create light weight services of 1000 lines of code or less. Each service alone provides a particular focused solution or set of solutions. These small services can be used as the foundational building blocks in the creation of a larger system.

The concept of microservices is not new, this is the reimagination of service orientied architecture but with an approach more holistically aligned with unix processes and pipes. For those of us with extensive experience in this field we're somewhat biased and feel this is an incredibly beneficial approach to system design at large and developer productivity.

The goal of **Micro** is to try introduce the concept of microservices to the developer community and start a conversation around a protocol of sorts. Hopefully this will allow us to come to a consensus on requirements and begin writing libraries in every language so this architecture pattern can be used by anyone and everyone.

It's still early days but a proof of concept library has been started in Go which can be found at [github.com/asim/go-micro](https://github.com/asim/go-micro)

## Requirements

The foundation of a library enabling microservices is based around the following requirements:

- Server - an ability to define handlers and serve requests 
- Client - an ability to make requests to another service
- Discovery - a mechanism by which to discover other services

## Resources

[A microservices primer by Martin Fowler](http://martinfowler.com/articles/microservices.html)
