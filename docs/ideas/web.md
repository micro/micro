# Web Evolution

Micro web must evolve from being something that serves a http proxy to the true vision of web apps as microservices.

## Overview

Micro web is something of an anomaly that expects us to write separate web apps but in reality micro web should 
dynamically generate a web UI for each service and allow us to augment that as needed. In an ideal world we 
are only programming backend services and then creating small UI elements, fragments, widgets to visualise those 
or having no UI at all.

In the event something more complete is needed people should go build using the JAM Stack and speak to services using 
the api or potentially micro web enables the use of grpc-web.

## Design

Let's assume micro web maintains its current home screen icon view like dashboard but each app that you click 
through to visualise the backend service to you in a dynamic way, either by generating a table or rendering 
the fragment element you've including in your service at the endpoint [Service].Web e.g User.Web.

The query model in this case would be fairly rudimentary but an extension of the existing form we have 
already built into the web UI. What this would allow is simple querying of services from the UI for anyone 
but without having to actually build a full fledged UI.

## Inspiration

-  https://walkthechat.com/wechat-mini-programs-simple-introduction/
- https://github.com/gogo/letmegrpc
-  https://github.com/fullstorydev/grpcui
-  https://github.com/grpc/grpc-web
- https://github.com/grpc-ecosystem/awesome-grpc#gui
