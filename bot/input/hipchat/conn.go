package hipchat

import (
	"errors"

	"github.com/micro/hipchat"
	"github.com/micro/micro/bot/input"
)

type hipchatConn struct {
	exit   chan bool
	client *hipchat.Client
}

func (c *hipchatConn) Close() error {
	select {
	case <-c.exit:
		return nil
	default:
		close(c.exit)
	}
	return nil
}

func (c *hipchatConn) Recv(event *input.Event) error {
	if event == nil {
		return errors.New("event cannot be nil")
	}

	messages := c.client.Messages()

	for {
		select {
		case <-c.exit:
			return errors.New("connection closed")
		case msg := <-messages:
			if event.Meta == nil {
				event.Meta = make(map[string]interface{})
			}
			event.Type = input.TextEvent
			event.Data = []byte(msg.Body)
			event.Meta["reply"] = msg
			return nil
		}
	}
}

func (c *hipchatConn) Send(event *input.Event) error {
	ev, ok := event.Meta["reply"]
	if !ok {
		return errors.New("can't correlate")
	}

	msg := ev.(*hipchat.Message)
	c.client.Say(msg.From, "", string(event.Data))
	return nil
}
