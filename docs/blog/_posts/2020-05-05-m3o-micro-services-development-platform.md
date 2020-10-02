---
layout:	post
author: Janos Dobronszki
title:	"Announcing: M3O - A Micro services development platform"
date:	2020-05-05 10:00:00
---

<br>
Today we're announcing [M3O](https://m3o.com) - a cloud native platform for Micro services development. A vastly simpler method of building distributed systems in the Cloud and beyond without having to manage the infrastructure. M3O (pronounced "em-3-oh" and derived from the word `M[icr]o`) is the culmination of many years of experience doing distributed systems development and today we want to shed more light on what we're working on.

Before we do, let's talk about how we got here and what developers have had to deal with over the past decade, as all attempts to focus on software development has been thwarted by the new era of Cloud computing.
<br>

## Cloud native development is just too complex

Let's surmise in one sentence, we think modern application development is just too complicated. Building for the cloud has now gone from spinning up a VM via an API to throwing everything and the kitchen sink at container orchestration platforms and the surrounding dumbfounding ecosystem of tools governed by the [CNCF](https://www.cncf.io/).

Gone are the days of launching a successful company on a shared PHP host, or being content with [restarting ones Ruby server every hour](https://books.google.hu/books?id=ja1KDAAAQBAJ&pg=PA134&lpg=PA134) due to memory leaks, and still becoming the literal behemoth and utility that Twitter became back in 09.

Somehow between the disparate contribution of technologies to the Cloud ecosystem from startups and the thousands of features now promoted by AWS, the ease of use for the developer has drastically suffered. In fact we argue, the cognitive load on the developer is now so bad that we cannot even begin to start building software for the cloud without 3-6 months of wrangling a kubernetes cluster on one of the incumbent cloud providers, while muddling our way through CloudFormation and walls of YAML...

Or perhaps it's a form of decision paralysis due to the explosion of available technologies. In this fast moving industry one often can't help but follow the technologies provided by FAANGs in a "Nobody got fired for using IBM" spirit. Or the thousand open source projects listed on the CNCF landscape which give us FOMO and a headache all at the same time as we figure out if any of its needed or all of it.

## So what now?

We are self aware enough to realise our thinking might run critically close to the [XKCD comic about standards](https://xkcd.com/927/). We are also idealistic, (or crazy/bold/ambitious, it's for the dear reader to decide), enough to believe, we as a tiny startup of misfits can provide a leaner approach to other startups to build services for the Cloud era.

Being a small team, we hope to be the David against the Goliaths and provide a unifying, overarching vision for the whole workflow from local to cloud and all steps in between, while adamantly keeping the Holy Grail of developer happiness in mind.

We did not always want to get into the platform business. In fact we desperately tried to avoid doing it at all costs, knowing how PaaS has played out before us. Unfortunately seeing what the status quo has become and having experienced a better world at companies like Google, Hailo and Monzo, we knew it was on us to do something about it.

There were some times over the past few years where we thought technologies like AWS lambda and Google Cloud functions would push serverless in a direction that would solve for all the developers problems but this was wishful thinking. What we came away with was the realisation that developers don't want to and shouldn't have to manage infrastructure. But functions was the wrong development model.

## A platform for Micro services development

Micro started life as a Go framework for microservices development. To solve the problem of distributed systems development with a developer first focus. It took the approach of solving this problem in the smallest way possible. Now though, we find ourselves evolving and moving on to becoming both a [Runtime](https://m3o.com/blog/2020/05/04/introducing-micro-server.html) and a Platform.

M3O is a platform for Micro services development. A fully managed cloud platform which eliminates the need for developers to touch infrastructure and lets them get back to focusing on product and service development. The best kept secret of every successful technology company is now being opened to the world to use as a simply priced hosted platform.

Our goal is to focus on the following:

- Developer productivity - enabling devs to just focus on building services
- Collaboration and reuse - the platform is built for teams and sharing of services between them
- Everything as a Service - all applications on the platform are built as Micro services
- Velocity of development - we're building an environment that allows you to move at a blistering pace

Judging the industry by our primary focus on simplicity, developer productivity and happiness (concepts we believe to be very correlated), we see constant ebbs and flows. One stride towards a positive direction is cancelled out by the introduction and promotion of overly complex tools, or at least the promotion of overly complex tools to the wrong audience and for the wrong usecase.

Naturally, with millions of us working in this industry, we can't and don't expect everything to go the way we envision - we just enthusiastically wish for it. Our industry is still very young, and we all work hard to make it and the world a better place. That being said, we hope our hard work and decades of collective experience with microservices will result in something that the users will love and put to productive use.

## Who are we?

We are a collective of engineers who have experienced the woes of cloud-native complexity, built platforms before the era of containers and fire fought battles with Kubernetes on multiple cloud providers. We are the every-man, the every-day engineer, who just so happens to now want to do something to combat the complexity that AWS, Google and others have introduced to the world.

Micro is and always was, an opinionated framework and ecosystem. Convention over configuration. Easy bootstrapping with zero dependencies locally. Filling in blanks as demands of scaling and resiliency comes up - by switching out implementations of interfaces with more sophisticated ones - that was always the Micro way.

With M3O we plan to keep this approach, working both ways: moving things from local to prod, or prod to local should be a breeze. Micro is particularly well suited for this as an ecosystem built around the use of pluggable abstractions. One of our main goals is to make handling multi environments as easy as possible - regardless of where and how the services are being run, managing and interacting with them should feel just as native as running them on your laptop or PC.

We try to make it so when someone learns how to use Micro locally, deploying to M3O will not be further away then a CLI command (assuming an account on M3O exists). This reuse of the well known and already useful things that are part of the daily workflow will hopefully provide the easiest way for developers to migrate from local to prod or vice versa.

We make a promise to you. We will never ask you to run or touch Kubernetes. You have our word. Because we also know that pain.

## What is M3O?

The M3O platform initially will be providing [Micro](https://github.com/micro/micro) as a Service. A hosted offering of the open source project which will be billed on a monthly subscription basis. We'll build the platform so you can build the product. Rather than thinking about infrastructure management, we want you to focus on service development. You'll never have to touch infrastructure ever again.

Later on we'll look to introduce collaboration features and value add services so you don't have to build them yourself. Email, sms, payments, user management, etc. We've got you covered.

For existing users who already have their own live setup, perhaps easy to configure custom environments (staging, testing or even per engineer environments) will provide enough value to try us out. Or if you're perfectly happy with your current setup, maybe they using M3O for the next project to save bootstrapping time and cost.

## What next?

We invite the reader to [signup](https://m3o.com) to the beta which we'll be providing access to in the coming weeks. We'll have a community tier for open source software so M3O will be easy to test drive without strings attached. If you're interested to learn more come join us on [slack](https://slack.m3o.com) in the #platform channel.

Thanks for reading.

Cheers

The Micro Team
