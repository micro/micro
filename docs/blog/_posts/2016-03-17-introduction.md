---
layout: post
title: Introducing Micro - a microservice ecosystem
date:   2016-03-17 09:00:00
---
<br>
Hello World!

Let's talk about the future of software development.

Change is afoot. We're moving more and more towards a world driven by technology at 
the heart of every business. Maintaining a competitive edge in this day and age is 
becoming difficult. Organisation's ability to execute can grind to a halt as they 
try to scale with inefficient platforms, processes and structure. Decade old technology 
companies have undergone these scaling pains and most have used the same methods 
to overcome those challenges.

It's time to bring the competitive advantages of the most successful companies in the world to 
everyone else. So with that, let's talk microservices, a way to create competitive edge.

###### What are microservices?

Microservices is an software architecture pattern used to breakdown large monolithic applications 
into smaller manageable independent services which communicate via language agnostic protocols and 
each focus on doing one thing well.

Definition of microservices from industry experts.

> Loosely coupled service oriented architecture with a bounded context <br/>
> <sub>Adrian Cockcroft</sub>

<p/>

> An approach to developing a single application as a suite of small services, 
each running in its own process and communicating with lightweight mechanisms <br/>
> <sub>Martin Fowler</sub>

The concept of microservices is not new, this is the reimagining of service orientied architecture 
but with an approach more holistically aligned with unix processes and pipes.

The philosophy of Microservices architecture:

- The services are small - fine-grained as a singular business purpose similar to unix philosophy of "Do one thing and do it well"
- The organization culture should embrace automation of deployment and testing. This eases the burden on management and operations
- The culture and design principles should embrace failure and faults, similar to anti-fragile systems.

###### Why microservices?

As organisations scale both technology and head count it becomes much more difficult to manage monolithic 
code bases. We all became accustomed to the Twitter fail whale for a period of time as they attempted to 
scale their user base and product feature set with a monolithic system. Microservices enabled Twitter to 
decompose their application into smaller services that could be managed separately by many different teams. 
Each team being responsible for a business function composed of many microservices which can be deployed 
independently from other teams.

<div class="text-center">
<img src="{{ site.baseurl }}/blog/images/micro-service-architecture.png" style="width: 100%; height: auto;" />
</div>

We've seen through first hand experience that microservice systems enable faster development cycles, 
improved productivity and superior scalable systems.

Let's talk about some of the benefits:

1. **Easier to scale development** - teams organise around different business requirements and manage their own services.

2. **Easier to understand** - microservices are much smaller, usually 1000 LOC or less.

3. **Easier to deploy new versions of services frequently** - services can be deployed, scaled and managed independently.

4. **Improved fault tolerance and isolation** - separation of concerns minimises the impact of issues in one service from another.

5. **Improved speed of execution** - teams deliver on business requirements faster by developing, deploying and managing microservices independently.

6. **Reusable services and rapid prototyping** - the unix philosophy ingrained in microservices allow you to reuse existing services and 
build entirely new functionality on top much quicker.

###### What is Micro?

Micro is a microservices ecosystem focused on providing products, services and solutions to enable 
innovation in modern software driven enterprises. We plan to be the defacto resource for anything 
microservices related and will look to enable companies to leverage this technology for their 
own businesses. From early stage prototyping all the way through to large scale production 
deployments.

We've seen a fundamental shift coming in the industry. Moore's law is in effect and we're gaining access 
to more and more compute power everyday. Yet we're not able to fully realise this new capacity. Existing 
tools and development practices do not scale in this new era. Developers are not being provided the tools 
to move from monolithic code bases to more efficient design patterns. Most companies inevitably reach a 
point of diminishing returns with monolithic designs and have to undergo massive R&D reengineering efforts. 
Netfix, Twitter, Gilt and Hailo are all prime examples of this. All ended up building their own microservice 
platforms.

Our vision is to provide the fundamental building blocks to make it easier for anyone to adopt microservices.

We're kicking all this off with an open source microservice toolkit, also called <b>[Micro](https://github.com/micro/micro)</b>.
Expect a follow up blog post detailing the toolkit soon.

###### What now?

When you hear the word microservices we want you to think, Micro - the microservices ecosystem. <br/>

If you want to learn more about the services we offer or microservices, checkout the website [micro.mu](https://m3o.com) or 
the github [repo](https://github.com/micro/micro).

Follow us on Twitter at [@MicroHQ](https://twitter.com/m3ocloud) or join the [Slack](https://slack.m3o.com) 
community [here](http://slack.m3o.com).


