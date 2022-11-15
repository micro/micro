---
layout: post
title:  "Micro Web Reborn"
author: Asim Aslam
date:   2022-11-15 10:00:00
---
Micro Web was a dashboard built for querying Micro services in a browser. In v3 we largely sidelined it to focus 
on a headless server driven by a CLI and consumed via APIs. Today though we're bringing it back in an all new way. 
Micro Web is Reborn as a pure JS experience.

## Home

The first thing you'll notice is the Home screen looks pretty sparse. We've taken it back to 
basics to focus on just querying services. Everything relating to admin is done through the CLI, the 
web is now focused purely on consumption of the services run by Micro.

<center>
<img src="{{ site.baseurl }}/blog/images/micro-web.png" style="width:100%; max-width: 800px; border: 1px solid; height=auto" />
<p>The home screen</p>
</center>

## Search

The home screen provides a simple search prompt to find services, it will start to render a list as you type to help find what you need. 
If you already know what you're looking for, a query can be made directly in the prompt almost like a command line.

<center>
<img src="{{ site.baseurl }}/blog/images/micro-web-prompt.png" style="width:100%; max-width: 800px; border: 1px solid; height=auto" />
<p>The search prompt</p>
</center>

## Query

For example if I'm looking to query the Twitter timeline for Elon Musk `twitter timeline username=elonmusk` will immediately execute the 
query and render a result like so. Recent queries are displayed under the prompt for conveniently getting back to them at a later date.

<center>
<img src="{{ site.baseurl }}/blog/images/micro-web-text.png" style="width:100%; max-width: 800px; border: 1px solid; height=auto" />
<p>search: twitter timeline username=elonmusk</p>
</center>

## Format

Each service endpoint is given a dynamically generated form and defaults to a text based output with the option of also rendering JSON.

<center>
<img src="{{ site.baseurl }}/blog/images/micro-web-json.png" style="width:100%; max-width: 800px; border: 1px solid; height=auto" />
<p>JSON formatted response</p>
</center>

## Services

Services are listed on a separate page for convenience as a directory or catalog view. It's not always clear you know what you're looking for 
so seeing the index listing is a useful feature. This will link directly through to endpoints for each service.

<center>
<img src="{{ site.baseurl }}/blog/images/micro-web-services.png" style="width:100%; max-width: 800px; border: 1px solid; height=auto" />
<p>The service catalog</p>
</center>

## Future

The goal for Micro Web is to enable consumption of Micro services through a browser or mobile device. The idea being that we can facilitate 
immediate access to query the data or execute API endpoints without the need for dev tools like the CLI or client libraries. In future 
we'll look to extend this by providing more powerful visualisation via widgets and embedding into external web apps and services.

Find the source on [GitHub](https://github.com/micro/micro). Follow us on [Twitter](https://twitter.com/MicroDotDev) for updates.
