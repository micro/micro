package input

import (
	"github.com/micro/cli"
)

type EventType string

const (
	TextEvent EventType = "text"
)

var (
	Inputs = map[string]Input{}
)

// Event is the unit sent and received
type Event struct {
	Type EventType
	Data []byte
	Meta map[string]interface{}
}

// Input is an interface for sources which
// provide a way to communicate with the bot.
// Slack, HipChat, XMPP, etc.
type Input interface {
	// Provide cli flags
	Flags() []cli.Flag
	// Initialise input using cli context
	Init(*cli.Context) error
	// Connect to the input to
	// sendd and receive events
	Connect() (Conn, error)
	// Start the input
	Start() error
	// Stop the input
	Stop() error
	// name of the input
	String() string
}

type Conn interface {
	Close() error
	Recv(*Event) error
	Send(*Event) error
}
