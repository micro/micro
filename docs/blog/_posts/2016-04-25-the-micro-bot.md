---
layout: post
title:  The Micro Bot - ChatOps for microservices
date:   2016-04-25 00:00:00
---
<br>
Today I want to talk to you about bots.

###### Bots? Really...

Now I know what you're thinking. There's a lot of hype around bots at the moment. If you're familiary with 
chatterbots you'll know it's not a new concept and in fact goes way back to the days of [Eliza](https://en.wikipedia.org/wiki/ELIZA). 
The renewed fascination with bots has really emerged because we've found more useful functions for them now beyond sheer amusement. They've also 
shined a light on what may become the next dominant interface for interaction beyond apps, the conversational UI.

In the engineering world though, bots are not just for conversational purposes, they can be incredibily useful for operational tasks. 
So much so that most of us techies have become familiar with the term ChatOps. GitHub have been credited with the origins of this term 
since publicising the creation and use of [Hubot](https://hubot.github.com/), a chatbot for managing technical and business 
operation tasks.

Check out this presentation by Jesse Newland on [ChatOpts at GitHub](https://www.youtube.com/watch?v=NST3u-GjjFw) to learn more about it.

Hubot and bots like it have proven to be incredibly useful in technical organisations and become a staple in the world of ops and automation. 
The ability to instruct a bot to perform task through HipChat or Slack which you would otherwise perform manually or cron a script for is 
rather powerful. There's immediate value in the visibility it offers to the entire team. Everyone can see what you're doing and 
what the effects are. 

###### What does this have to do with micro services?

[**Micro**](https://github.com/micro/micro), the microservice toolkit, includes a number of services which provide entry points into your 
running systems. The API, Web Dashboard, CLI, etc. These are all fixed points of entry to interact and observe your microservices environment. 
Over the past few months it became clear that the Bot is another form of entry point to interact and observe and that it should be a first 
class citizen in the Micro world.

So with that...

###### <center>Introducing The Micro Bot<center>
<p align="center">
  <img src="{{ site.baseurl }}/blog/images/micro_bot.png" />
</p>

Let me start by saying, The Micro Bot is a VERY early stage prototype and is currently focused on providing feature parity with the CLI. 
We're not boasting AI based ChatOps here... but maybe someday...

The Micro Bot includes hubot like semantics for scripting as [Commands](#commands) and a way of implementing new [Inputs](#inputs) like Slack and 
Hipchat. It's a rough version 1 but I'm a big believer in shipping as soon as something works and think that by doing so it will open us up to a more rapid 
paced effort to improve the bot. Hopefully with community contributions!

The Bot includes all the CLI commands and Inputs for Slack and HipChat. Our community Bot currently runs in the [demo](http://web.micro.pm) 
environment and is in the Micro Slack right now! Join us [here](http://slack.m3o.com) if you want to check it out.

In the near term we'll look to add more Input plugins like IRC, XMPP and Commands that simplify managing micro services in a running 
environment. If you have ideas for other Inputs or Commands or would like to submit a PR for something please do, contributions are 
more than welcome.
Additional plugins can be found in [github.com/micro/go-plugins/bot](https://github.com/micro/go-plugins/tree/master/bot).

This is really the foundational framework for a programmable bot for the Micro ecosystem. Given the pluggable nature of the entire toolkit it 
only makes sense to provide something similar in bot form.

Let's move on to discussing how Inputs and Commands work.

###### Inputs

Inputs are how the micro bot connects to hipchat, slack, irc, xmpp, etc. We've currently got implementations for 
HipChat and Slack as mentioned above, which seem to cover a significant number of users.


Here's the Input interface.

```
type Input interface {
	// Provide cli flags
	Flags() []cli.Flag
	// Initialise input using cli context
	Init(*cli.Context) error
	// Stream events from the input
	Stream() (Conn, error)
	// Start the input
	Start() error
	// Stop the input
	Stop() error
	// name of the input
	String() string
}
```

The input provides a convenient feature for adding your own command line flags and processing the arguments. 
The Flags() method is used before initialisation and any flags specified will be added to the global flags list. 

After the flags have been parsed, Init() will be called next so that any context for the Input can be initialised. 
Once everything is setup, the Bot will call Start() and then Stream() method to create a connection to the Input.


Here's Conn interface returned by Stream method.


```
type Conn interface {
	Close() error
	Recv(*Event) error
	Send(*Event) error
}
```

The bot will continuosly call Recv() waiting for events. Recv() should essentially be a blocking call otherwise 
we'll end up in a spin loop that will chew up the CPU. Once the Bot has processed the event it will return 
some resulting event using the Send() method.


An Event is the common type sent back and forth between the bot and inputs. It allows us to translate various message types 
of the inputs into a common format. There's currently only a TextEvent type but in the future we have more. 

The bot knows nothing about whether something is from Slack, HipChat or anywhere else. It just knows it's received an event and 
has to do something with it. It's a great way of separating responsibility of the bot from the input.


Here's the Event type.

```
type Event struct {
	Type EventType
	From string
	To   string
	Data []byte
	Meta map[string]interface{}
}
```

###### Commands

Commands are functions that can be executed by the bot. It's that simple. They're stored in a map, keyed by regex, that will be matched 
against text events received from the inputs. If a regex matches the event, the associated function will be executed. The command response 
is then sent back to the input from which the Event originated. If the From field of the originating Event is not blank it will be sent 
as the To field for the response. You can see how this would then allow the bot to directly communicate with someone or something.

The current interface for the Command is fairly straight forward but could potentially change in the future if more complex use cases 
arise. 

The command interface

```
type Command interface {
	// Executes the command with args passed in
	Exec(args ...string) ([]byte, error)
	// Usage of the command
	Usage() string
	// Description of the command
	Description() string
	// Name of the command
	String() string
}
```

Here's an example Echo command

```
// Echo returns the same message
func Echo(ctx *cli.Context) Command {
	usage := "echo [text]"
	desc := "Returns the [text]"

	return NewCommand("echo", usage, desc, func(args ...string) ([]byte, error) {
		if len(args) < 2 {
			return []byte("echo what?"), nil
		}
		return []byte(strings.Join(args[1:], " ")), nil
	})
}
```

And that's the initial building blocks of the bot required to create a conversational interface for Micro services.

###### What else?

It's not enough to just have Inputs and Commands. What if we want to defer some process to a later date? What if 
we want to persist some state in the bots memory? What about real two way dialog rather than just canned responses? 

It must be built!

We're still in the early phases of developing the bot framework and it's an opportunity to contribute to what the 
foundational interface should look like.

The next step is to provide an interface for <i>streams of consciousness</i>. Sounds...abstract. Being a little more serious, 
we need a `Stream` interface or something similar, which overlays `Input.Conn` so that we can write plugins for processing 
all input event streams.

This should potentially allow the ability to use multiple input streams at the same time. So that we may take events from 
one stream, process, and then respond elsewhere.

An example would be receiving a message on Slack, querying a micro service in the platform and sending a summary email.

###### Where does it run?

The micro bot runs in your environment alongside other services. In a way it's just like any other service. It will register 
with service discovery and leverage it to see everything that's running.

<p align="center">
  <img src="{{ site.baseurl }}/blog/images/bot.png" style="width: 100%; height: auto;" />
</p>

###### How do I run it?

Because the bot behaves like any other service you'll need to be running a service discovery mechanism. The default is consul.

Using it with Slack is as simple as

```
micro bot --inputs=slack --slack_token=SLACK_TOKEN
```

And with HipChat

```
micro bot --inputs=hipchat --hipchat_username=XMPP_USERNAME --hipchat_password=XMPP_PASSWORD
```

###### The bot in action

Here's some screengrabs to give you an idea of what it looks like in action. As you can see, it replicates the features of the 
micro CLI. We've got some extra commands like animate and geocode in the Micro Slack just for kicks. They're 
in [github.com/micro/go-plugins](https://github.com/micro/go-plugins) if you want to add them yourself.

<p align="center">
<img src="{{ site.baseurl }}/blog/images/slack.png" style="width: 90%; height: auto;" />
</p>

<p align="center">
<img src="{{ site.baseurl }}/blog/images/hipchat.png" style="width: 90%; height: auto;" />
</p>

###### Adding new Commands

Commands are functions executed by the bot using text based pattern matching, similar to Hubot or any other ChatOps bot.

Here's how to write a simple ping command.

###### 1. Write a Command

Firstly create a command using the NewCommand helper. It's basically a quick start for creating commands. 
You can implement the Command interface yourself too if you like.

```go
import "github.com/micro/micro/bot/command"

func Ping() command.Command {
	usage := "ping"
	description := "Returns pong"

	return command.NewCommand("ping", usage, desc, func(args ...string) ([]byte, error) {
		return []byte("pong"), nil
	})
}
```

###### 2. Register the command

Add the command to the Commands map with a pattern key that can be matched by [golang/regexp.Match](https://golang.org/pkg/regexp/#Match). 

Here we're saying that we'll only respond to the word "ping".

```go
import "github.com/micro/micro/bot/command"

func init() {
	command.Commands["^ping$"] = Ping()
}
```

###### 3. Link the Command

Create a link with an import for your command.

link_command.go:

```go
import _ "path/to/import"
```

Build micro with your command

```go
cd github.com/micro/micro
go build -o micro main.go link_command.go
```

And that's all there is to creating a command.

###### Adding new Inputs

Inputs are plugins for communication e.g Slack, HipChat, XMPP, IRC, SMTP, etc, etc. 

New inputs can be added in the following way.

###### 1. Write an Input

Write an input that satisfies the Input interface.

```go
type Input interface {
	// Provide cli flags
	Flags() []cli.Flag
	// Initialise input using cli context
	Init(*cli.Context) error
	// Stream events from the input
	Stream() (Conn, error)
	// Start the input
	Start() error
	// Stop the input
	Stop() error
	// name of the input
	String() string
}
```

###### 2. Register the input

Add the input to the Inputs map.

```go
import "github.com/micro/micro/bot/input"

func init() {
	input.Inputs["name"] = MyInput
}
```

###### 3. Link the input

Create a link with an import for your input plugin.

link_input.go:

```go
import _ "path/to/import"
```

Build micro with your input

```go
cd github.com/micro/micro
go build -o micro main.go link_input.go
```

Inputs are a little tricker to implement than Commands but that's the gist of it.

###### What next?

Making sense of a microservices world isn't easy. It requires a different set of tools and a focus on observability. 
Monitoring, distributed tracing, structured logging and metrics all play a role but even then it can be difficult. 

Imagine a world in which bots are capable of making sense of distributed systems. Providing 
feedback when we really need it rather than having to stare at dashboards and deal with false alerts. 
You've heard of NoOps right? Well what if it was BotOps? What if you never had to be on-call ever again? 
What if they could be used as 1st point of call or run through a set of procedures to rule out common issues during 
outages. Just some food for thought.

Some outlandish ideas and definitely some ways off. At the very least, look for future integrations into Kubernetes, 
Mesos, etc for managing your services directly from HipChat or Slack and automation of other common tasks.

###### Summary

The bot revolution is upon is. The landscape of infrastructure and automation is changing. We believe bots can play a 
vital role, initially in a classic ChatOps form but longer term achieving much more.

Bots should be treated as first class citizens along side configuration management, command line interfaces and APIs. 
We're doing just that in the Micro ecosystem by including a bot as part of the [Micro toolkit](https://github.com/micro/micro).

It's still early days but looking very promising thus far.

If you want to learn more about the services we offer or microservices, check out the [blog](/), the  website 
[micro.mu](https://m3o.com) or the github [repo](https://github.com/micro/micro).

Follow us on Twitter at [@MicroHQ](https://twitter.com/m3ocloud) or join the [Slack](https://slack.m3o.com) 
community [here](http://slack.m3o.com).

