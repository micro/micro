---
layout:	post
title:	Deprecating Consul in favour of Etcd
date:	2019-10-04 09:00:00
---
<br>
For over 4 years Consul has served us well as one of the default service discovery systems in Micro. It was 
in fact in the very beginning the default mechanism used for the registry and the only underlying 
dependency required to get started.

Since then the world has moved on and cloud-native technologies have evolved. We've found a number of issues 
at scale related to the way in which we use Consul. This is not a knock on Consul but a reflection on our 
use cases and the need to move on to something different.

For example we binary encode, compress and base64 encode our metadata and service endpoint information 
before storing them as Consul tags because there just wasn't any other way to do so. We were also very 
heavily abusing the distributed properties of Consul which caused a number of issues with raft consensus.

Unfortunately we've found that its now time to move on.

Since 2014 kubernetes has really become a reckoning force in the landscape of container orchestration and the 
base level platform for services. With that etcd became their key-value storage of choice, a distributed key-value 
store built with raft consensus. It has evolve to cater to the scale requirements of kubernetes and has since 
been battle tested in a way few other open source projects have.

Etcd also being a very standard Get/Put/Delete store for binary data means we're easily able to encode and store 
our service metadata with zero issues. It has no opinions about the format of the data being stored.

We've in the past week moved etcd to become one of the default service discovery mechanisms in Micro and will be 
looking to deprecate Consul in the coming weeks. What does this mean? Well we'll be moving consul to our 
community maintained [go-plugins](https://github.com/micro/go-plugins) repository and focusing on supporting 
etcd.

We know a number of users are using Consul and this may cause disruption. This to us is a breaking change on our 
path to v2 and so our next release will be tagged as v2. You can be assured that your v1 releases will continue 
to operate as is but expect that the next release we do is micro v2.0.0.

<center>...</center>
<br>
To learn more check out the [website](https://m3o.com), follow us on [twitter](https://twitter.com/m3ocloud) or 
join the [slack](https://slack.m3o.com) community.

<h6><a href="https://github.com/micro/micro"><i class="fab fa-github fa-2x"></i> Micro</a></h6>
