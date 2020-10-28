---
layout:	post
author: Janos Dobronszki
title:	Introducing Micro Server
date:	2020-05-04 09:00:00
---
<br>
In 2015, `go-micro`, a Go microservices framework was announced. Today we introduce the `micro server`, which builds on top of `go-micro`, and enables you to run and manage microservices with ease, both locally and across different environments.

Whether you are running the `micro server` locally with zero dependencies (using memory or files), or on k8s (with highly available distributed systems and other third party tools), the `micro server` and the micro cli should provide a straightforward and runtime agnostic experience that is the same across any environment.


For those who prefer action to words, the [getting started guide](https://m3o.com/docs/getting-started.html)! can give a taste of what the `micro server` is about or visit the project on [github](https://github.com/micro/micro).

<br>
<div style="text-align: center; width: 100%;">
  <img src="https://m3o.com/images/runtime10.svg" />
</div>
<br>

## Combining a decade of research

At Google in 2011, Asim Aslam, the creator of Micro, experienced what it was like to build systems at Google scale. Google operated a microservices architecture before most ever knew what the term meant. They did this because of an organisational and technical need to build systems that could be developed and managed independently.

In 2013 a few of us worked at a company called Hailo - a European Uber competitor - one of the few companies that followed Neflix's model to build an organisational architecture based on microservices. This enabled us to become - at the time - to be incredibly productive and arguably the most successful taxi app in Europe.

Years later, ex-Hailo members brought the fruits of years of microservices R&D done at Hailo to companies like Monzo, Sixt and many others. We at the Micro team are working on bringing these ideas to the Open Source community. The source of our passion is both the fact that microservices enable companies to be successful and also that we have collectively seen many dozens of companies building similar systems from scratch. Some successfully, some not so successfully.

## Complexity in the age of cloud computing

We primarily exist to ensure company and developer success - the advantages of microservices became clear to the industry in the past years, but the tooling around it was and is still in infancy. Breaking up the monolith to different processes comes with an increase in [accidental complexity](https://en.wikipedia.org/wiki/No_Silver_Bullet), which, without appropriate tooling can cause non-negligible amounts of pain. Micro's mission is to make this pain go away, and unlock a whole new era of computing where services play as nicely together as functions in a programming language do.

We also aim to help working with different environments (ie local/custom envs/prod) and tools (different databases, orchestration systems etc.), so similar concepts and interfaces can be reused across them.

There is no question that most local and production setups vastly differ.
The differences dictated by resiliency and scaling requirements and available computing resources creates a disconnect between the different steps during the lifecycle of a code change - from the moment of its birth locally to different environments and ending up in production.

A simple file backed persistent on a single node is "not enough" in production settings, but similarly you might not want to install and maintain kubernetes or a different heavy container based solution locally.

This is one of the many areas of modern microservices based workflows that Micro aims to simplify, and our current focus with the `micro server` release. There are many more concepts we plan to explore, so if you are interested, stay tuned for further developments.

## To end

If you're interested in the future of cloud-native development, microservices and distributed systems come join us in the community. You can find us in our [Slack channel](https://slack.m3o.com) or on [Discourse](https://community.micro.mu/). We can help you if for any reason you get stuck with our [docs](https://m3o.com/docs) or the [getting started guide](https://m3o.com/docs/getting-started.html).

Also check out the project on [github](https://github.com/micro/micro) and give us a ⭐ if you like it.

Cheers

From the team at Micro
