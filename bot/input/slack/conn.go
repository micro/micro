package slack

import (
	"errors"
	"fmt"
	"strings"

	"github.com/micro/micro/bot/input"
	"github.com/nlopes/slack"
)

// Satisfies the input.Conn interface
type slackConn struct {
	auth *slack.AuthTestResponse
	rtm  *slack.RTM
	exit chan bool
}

func (s *slackConn) Close() error {
	select {
	case <-s.exit:
		return nil
	default:
		close(s.exit)
	}
	return nil
}

func (s *slackConn) Recv(event *input.Event) error {
	if event == nil {
		return errors.New("event cannot be nil")
	}

	for {
		select {
		case <-s.exit:
			return errors.New("connection closed")
		case e := <-s.rtm.IncomingEvents:
			switch ev := e.Data.(type) {
			case *slack.MessageEvent:
				// only accept type message?
				if ev.Type != "message" {
					continue
				}

				// only accept DMs or messages to me?
				switch {
				case strings.HasPrefix(ev.Channel, "D"):
				case strings.HasPrefix(ev.Text, s.auth.User):
				case strings.HasPrefix(ev.Text, fmt.Sprintf("<@%s>", s.auth.UserID)):
				default:
					continue
				}

				if event.Meta == nil {
					event.Meta = make(map[string]interface{})
				}

				// fill in the blanks
				event.Type = input.TextEvent
				event.Data = []byte(ev.Text)
				event.Meta["reply"] = ev
				return nil
			case *slack.InvalidAuthEvent:
				return errors.New("invalid credentials")
			}
		}
	}
}

func (s *slackConn) Send(event *input.Event) error {
	ev, ok := event.Meta["reply"]
	if !ok {
		return errors.New("can't correlate")
	}
	msg := s.rtm.NewOutgoingMessage(string(event.Data), ev.(*slack.MessageEvent).Channel)
	s.rtm.SendMessage(msg)
	return nil
}
