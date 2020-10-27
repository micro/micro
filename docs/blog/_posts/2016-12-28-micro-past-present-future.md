---
layout: post
title: Micro - Past, Present and Future
date:   2016-12-28 09:00:00
---
<br>
2016 has been a heck of a year for Micro. In this post we'll recap on the past year and where we're going but to begin with let's talk 
about Micro's origin.

### The origin of Micro

[Micro](https://github.com/micro/micro) or more specifically [go-micro](https://github.com/micro/go-micro) began as a much needed platform 
library for an infrastructure based side project back in late 2014. At the time I was working at Hailo where we had built a microservice platform 
called H2 for our own needs.

Starting on this new project I knew I wanted to use Go and a microservice architecture, but it became clear 
very quickly there wasn't anything quite like the Hailo H2 platform out there. So the painstaking work began to recreate it.

The side project unfortunately never really went anywhere but a few months later I had an early version of go-micro and decided it should 
at least be open sourced on github. The initial commit to github can be found [here](https://github.com/micro/go-micro/commit/8e55cde513cbec9f3e3aaf0ed7beb637510fc9b2). 

<center>
  <img style="width: 100%;" src="{{ site.baseurl }}/blog/images/first-commit.png" />
</center>

Back then go-micro had a fairly simple and straightforward purpose. An RPC framework for microservices. An attempt to address the core 
requirements of microservices and simplify distributed systems development.

Having helped build and leverage a microservice platform at Hailo I knew how powerful it could be in software development and this really 
sparked my own desire to try build something similar as an open source project.

### The evolution of Micro

Go-micro the pluggable RPC framework evolved into [Micro](https://github.com/micro/micro), a microservice toolkit. The toolkit began life 
as an idea to develop a microservice protocol or a specification, as can be seen from the initial commits. [[1]](https://github.com/micro/micro/commit/6259b301a8f2bc8e143910fd320caf8071426504), [[2]](https://github.com/micro/micro/commit/d0610303bb44561064fb450c19264187d06d82f2), [[3]](https://github.com/micro/micro/commit/6876df2caf028b6dbff7d3a5b672d6a0aa7e6ef2)

For whatever reason I lost interest in the idea of a protocol and started to rethink what was actually important. Microservices being a 
architecture pattern, it became clear that the tooling should really guide that architecture and simplify microservices in general.

So appeared the API Gateway and CLI. The commit can be found [here](https://github.com/micro/micro/commit/b2ef68df31644442ad2de6398caf661328542a4f).

It all progressed rapidly from there and became something more featureful.

<center>
  <img  src="{{ site.baseurl }}/blog/images/micro-diag.png" />
  <br>
</center>

I'll spare you the details of how things unfolded throughout 2015 but Micro turned into a full fledged toolkit for microservices. The goal, to provide the fundamental requirements for building and managing microservices. 

### 2016

Near the end of 2015, Micro got noticed by a few companies and even saw production use before the end of the year. That was pretty awesome!

I forgot to mention, a few months prior I had actually quit my job at Hailo to pursue working on Micro full time. Knowing I was burning 
savings as each month went by, I worked tirelessly to address the core requirements and drive initial adoption.

###### Sponsorship

A colleague and friend who helped build the platform at Hailo took notice. He had also moved on and was pursuing the rather 
ambitious goal of transforming a century old enterprise into a modern technology driven company. 

As he began to lay the groundwork for his own project he was also struck by the same issue, a lack of Go based open source tooling for microservices like 
what we had built at Hailo.

This started the discussion about where Micro was headed and how we could work together to create a sustainable future for the project 
which included allowing me the opportunity to work on it full time.

<center>
  <a href="{{ site.baseurl }}/2016/04/25/announcing-sixt-sponsorship.html">
    <img style="width: 50%;" src="{{ site.baseurl }}/blog/images/micro_sixt.png" />
  </a>
  <br>
</center>

In early 2016, [Sixt](https://en.wikipedia.org/wiki/Sixt) the German car rental enterprise agreed to sponsor the open source development of Micro. Read the announcement [here]({{ site.baseurl }}/2016/04/25/announcing-sixt-sponsorship.html). 


[Boyan Dimitrov](https://twitter.com/nathariel), platform director at Sixt and the friend I mentioned, was really the driving force behind the 
sponsorship and it's early adoption in the enterprise.

Sustainable open source development is difficult without sponsorship or other forms of funding. The relationship with Sixt has really been fundamental 
in the development, adoption and success of Micro.

###### The Blog

In the world of software, there's a commonly uttered passive aggressive phrase known as, "read the code". Author's of projects intimately knowing 
their code sometimes believe others should just as easily understand it.

Code on its own though is rarely enough.

Not only knowing this but also having the desire to raise awareness for Micro and touch on higher level ideas I decided it would be important to start a blog.

<center>
<blockquote class="twitter-tweet" data-lang="en"><p lang="en" dir="ltr">The follow up blog post detailing Micro - a microservice toolkit. <a href="https://t.co/TEt8MSgBtL">https://t.co/TEt8MSgBtL</a> <a href="https://twitter.com/hashtag/microservices?src=hash">#microservices</a> <a href="https://twitter.com/hashtag/golang?src=hash">#golang</a></p>&mdash; Micro (@MicroHQ) <a href="https://twitter.com/MicroHQ/status/711886548527157249">March 21, 2016</a></blockquote>
<script async src="//platform.twitter.com/widgets.js" charset="utf-8"></script>
<br>
</center>

The micro blog has served as a place to explore the toolkit and microservices in depth while also delving into other aspects such as the requirements for company adoption and it's impact. Much more could be written about the non-technical requirements of microservices in organisations.

<center>
<blockquote class="twitter-tweet" data-lang="en"><p lang="en" dir="ltr">Micro architecture &amp; design patterns for microservices <a href="https://t.co/cqV226ndZi">https://t.co/cqV226ndZi</a> <a href="https://twitter.com/hashtag/microservices?src=hash">#microservices</a> <a href="https://twitter.com/hashtag/golang?src=hash">#golang</a></p>&mdash; Micro (@MicroHQ) <a href="https://twitter.com/MicroHQ/status/722004109495267328">April 18, 2016</a></blockquote>
<script async src="//platform.twitter.com/widgets.js" charset="utf-8"></script>

<blockquote class="twitter-tweet" data-lang="en"><p lang="en" dir="ltr">My post on why companies adopt microservices is getting decent reads on Medium. Think I need to move the blog there. <a href="https://t.co/eP2QqrYHhY">https://t.co/eP2QqrYHhY</a></p>&mdash; Asim Aslam (@chuhnk) <a href="https://twitter.com/chuhnk/status/750738451931226112">July 6, 2016</a></blockquote>
<script async src="//platform.twitter.com/widgets.js" charset="utf-8"></script>
</center>

The blog has largely been a success, a great way to clarify my own thinking on technical topics and the feedback mostly positive. One of the future 
goals will be to introduce guest writers or topics of interest beyond microservices. 

I also started a publication on [Medium](https://m3o.com/blog) which has helped increase readership.


###### Plugins

Go-micro being self-coined a pluggable framework was all good and well but it's not all that valuable without the plugins. It took me a while to realise this.

Go-micro started out with [consul](https://www.consul.io/) for service discovery and point-to-point http for messaging. Over time new plugins 
were introduced and even a few contributed by the community.

Today there's ample choice with most having been battle tested in production settings. This allows anyone to confidently architect the stack they 
want with minimal code changes.

Here are the service discovery plugins so far:
<center>
  <img style="width: 100%;" src="{{ site.baseurl }}/blog/images/registry_plugins.png" />
  <br>
</center>

And the message broker plugins:
<center>
  <img style="width: 100%;" src="{{ site.baseurl }}/blog/images/broker_plugins.png" />
  <br>
</center>

Every aspect of go-micro is pluggable and the Micro toolkit also boasts pluggability too. Find all the plugins at [github.com/micro/go-plugins](https://github.com/micro/go-plugins).

###### Golang UK Conf

In August, I had the awesome opportunity to speak at the [Golang UK Conf](http://golanguk.com/) about 
[Simplifying Microservices with Micro](https://www.youtube.com/watch?v=xspaDovwk34). It was the first time I had personally spoken 
at a conference and only my third ever speaking engagement.

The talk was an opportunity to raise awareness of the Micro toolkit, explain it's purpose at a high level and dig into the tooling itself 
in a bit of detail.

<center>
  <a href="https://www.youtube.com/watch?v=xspaDovwk34">
    <img style="width: 70%;" src="{{ site.baseurl }}/blog/images/talk.png" />
  </a>
  <br>
</center>

It's safe to say I was pretty nervous but it was good validation for all the hard work that had gone into Micro and a few people approached 
afterwards with positive comments and questions which I'm grateful for.

###### Community

A huge aspect of Micro which must be discussed is the community. The Micro community exists on Slack and is now over 500 members strong.

Open source is driven by developers and without their blessing projects rarely thrive. As I mentioned before, code on it's own is never 
enough. Building an online community around an open source project is what helps drive growth and adoption. It's also the place 
where the most learnings occur and are transferred amongst users of the project.

The Micro Slack has become a place where companies or individual developers discuss their problems, requirements, etc and can receive immediate 
feedback.

I want to thank the members of the community for believing in the project and taking the time to contribute back.

Join us at [slack.m3o.com](http://slack.m3o.com/)

With 2016 in the books, let's look to the future.

### 2017

Micro/go-micro have stabilised over the past year and have been proven out in a production setting. They serve a specific function at 
the core of the overall vision for Micro. The foundational requirements for building microservices.

While the toolkit and framework will continue to evolve, we want to stick with the unix and microservices philosophies by keeping things 
simple and logically separate out any further features to separate tools.

###### Micro OS

At a certain scale distributed systems pose a challenge to manage and make sense of. Existing tools built for fewer applications or a different 
era of software development prove to be less relevant in this new world of microservices. 

A single request may fan out to dozens if not hundreds of microservices, in the case of failure it can be difficult to track down where an issue 
occurred. Multiple copies of each service may be running globally with rescheduling occuring at any time making old host based monitoring 
tough. And with data being distributed to support the scale, transactions are likely not supported by the database but we still need some 
form of coordination.

There are a number of tools emerging to solve many of these problems but nothing offered as a cohesive whole. 
[Micro OS](https://github.com/micro/os) is a microservice operating system which attempts to address all these requirements.

<center>
  <img  src="{{ site.baseurl }}/blog/images/os-diag.png" />
  <br>
</center>

Micro OS is currently a work in progress but based on previously well executed models for providing an entire platform runtime for 
microservices. It will include service discovery, global load balacing, distributed tracing, dynamic configuration, distributed coordination
and much more.

###### Micro Edge

My personal interests in computer networks and distributed systems spans over a decade, going back to my studies at university. There I 
learned about the internet architecture itself and became fascinated with large scale systems.

<center>
<blockquote class="twitter-tweet" data-lang="en"><p lang="en" dir="ltr">My technology interest over time:<br>2002-2006 internet<br>2007-2011 c10k<br>2011-2013 google<br>2013-2016 microservices<br>2016-201x edge compute</p>&mdash; Asim Aslam (@chuhnk) <a href="https://twitter.com/chuhnk/status/808669470109683712">December 13, 2016</a></blockquote>
<script async src="//platform.twitter.com/widgets.js" charset="utf-8"></script>
</center>

Over time as I saw the evolution of those large scale systems and how most problems had been solved by the likes Google, 
I became drawn to new problems pushing beyond the datacenter and driving us to even greater scaling requirements.

In the future we'll find intelligent systems almost everywhere with the need to process large amounts of data and act offline or at a rate 
which is faster than what's possible if we sent requests back to a centralised cloud provider.

Edge computing looks to solve this problem by pushing applications, data processing and services closer to the logical extremes of a network. It is not 
a new concept but something which is becoming more relevant as our needs evolve.

In the second half of 2017 I will begin work on Micro Edge, a set of tools to simplify development of edge computing services. Reach out to learn more.

### Highlights

It's been a heck of a ride so far with a lot more to come in 2017. Stay tuned! And with that said, I'll leave you with a few of the greatest hits from the 2016.

Happy new year!

- [Micro - a microservice toolkit](https://m3o.com/blog/micro-a-microservices-toolkit-c403145b65c1)
- [Micro on NATS - microservices with messaging](https://m3o.com/blog/micro-on-nats-microservices-with-messaging-2dcc248fb5b9)
- [Micro architecture & design patterns for microservices](https://m3o.com/blog/micro-architecture-design-patterns-for-microservices-37f4b9049ad3)
- [Building Resilient and Fault Tolerant Applications with Micro](https://m3o.com/blog/building-resilient-and-fault-tolerant-applications-with-micro-53b454a8e8eb)
- [Why Companies Adopt Microservices and How They Succeed](https://medium.com/@asimaslam/why-companies-adopt-microservices-and-how-they-succeed-2ad32f39c65a)
- [Simplifying Microservices with Micro](https://www.youtube.com/watch?v=xspaDovwk34) [Video]  [Golang UK Conf 2016]

If you want to learn more about the services we offer or microservices, checkout the website [micro.mu](https://m3o.com) or 
the github [repo](https://github.com/micro/micro).

Follow us on Twitter at [@MicroHQ](https://twitter.com/m3ocloud) or join the [Slack](https://slack.m3o.com) 
community [here](http://slack.m3o.com).


