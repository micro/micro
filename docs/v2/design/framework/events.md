# Events

Micro events is a timeseries database for event streaming

## Overview

Micro has the concepts of key-value storage through the Store interface and asynchronous messaging 
through the Broker interface. These things individually are useful but often we're looking for 
something more. Essentially a timeseries based event stream of data thats persistent.

## Design

We always start with:

- Go Micro interface
  * Zero dep implementation
  * Distributed system equivalent
  * Service implementation

Our goal is to define a go-micro/event package with Event interface that supports reading and writing 
events and being able to playback from a specific Offset which is a timestamp. The first attempt 
at this was file based style [sync/event](https://github.com/micro/go-micro/blob/master/sync/event/event.go).

```
// Event provides a distributed log interface
type Event interface {
	// Log retrieves the log with an id/name
	Log(id string) (Log, error)
}

// Log is an individual event log
type Log interface {
	// Close the log handle
	Close() error
	// Log ID
	Id() string
	// Read will read the next record
	Read() (*Record, error)
	// Go to an offset
	Seek(offset int64) error
	// Write an event to the log
	Write(*Record) error
}

type Record struct {
	Metadata map[string]interface{}
	Data     []byte
}
```

We probably want something similar but new ideas are welcome.

## Service

After the interface has been implemented (which acts as a building block) we want to build 
the service equivalent. The service will effectively provide an RPC interface for event 
streaming based on our go-micro interface.

Ideally we're writing a stream of events to a database where keys are timeseries with 
key-prefix:timestamp (truncated to hour or configurable).

Timeseries ontop of cassandra as example [cassandra/timeseries](https://github.com/HailoOSS/service/tree/master/cassandra/timeseries)


