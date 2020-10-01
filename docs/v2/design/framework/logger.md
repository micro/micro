# Logging

Micro need usable logger to be able to write messages about internal state and errors and also provide useful 
logger interface for end-user.

## Overview

Logger must provide minimal interface to write logs at specific levels with specific fields.

## Implemenations

We have 4 implemntations:
* micro internal that writes to console and also to internal in-memory ring buffer
* zap
* logrus
* zerolog

## Design

```go
type Logger interface {
    // Init initialises options
    Init(options ...Option) error
    // The Logger options
    Options() Options
    // Fields set fields to always be logged
    Fields(fields map[string]interface{}) Logger
    // Log writes a log entry
    Log(level Level, v ...interface{})
    // Logf writes a formatted log entry
    Logf(level Level, format string, v ...interface{})
    // String returns the name of logger
    String() string
}
```

Also we have helper functions that automatic uses specified log-level:
* Warn/Warnf
* Error/Errorf
* Debug/Debugf
* Info/Infof
* Fatal/Fatalf
* Trace/Tracef

This is enought for internal micro usage. Additional helper functions implemented via Helper struct and interface inheritance.

```go

type Helper struct {
    Logger
}

func NewHelper(log Logger) *Helper {
    return &Helper{log: log}
}

func (h *Helper) Info(args...interface{}) {}
.....
```

## Benefits

We don't need to implemet helper functions in all all loggers, but internally helper uses only one Log/Logf

## Expected usage

/main.go:

```go
    ctx, cancel := context.WitchCancel(context.Bacground())
    defer cancel()

    ...
    log := zerolog.NewLogger(logger.WithOutput(os.Stdout), logger.WithLevel(logger.DebugLevel))
    logger.NewContext(ctx, log)

    ...
    handler.RegisterHelloHandler(service.Server(), new(handler.Hello))
    ....
```

/handler/hello.go:

```go
    ...
    func (h *Hello) Call(ctx context.Context, req *xx, rsp *yy) error {
        log := logger.FromContext(ctx)
        l := logger.NewHelper(log.Fileds(map[string]interace{}{"reqid":req.Id}))
        ...
        l.Debug("process request")
        ...
        return nil
    }

```
