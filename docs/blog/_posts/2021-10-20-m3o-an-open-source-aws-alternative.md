---
layout: post
title:  "M3O - An open source AWS alternative"
author: Asim Aslam
date:   2021-10-20 10:00:00
---
<p style="text-align: center;">
  <a href="https://m3o.com"><img src="https://raw.githubusercontent.com/m3o/m3o/main/images/m3o.png" /></a>
  <br>
</p>
<p align="center">M3O is an open source public cloud platform.<br>We are building an AWS alternative for the next generation of developers.</p>

---

## Overview

AWS was a first generation public cloud provider started in 2006. It's infrastructure services and pay as go pricing model made it an incredibly 
compelling choice for a previous generation of developers. But what about the future? 

M3O is an attempt to build a new public cloud platform with higher level building blocks for the next generation of developers. 
M3O is powered by the open source [Micro](https://github.com/micro/micro) platform and programmable real world [Micro Services](https://github.com/micro/services).

## Features

- **🔥 Dev UX** - The developer experience is first priority. A slick new UX for the next generation of developers.
- **☝️ One Token** - Use one Micro API token to fulfill all your API needs. Access multiple public APIs with a single token.
- **⚡ Fast Access** - Using a new API is easy - no need to learn yet another API, it's all the same Micro developer experience.
- **🆓 Free to start** - It's a simple pay as you go model and everything is priced per request. Top up your account and start making calls.
- **🚫 Anti AWS Billing** - Don't get lost in a sea of infinite cloud billing. We show you exactly what you use and don't hide any of the costs.
- **✔️ Open Source Software** - Built on an open source foundation and services which anyone can contribute to with a simple PR.

## Rationale

AWS is a fairly complex beast which makes it hard for new developers to get started. In the past we needed VMs and file storage, but today with the Jamstack 
and other modern development models, the building blocks we're looking for are changing. They're mostly now third party APIs. M3O is looking to 
aggregate all those third party public APIs into a single uniform offering with a slick new dev UX.

## Services

So far there are over 35+ services. Here are some of the highlights:

### Backend

- [**Cache**](https://m3o.com/cache) - Quick access key-value storage
- [**DB**](https://m3o.com/db) - Simple database service
- [**Events**](https://m3o.com/event) - Publish and subscribe to messages
- [**Functions**](https://m3o.com/function) - Serverless compute as a service
- [**User**](https://m3o.com/user) - User management and authentication
- [**File**](https://m3o.com/file) - Store, list, and retrieve text files

### Logistics

- [**Address**](https://m3o.com/address) - Address lookup by postcode
- [**Geocoding**](https://m3o.com/geocoding) - Geocode an address to gps location and the reverse.
- [**Location**](https://m3o.com/location) - Real time GPS location tracking and search
- [**Routing**](https://m3o.com/routing) - Etas, routes and turn by turn directions
- [**IP**](https://m3o.com/ip) - IP to geolocation lookup

### Web

- [**Email**](https://m3o.com/email) - Send emails in a flash
- [**Image**](https://m3o.com/image) - Quickly upload, resize, and convert images
- [**OTP**](https://m3o.com/otp) - One time password generation
- [**QR Codes**](https://m3o.com/qr) - QR code generator
- [**SMS**](https://m3o.com/sms) - Send an SMS message
- [**Weather**](https://m3o.com/weather) - Real time weather forecast

See the full list at [m3o.com/explore](https://m3o.com/explore) or the source at [github.com/micro/services](https://github.com/micro/services).

## Getting Started

- Head to [m3o.com](https://m3o.com) and signup for a free account.
- Generate an API key on the [Settings page](https://m3o.com/settings/keys).
- Browse the APIs on the [Explore page](https://m3o.com/explore).
- Call any API using your token in the `Authorization: Bearer [Token]` header and `https://api.m3o.com/v1/[service]/[endpoint]` url.
- See the [m3o/cli](https://github.com/m3o/m3o/tree/main/cli) for command line usage.

## Learn More

- Checkout the [Examples](examples) to see what you can build
- Follow us on [Twitter](https://twitter.com/m3oservices)
- Join the [Slack](https://slack.m3o.com) community
- Join the [Discord](https://discord.gg/TBR9bRjd6Z) server
- Email [Support](mailto:support@m3o.com) for help

## How it Works

M3O is built on existing public cloud infrastructure using managed kubernetes along with our own [infrastructure automation](https://github.com/m3o/platform) 
and abstraction layer for existing third party public APIs. We host the open source [Micro](https://github.com/micro/micro) project as our base OS and 
use it to power all the [Micro Services](https://github.com/micro/services), which provide simpler building blocks for existing cloud primitives.

### UX

We host our own custom dev UX ([m3o/cloud](https://github.com/m3o/cloud)) on top of the infrastructure stack and a [backend](https://github.com/m3o/backend) 
which acts as the management control plane.

### Services

Developers build and contribute to services in [github.com/micro/services](https://github.com/micro/services), a vendor neutral home. We then automate the 
building and publishing of those services and client libraries. This creates a shared and fully managed platform for everyone to leverage.

### Infrastructure

We primarily use existing open source software, fully managed services and third party public APIs as the backing infrastructure then layer a standard interface 
on top. With all the services on one platform, accessible with one API token, we drastically improve the Dev UX.

## Development

This project is a combination of open source projects and a platform managed by the Micro team. Our goal is to enable any developer to 
contribute to the open source while benefiting from the platform as a shared resource.

### Cloud Hosting

The cloud hosted providers of Micro services:

- [m3o.com](https://m3o.com) - a fully managed offering of micro services

### Open Source

The core cloud OS and services exists in a vendor neutral org

- [micro/micro](https://github.com/micro/micro) - an open source operating system for the cloud
- [micro/services](https://github.com/micro/services) - open source micro services powering m3o.com

### M3O Dev UX

The hosting of Micro services on [m3o.com](https://m3o.com) is powered by the following:

- [m3o/cloud](https://github.com/m3o/cloud) - locally hostable angular based dev UX for the website
- [m3o/platform](https://github.com/m3o/platform) - the infrastructure automation for cloud hosted stack
- [m3o/backend](https://github.com/m3o/backend) - the services which power the m3o.com product backend

## Publish APIs

If you'd like to publish your own APIs on the M3O platform [fill in this form](https://forms.gle/9SQV6DdLNDzSRQ477) and we'll get back to you.

## Contributing

- Write examples - [m3o/examples](https://github.com/m3o/m3o/tree/main/examples) provides a place to showcase things built on top which we'll feature on the website
- Write services - [micro/services](https://github.com/micro/services) are all the services we host and hope for more devs to help
- Write docs - [m3o/docs](https://github.com/m3o/m3o/tree/main/docs) is where our docs will live and we know without great docs this isn't going to work
- Show support - Tweet [@m3oservices](https://twitter.com/m3oservices) and let the world know what we're building so that more people can get involved

## Thanks for reading

Hopefully you buy into what we're talking about and the need for something new in the public cloud space.
If you like what you’re hearing, [Signup for Free](https://m3o.com) or send us some [feedback](mailto:contact@m3o.com). 
Reach out on [slack](https://slack.m3o.com) or [twitter](https://twitter.com/m3oservices) if you have any questions.


