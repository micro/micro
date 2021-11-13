---
layout: post
title:  "Why developers need an AWS alternative"
author: Asim Aslam
date:   2021-11-05 10:00:00
image: "/assets/images/xaas.jpg"
---
<center>
<img src="{{ site.baseurl }}/blog/images/xaas.jpg" style="width: 100%; height: auto;" />
</center>

*Author: Asim Aslam, founder of Micro Services, Inc. (Micro). Micro is building M3O, an open source public cloud platform. An AWS alternative for the next generation of developers. Consume 
popular public APIs as higher level building blocks, all on one platform, for a 10x developer experience.*

AWS was launched over 15 years ago, imagined as an operating system for the internet at a time before the cloud even existed as a concept. It was built to provide on-demand 
access to compute, storage and infrastructure services. What previously might have taken 6 weeks to provision was now being done in minutes. It largely unlocked productivity 
in a way that had not yet been imagined and enabled developers to quickly scale web services as users were coming online.

Yet for all it's worth, AWS has largely maintained this experience and failed to keep up with modern day needs of developers. In 2006, we expected a level of server 
and database management as developers. What we classified as sysadmins (not yet devops) was a skill we'd gladly learn if it meant being able to ship the next shiny 
Rails app. 

Today we're looking for more. We're looking for not just fully managed (which AWS attempts to convince us they are), but an entirely serverless experience. We don't 
want to have to deploy that next database or spin up containers. We don't want to deal with the issues that arise when dealing with the complexities of this new 
fangled Kubernetes. As developers, all we want to focus on is building that next product and leveraging the APIs that let us do that.

In 2021, AWS is slowly being eaten by third party API providers like Algolia, Elastic, CockroachDB, Twilio, Stripe, Sendgrid, Segment and so many many more. We're looking 
for entirely API first experiences in the cloud that don't require us to think about the infrastructure. We're looking for platforms that compliment our modern day 
Jamstack architectures powered by the Netlify's and Vercel's of the world.

AWS now leaves a lot to be desired for the next generation of devs. Can they do anything to fix that?

We personally don't believe so.

## Who is AWS for?

Then if AWS isn't built for developers, who's it for? AWS was never built for developers in the first place, let's be clear about that. AWS was about provisioning 
infrastructure services on which we could then run our software which was still automated by the sysadmins in our companies. How do we know that? Because we were 
heavily invested in AWS in our prior companies.

We contended with the complexity of the bare metal era before AWS and then what came after managed largely by a disarray of hand crafted bash scripts, python libraries and 
eventually configuration management tools like chef and puppet. We escaped just as the DevOps movement took off but continued to witness the extraordinary pains of 
building systems for the cloud as a software engineer.

Yet in all that time, we never once saw developers personally touch CloudFormation, or swim the sea of endless complexity unless they truly had to. No, those 
developers would gladly choose a Heroku long before an AWS, but if you worked at a startup that was scaling, at some point in the lifetime of the company 
you could expect an infrastructure engineer to join and quickly replatform you to AWS.

The truth of the matter is. AWS was built for operators, not the developers.

There's enough people now shouting at the screen pointing to services like AWS Lambda or Fargate talking about it's serverless nature or how it was built for 
developers but I'd argue, this is just AWS pandering to an existing Enterprise audience and checking off boxes. AWS is about building the 80% solution to 
keep existing customers happy, that doesn't mean the actual users in those companies are happy using them.

AWS was the "just good enough" solution for a time in which we had nothing else. The book store that started a cloud computing business even admits to 
being shocked that they had a 7 year head start on everyone else. Had Google gotten their act together and shipped their superior internal tools as 
public services long ago, we'd be having an entirely different conversation.

The fact of the matter is, AWS doesn't understand developers and the harder they try the more complex their offering becomes. As a developer AWS is 
an overwhelming and anxiety inducing experience.

## What do developers need?

What we need is a clean slate approach. As developers we need a new experience for cloud services. We need a new public cloud platform. One that focuses 
entirely on the developer experience. Higher level building blocks for existing public APIs.

Replicating an AWS isn't the answer. VMs (EC2) and file storage (S3) are not the primitives developers need today. We need to start with next 
level building blocks for the next generation of devs. Today we are all about the Jamstack and leveraging third party APIs as the backend. 

We need a public API platform that aggregates the existing market and provides a new clean abstraction layer on top, all through a 
single unified offering. One that simplifies the pricing model rather than requiring a pricing calculator to know what you're spending.

For what it's worth, we thank AWS for getting us to this point, but now it's time to hand the torch to someone else. Someone who understands 
what developers need and provide the next level building blocks for new types of software we'll come to use.

The world is no longer talking about building mobile apps or web services but instead, crypto networks and the metaverse. Your grandparents can barely 
use a mobile phone, are we really expecting AWS and others to help us build the metaverse? 

It's up to developers to build the future and with it decide the kinds of platforms we want to build on. We're now more than ever interested 
in open platforms. Not just in the case of Web3 but more so in regards to "open source eating everything". It's not enough that just the services 
you run are open source, the entire system also needs to be so.

AWS built in an era before GitHub, and the explosive nature of open source, is not. AWS is a silo, and a ship filled with containers of 
teams all building APIs in isolation. Their control plane is not open source, their platform is not open source, their system is not open source. 
AWS is not open source.

We are a generation of developers who are looking for a new platform, one that aligns with our goals, beliefs and mantras and one that is entirely 
based on open source software.

I'm Asim Aslam, the founder of Micro, and we're building [M3O](https://m3o.com), a new open source public cloud platform, an AWS alternative for the next 
generation of developers. Come join me in deciding how, where and what we're going to build the future on.

<center>
See the source on <a href="https://github.com/m3o/m3o">GitHub</a>.
</center>
