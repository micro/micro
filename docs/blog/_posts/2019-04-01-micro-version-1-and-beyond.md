---
layout:	post
title:	Micro 1.0.0 release and beyond
date:	2019-04-01 09:00:00
---
<br>
Over the past 4 years we've focused on creating the simplest experience for microservice development. To do this 
we built a strongly opinionated open source framework called [**Go Micro**](https://github.com/micro/go-micro) and 
[**Micro**](https://github.com/micro/micro), a microservice toolkit built to explore, query and 
interact with those services via an API Gateway, CLI, Slack and Web Dashboard.

<img src="https://m3o.com/docs/images/go-micro.svg" style="max-width: 100%; margin: 0;" />

## Version 1.0.0

Last month we released **version 1.0.0** of both of these tools. This signifies a huge moment for Micro and the community. We've been 
running in production since 2015 and have become vital for companies like our sponsor Sixt, the german car rental enterprise, who are 
running hundreds of microservices in production.

Micro enables teams to scale microservices development while abstracting away the complexity of distributed systems and cloud-native infrastructure. 
It provides a pluggable and runtime agnostic architecture with sane zero dependency defaults.

<center>
<img src="https://m3o.com/micro-diag.svg" style="max-width: 100%; margin: 0;" />
</center>

<br>
We've considered Micro production ready for a long time but the release of 1.0.0 solidifies the maturity and stability of our tooling. And 
we believe it's the right time for everyone to adopt Micro as the defacto standard for microservice development.

## Usage

Micro has largely grown organically. We've not yet actively engaged in speaking at conferences, meetups or any other form of outreach. Instead 
we focused on solving a real problem and it's shown in the numbers.

<center>
<img src="{{ site.baseurl }}/blog/images/stars.png" style="max-width: 75%; margin: 0;" />
</center>
<br>
It's difficult to track active usage of a library or framework but what's clear from all we can see is that Micro has really resonated with 
the developer community who want a simpler path to adopting microservices with Go.

## Beyond 1.0

The announcement of version 1.0.0 is not just a marker for maturity and stability to run Micro in production but it also signals that this release version 
will not suffer any further breaking API changes. This now also allows us to take stock of all the learnings of Micro's usage of the past 4 
years, how technology has evolved in the industry and what version 2 might start to look like.

When we started, kubernetes was just in it's infancy and gRPC had only recently been released. We're seeing these trends along with service mesh 
and much more. 

Because Micro is pluggable we've always been able to adapt to the changing needs of developers while continuing to provide 
simpler abstractions on top for microservice development. With version 2.0 we have the ability to create an even more frictionless and streamlined
experience.

Some of these ideas will revolve around using gRPC by default, allowing a drop-in experience on kubernetes and potentially a runtime 
for those who don't want to deal with the complexity of cloud-native systems or any dependency management at all.

We'll also be thinking about how to move beyond Go to support multiple languages.

## Collaboration

Slack has served us well for realtime collaboration but we need a medium aligned with open source to push this much further, to be far more inclusive 
and to provide a historic record for newcomers to explore easily.

We're going to work with the community by using GitHub to create an open source location to share ideas, discussion and the roadmap for the 
[**Development**](https://github.com/micro/development) of features for 2.0 and beyond.

To all those interested in contributing and collaboration, create an issue for feature requests, a pull request to share design ideas and we'll work 
together to shape the roadmap.

## Thanks

I want to finish by saying thank you to the Micro community and all who've used or supported us over the past 4 years. It's been a long hard 
but incredibly rewarding journey with so much more left to do. Without the community Micro would not be where it is today. We're 1.6k+ members 
strong in Slack with thousands more across other forums.

<center>
<blockquote class="twitter-tweet" data-cards="hidden" data-lang="en"><p lang="en" dir="ltr">Today I released v1 of micro and go-micro. 4 years of hard work. Thanks to all that supported me along the way. <a href="https://t.co/blI1pJ3hBl">https://t.co/blI1pJ3hBl</a></p>&mdash; Asim Aslam (@chuhnk) <a href="https://twitter.com/chuhnk/status/1102992210088378369?ref_src=twsrc%5Etfw">March 5, 2019</a></blockquote>
<script async src="https://platform.twitter.com/widgets.js" charset="utf-8"></script>
</center>
<br>
Thank you for all your support and contributions. We hope we can return the favour by providing everyone the most inclusive and collaborative 
place for all things microservices.

<center>...</center>
<br>
Micro is the simplest way to build microservices. If you're thinking about microservice development we want to help enable you on that journey. 
To learn more check out the [website](https://m3o.com), follow us on [twitter](https://twitter.com/m3ocloud) or 
join the [slack](https://slack.m3o.com) community.

<h6><a href="https://github.com/micro/micro"><i class="fab fa-github fa-2x"></i> Micro</a></h6>
