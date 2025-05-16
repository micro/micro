---
layout: post
title:  "Real World Micro Services"
author: Asim Aslam
date:   2022-09-27 10:00:00
---
Over the years I've become pretty frustrated by the state of tech and engineering in general. 
One of the biggest issues we face in the industry is the lack of reusability in software. 
GitHub made a major revolutionary change for developers, enabling all of us to reuse libraries, 
and code through reuse rather than writing everything from scratch. Yet it never felt like that 
made it any further than that.

We've got a significant number of projects and services which are hugely abundant on GitHub 
that go beyond libraries, but often it feels like piecing together a lot of bespoke and fragmented 
systems. That was my platform engineering life. I automated as many open source pieces of software 
for companies as possible to avoid paying for SaaS services, and then developers layered systems 
on top of these for the business function. 

Inevitably each time we left a company we'd have to start from scratch. Those business systems written by
the devs were siloed, never to be seen again. In places like Google or Amazon, these systems thrived 
over multiple years to create services that became building blocks with compounding value. And no one
felt the pain more, than the engineers leaving these big orgs. While some of the infrastructure services 
made it out as open source, none of the business logic ever did.

If I had one wish, it's that we'd have a community led platform that we could all use which existed 
outside of these orgs, and a set of reusable services that lived beyond any one stint at a company. 
The platform is hard to do. Without the scale and budget of a Google, it's difficult to offer 
such an experience but I think the building block services might be doable.

So today I'm sharing [Micro Services](https://github.com/micro/services) with all of you. A set 
of reusable real world open source Micro services for everyday use.

[Micro](https://github.com/micro/micro) is an open source API first development platform which I've 
been working on since 2015. It formed out the idea of, what would happen if you built Rails or Spring 
for Go. Since then I went on to build [M3O](https://m3o.com), an API platform powered by Micro 
and the services I'm sharing with you today.

My hope by sharing these services with you is not really to sell anything to you, but to try raise awareness 
for the potential of a model where we all reuse the same services as building blocks for future software. 
By doing so we could create truly open services and systems that are portable across clouds.

See the source on <a href="https://github.com/micro/services">GitHub</a>

Cheers, [Asim Aslam](https://github.com/asim)
