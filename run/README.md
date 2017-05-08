# micro run

The **micro run** command manages the lifecycle of a microservice. It fetches the source, builds a binary and executes it. 
It's a simple tool which can be used for local development. If no arguments are specified micro run operates as a service 
which can manage other services.

## Overview

```
# fetch, build, execute
micro run github.com/service/foo
# run service manager
micro run
# defer to service manager
micro run -x github.com/service/foo
# restart on death
micro run -r github.com/service/foo
# update source on fetch
micro run -u github.com/service/foo
# kill a service
micro run -k=uuid github.com/service/foo
```

## Usage

```
NAME:
   micro run - Run the micro runtime

USAGE:
   micro run [command options] [arguments...]

OPTIONS:
   -r	Restart if dies. Default: false
   -u	Update the source. Default: false
   -x	Defer run to service. Default: false
   -k	Kill a service with uuid
   
```
