# micro run

The **micro run** command manages the lifecycle of a microservice. It fetches the source, builds a binary and executes it. 
It's a simple tool which can be used for local development. If no arguments are specified micro run operates as a service 
which can manage other services.

Note: The default runtime (Go) requires the Go binary in PATH and GOPATH to be set.

## Overview

Run
```
micro run github.com/service/foo
```

Status
```
micro run -s github.com/service/foo
```

Kill
```
micro run -k github.com/service/foo
```

Run service manager
```
micro run
```

Defer run to service manager
```
micro run -x github.com/service/foo
```

Run and restart on death
```
micro run -r github.com/service/foo
```

Run and update source on fetch
```
micro run -u github.com/service/foo
```

## Usage

```
NAME:
   micro run - Run the micro runtime

USAGE:
   micro run [command options] [arguments...]

OPTIONS:
   -k	Kill service
   -r	Restart if dies. Default: false
   -u	Update the source. Default: false
   -x	Defer run to service. Default: false
   -s	Get service status
   
```

## TODO

- [ ] Accept args and env vars to service 
- [ ] Add Service interface to [go-run](https://github.com/micro/go-run)
- [ ] Support configurable runtimes beyond Go
- [ ] Rebuild with plugins
- [ ] Daemonization?
- [ ] Watch memory consumption and kill?
- [ ] Chroot the process?
