---
layout:	post
title:	What is Micro? It's just the future of microservices development.
date:	2019-11-18 09:00:00
---

[**Micro**](https://github.com/micro/micro) is an open source project focused on simplifying microservices development. 
It started life as [**go-micro**](https://github.com/micro/go-micro) - a Go framework for microservice development. But even 
before then, go-micro, was a hacked up tiny library created to enable the development of a "kubernetes as a service" 
project way back when in 2014 (see the first commit [here](https://gist.github.com/asim/a035820aec2d8cba5d73b5be12c6e707)).

**Go Micro** was an idea born out of that attempt to build kubernetes as a service, which was written as a set of microservices, 
but ultimately was too early and scrapped not long after being built. What was left behind though was the kernel of 
something else, a handful of packages which if you squinted hard enough looked like the basis for a framework.

# 2014: In the beginning

At the time microservices was a hot topic but tooling was sparse. People spoke of the benefits of this form 
of architecture and development at their organisations but nobody really had the opportunity to open source their tools 
including our team at Hailo.

I had noticed a pattern back then. A developer joined a company for a couple years, helped build a platform 
and set of services for the business, only to leave and then have to go do it all over again at the next company, 
with no ability to carry over the tools of the first. This really frustrated me. Especially 
because if the right tools existed as open source software we wouldn't have to continually go through this process 
and perhaps we'd be focusing on more interesting problems. Not to mention we'd probably save 6-9 months of our lives 
at the very least.

This started to get me thinking about a way for many companies to rally around a single solution. I knew 
a few things though. Every organisation had different skills, different infrastructure preferences and 
adopting new tools was often a big hurdle.

With that in mind my idea was to start with a very lightweight yet opinionated framework for microservices development. 
Knowing how this approach had benefited us at Hailo it felt like it might resonate with other developers too. So over 
the next couple of months I started to work on what eventually formed the initial go-micro framework.

# 2015: Open sourcing go-micro

In early 2015 I decided to open source go-micro. I was deathly afraid of the idea considering I'd never really 
actively publicised a project before and was worried about the 
quality of my code but there wasn't much to lose really. Go Micro felt like something worth sharing.

As with a lot of open source projects, I posted to [Hacker News](https://news.ycombinator.com/item?id=8895794). 
It didn't really get any comments, I can't even remember if it 
was on the front page but what I do remember is hitting 300 stars within a few days! 

I don't have a quick view of the earliest release on github but thankfully Brian Ketelsen, a good friend and strong 
advocate of Micro forked it back then. You can see that code at [github.com/bketelsen/go-micro](https://github.com/bketelsen/go-micro) and from it it's clear there was a few packages outlining a method of microservices communication.

<center>
<a href="https://github.com/bketelsen/go-micro">
  <img src="{{ site.baseurl }}/blog/images/fork.png" style="height: auto; width: 80%; margin: 0" />
</a>
</center>

<br>
Go Micro at the time included a **registry** for service discovery, **server** for RPC and protobuf based request 
handling and a **client** to call those services by name. It even included a key-value storage package but 
we later removed this to focus entirely on communication first (we've recently added it back in).

# Micro: A microservices toolkit

Somewhere in mid-2015 I came to the realisation that a framework was not enough. Once you'd written those services 
there needed to be a way to access them, to serve them, and to consume them by traditional means. This is where 
I began to think about a toolkit.

In a lot of cases we see open source tools which try to solve one problem. State, load balancing, messaging, etc but 
in the case of microservices you really needed a holistic system that would cover all the bases in a seamless way. 
Something that would essentially form the foundations of a platform.

In that [**Micro**](https://github.com/micro/micro) was born. Micro was built as a toolkit to enable the development of 
a microservices platform. It contained a CLI, Web dashboard and API gateway along with a sidecar for non Go based 
applications. That sidecar pattern has now evolved into something called "service mesh" but back then Netflix 
had this thing called [Prana](https://github.com/Netflix/Prana) which is what the Micro sidecar was based on.

Micro and Go Micro were my full focus for the rest of 2015 and took a significant period of time to develop but 
in Autumn of that year a few companies started to use it in production which gave me hope that it may thrive 
in the years to come.

# 2016: Validating the tooling

In 2016 I decided it was time to test the waters once more. To let the world know about Micro and drum up some traction.
I went to Hacker News once more, only this time, things went a bit differently 
[https://news.ycombinator.com/item?id=11327679](https://news.ycombinator.com/item?id=11327679).

Hacker News responded positively and Micro shot to the top of the front page. Here's the original blog post 
for those interested in reading it [https://m3o.com/blog/2016/03/20/micro.html](https://m3o.com/blog/2016/03/20/micro.html).

<center>
<a href="https://m3o.com/blog/2016/03/20/micro.html">
  <img src="{{ site.baseurl }}/blog/images/micro-post.png" style="height: auto; width: 80%; margin: 0" />
</a>
</center>
<br>

It was clear there was something here, that there might be a demand for such a set of tools, and I wanted to pursue 
it full time. Back then I got the opportunity to work with [Sixt](sixt.com) through a corporate sponsorship. This 
allowed me to work full time on Micro and use them as a feedback loop for it's features and development.

I'm incredibly grateful to Sixt for that opportunity and what it allowed Micro to become. Without them it's unclear 
if it would have made it to where it is today. The sponsorship let me continue to iterate on the tools 
as a solo effort for a few years. 3 years in fact. 

And in that time, Micro grew, from a small open source project, to one with a community of 1k+ members, thousands 
of GitHub stars, but more importantly use in the real world in production.

# 2019: The evolution of Micro

Fast forward to the present. Earlier this year I got the opportunity to take Micro from a solo bootstrapped open source 
project and turn it into a venture funded company with the potential to change microservices development on a much larger scale.

We're not ready to reveal all the details just yet but what I will say is it's enabled us to start executing on what many of 
us developers long for. The ability to build, share and collaborate on services in the cloud and beyond, without the 
hassle of managing infrastructure.

## Progress

The progress we've made as a small team in 6 months is pretty astounding. Having committed more times in that period than I had 
done in the entire 4 years of working on Micro alone.

<center>
  <img src="{{ site.baseurl }}/blog/images/commits.png" style="height: auto; width: 80%; margin: 0" />
</center>
<br>

And as you can see here, if GitHub stars are a measure of anything, it reflects in our awareness, popularity and usage. 
We recently passed the 10k star mark on the [go-micro](https://github.com/micro/go-micro) framework and it feels as though we're 
just getting started with what's possible.

<center>
  <img src="{{ site.baseurl }}/blog/images/10k-stars.png" style="height: auto; width: 80%; margin: 0" />
</center>

You can probably tell exactly where we went from 1 person to 2. Based on this progress I'm fairly confident in my previous assumption 
that [go-micro](https://github.com/micro/go-micro) will go on to become the most dominant Go framework and likely surpass Spring 
adoption globally within the next decade.

## Micro as a Runtime

[Micro](https://github.com/micro/micro) has also progressed significantly as we've moved on from just a sparse set of tools to 
something we're now calling **a microservice runtime environment**.

The idea behind this is to reorient the toolkit to be a full fledged 
environment for building microservices. One which provides a programmable abstraction layer for the underlying infrastructure built 
as microservices themselves.

This image is a little old but you'll get the idea. By abstracting away the underlying infrastructure and creating it as a set of 
services that all look the same, run the same, feel the same we end up with a programmable runtime that acts as a foundation for 
all development, whether it be local, in docker or on kubernetes in the cloud.

<center>
  <img src="{{ site.baseurl }}/blog/images/runtime.png" style="height: auto; width: 80%; margin: 0" />
</center>
<br>
We also redefine the boundaries between development and operations in a way that allows each side to focus on their roles without 
the cognitive load of understanding the other side. In the developers case, we no longer have to reason about infrastructure just code.

The feature set is fairly extensive and growing. 

<center>
  <img src="{{ site.baseurl }}/blog/images/feature-set.png" style="height: auto; width: 80%; margin: 0" />
</center>
<br>


## Micro as a Platform

Even still while Micro as a runtime and having a Go framework for development solves a lot of problems, this isn't enough. 
So Micro continues to evolve. It's no longer enough to just simply provide the tools for building microservices, we also need to
provide the environment in which to share and consume them. We at Micro are now building a **global shared microservices platform** for 
developers by developers.

What does that really mean? Well imagine the platform you're given to work on when you join a company or all of the things you 
have to do from an infrastructure perspective just to get up and running. We're going to provide this as a service to everyone.

A fully managed serverless platform for microservices development (that's a mouthful).

## Why?

I've become frustrated with the status quo and the way in which developers are now forced to reason about infrastructure 
and cloud-native complexity. The barrier to entry in just getting started is too high. Building services in the cloud 
should be getting vastly easier, not harder.

Just take a look at the cloud-native landscape...

<center>
<a href="https://landscape.cncf.io">
  <img src="{{ site.baseurl }}/blog/images/cncf.png" style="height: auto; width: 80%; margin: 0" />
</a>
</center>
<br>

Having to reason about this as a developer is horrible. All I want to do is write and ship software but now 
I'm expected to walk some arduous path of containers, container orchestration, docker, kubernetes, service mesh, etc, etc. 
Why can't I just write code and run it?

## Microservices

You're probably thinking. Ok that's great, I buy into this vision. Simpler app development without managing infrastructure 
but what's microservices got to do with this? 

We firmly believe that all forms of development at scale inevitably end up as distributed systems and the pattern 
for that development is now largely known as **microservices**.

Microservices unlock a huge productivity boost in the companies that adopt them and the velocity of their 
development is such that with every new service added their is compounding value in the system built.

I also believe that developers need a platform that enables this form of development for them to thrive. One in which
they do not have to reason about infrastructure and where they are provided the tools that empower them to build software at 
scale without having to worry about operating large scale systems.

One highly controversial example I want to share is from the startup bank [Monzo](https://monzo.com).

Monzo opted to pursue a microservices architecture from day 1. Knowing there were initial operational tradeoffs to this 
approach but with an insight from their time at Hailo, they knew if the company succeeded on the product side they'd 
need a scalable platform to help them grow and move fast.

This led to the creation of a platform that is now host to 1500 services. This might sound hard to reason about, but 
a shared platform where every developer has the ability to consume and reuse existing services is a fundamentally powerful thing.

Not only that, but when the platform is managed for you, developers can get back to focusing on what's really important. The 
product and the business.

<center>
<blockquote class="twitter-tweet"><p lang="en" dir="ltr">1500 microservices at <a href="https://twitter.com/monzo?ref_src=twsrc%5Etfw">@monzo</a>; every line is an enforced network rule allowing traffic <a href="https://t.co/2r2y9f6LYO">pic.twitter.com/2r2y9f6LYO</a></p>&mdash; Jack Kleeman (@JackKleeman) <a href="https://twitter.com/JackKleeman/status/1190354757308862468?ref_src=twsrc%5Etfw">November 1, 2019</a></blockquote> <script async src="https://platform.twitter.com/widgets.js" charset="utf-8"></script>
</center>
<br>

# Our solution

This form of development has largely been siloed at large tech cos capable of building such systems. But what 
if this was available to every developer as a shared system outside of those large orgs. What if we were able to 
collaborate across org and across teams. What would the velocity of our development as an industry look like as a whole? 

I would argue that all technology would advance faster than it ever has done in the decades that have come before us.
We would finally capture the true potential of the internet.

GitHub was a prime example of this collaboration and innovation in open source, massively reducing the pain of hosting 
source code and creating an environment for reusing code. However there's just one but, this source code largely sits 
at rest on their platform.

What if instead of just sharing code and running it in silos, we shared an environment for software development, one 
in which we could collaborate on services, reusing each others running applications where necessary and focusing 
on solving higher order problems.

It would have it's own quirks and challenges but the opportunities such a platform presents is immense. And something 
we want to explore at [Micro](https://m3o.com), the company.

So that's what we're setting out to do really. To build a global shared services platform for developers by developers. 
Where the pains of cloud, kubernetes and everything else will no longer be felt. An environment 
in which we can build, share and collaborate on micro services based on the [**go-micro**](https://github.com/micro/go-micro) framework.

# Closing

The future of Micro is one which involves rapidly reducing the friction for developers in harnessing the power of 
the cloud and to empower them to build microservices from anywhere, with anyone.

If this sounds interesting to you, come join our community on [**slack**](https://m3o.com/slack), kick the tyres 
on the [**go-micro**](https://github.com/micro/go-micro) framework or come help us make it a reality. We're hiring, 
just drop us an email at <a href="mailto:hello@micro.mu"><b>hello@micro.mu</b></a>.

Cheers
<br>
Asim

<h6><a href="https://m3o.com">Micro</a></h6>
