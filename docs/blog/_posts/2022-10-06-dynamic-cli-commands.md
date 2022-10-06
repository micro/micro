---
layout: post
title:  "Dynamic CLI commands"
author: Asim Aslam
date:   2022-10-06 10:00:00
---
One of the things I love the most about Micro is dynamic CLI commands. This is a feature that showed up in version 3.0. 
Micro started out with a fixed set of commands like most command line tools. The API however knows how to map HTTP requests
to service dynamically. We wondered, what if you could do the same for the CLI. It turns out, you can. Here's an example.

<center>
<img src="{{ site.baseurl }}/blog/images/dynamic-cli.png" style="width:100%; max-width: 600px; height=auto" />
</center>
<br>

To understand how this works we need to first understand a little bit about the internals of Micro. Every service is 
registered in a central registry including it's name, endpoints and request/response format. This enables us to inspect 
the registry for the details we need.

Calling `micro get service weather` we can see what this looks like.

```
service  weather

Endpoint: Weather.Now

Request: {
        location string
}

Response: {
        location string
        region string
        country string
        latitude float64
        longitude float64
        timezone string
        local_time string
        temp_c float64
        temp_f float64
        feels_like_c float64
        feels_like_f float64
        humidity int32
        cloud int32
        daytime bool
        condition string
        icon_url string
        wind_mph float64
        wind_kph float64
        wind_direction string
        wind_degree int32
}
```

As you can see, the registered data includes an endpoint called `Weather.Now` and the Request/Response which tells us what 
arguments we need to pass in and what comes back. Using this we can construct a CLI command parser to dynamically map a 
request.

In our case the CLI command format is `micro [service] [endpoint] [args]` e.g `micro weather now --location=london`.

We can actually map any command like this. For example, let's say we want to do something basic like get Marc Andreessen's tweets.

<center>
<img src="{{ site.baseurl }}/blog/images/twitter.png" style="width:100%; max-width: 600px; height=auto" />
</center>
<br>

I can't emphasize enough how powerful this has been since it's introduction. The ability to programmatically test and use 
services and APIs from the CLI in this way enables all sorts of scripting functionality but it has also just been a far more 
intuitive way to use services from the command line.

One of the lagging features which I'd love to address in future is unix style piping of responses from one service as requests 
into the next. The unix pipe provides composition like no other, hopefuly at some point we'll be able to replicate this 
as part of the Micro CLI.

Try out [Micro](https://github.com/micro/micro) for yourself!
