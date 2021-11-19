# Voice

Voice is the next major client of the platform

## Overview

The micro runtime currently includes an api, slack bot, cli, web dashboard and gRPC proxy. 
These are entrypoints into the system, a way of converting different forms of client 
interaction into RPC and service calls. Voice is the next frontier.

## Design

Micro Voice serves as the next client which takes an audio stream and converts it into 
an executable RPC call or event. We do this through the use of speech to text and 
then potentially converting back any response into speech.

We imagine voice like the bot to connect to existing systems such as Alexa, Google Home 
or whatever else is developer in the voice ecosystem.
