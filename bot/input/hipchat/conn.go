package hipchat

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/micro/hipchat"
	"github.com/micro/micro/bot/input"
)

type hipchatConn struct {
	exit   chan bool
	client *hipchat.Client
	recv   <-chan *hipchat.Message

	sync.Mutex
	names map[string]*hipchat.User
}

func newConn(c *hipchat.Client) *hipchatConn {
	// setup names
	c.RequestUsers()

	names := make(map[string]*hipchat.User)

	for _, user := range <-c.Users() {
		names[user.Id] = user
	}

	conn := &hipchatConn{
		exit:   make(chan bool),
		client: c,
		names:  names,
		recv:   c.Messages(),
	}

	go conn.run()

	return conn
}

func (c *hipchatConn) run() {
	// name ticker
	t := time.NewTicker(time.Minute)
	defer t.Stop()

	// request users
	users := c.client.Users()
	rooms := c.client.Rooms()

	// join rooms
	c.client.RequestRooms()

	// now we chat
	c.client.Status("chat")
	me := c.getName(c.client.Id).MentionName

	for {
		select {
		case <-c.exit:
			return
		case <-t.C:
			c.client.Ping()
			c.client.RequestRooms()
			c.client.RequestUsers()
		case r := <-rooms:
			for _, room := range r {
				c.client.Join(room.Id, strings.Title(me))
			}
		case u := <-users:
			names := make(map[string]*hipchat.User)
			for _, user := range u {
				names[user.Id] = user
			}
			c.Lock()
			c.names = names
			c.Unlock()
		}
	}
}

func (c *hipchatConn) getName(id string) *hipchat.User {
	c.Lock()
	u := c.names[id]
	c.Unlock()
	return u
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

	me := c.getName(c.client.Id).MentionName

	for {
		select {
		case <-c.exit:
			return errors.New("connection closed")
		case msg := <-c.recv:
			args := strings.Split(msg.Body, " ")
			switch msg.Type {
			// join the room
			case "join_groupchat":
				c.client.Join(msg.From, strings.Title(me))
				continue
			// this is a Room
			case "groupchat":
				// no args
				if len(args) < 2 {
					continue
				}

				// get first arg and check if its us
				name := strings.ToLower(args[0])
				args = args[1:]

				// is it for us?
				switch {
				case strings.HasPrefix(name, "@"+me), strings.HasPrefix(name, me):
				default:
					continue
				}

				// parse from
				parts := strings.Split(msg.From, "/")
				// from channel
				from := parts[0]

				c.Lock()
				names := c.names
				c.Unlock()

				// from user
				var user string
				if len(parts) > 1 {
					for _, u := range names {
						if u.Name == parts[len(parts)-1] {
							user = "@" + u.MentionName
						}
					}
				}

				event.From = from + ":" + user
				event.Data = []byte(strings.Join(args, " "))
				event.To = c.client.Id
			// this is a DM
			case "chat":
				// no args
				if len(args) < 1 {
					continue
				}

				// get first arg and check if its us
				name := strings.ToLower(args[0])

				// parse from
				parts := strings.Split(msg.From, "|")
				// from channel
				from := strings.Split(parts[0], "/")[0]

				// is it for us?
				switch {
				case strings.HasPrefix(name, "@"+me), strings.HasPrefix(name, me):
					args = args[1:]
				}

				event.From = from
				event.To = c.client.Id
				event.Data = []byte(strings.Join(args, " "))
			}

			if event.Meta == nil {
				event.Meta = make(map[string]interface{})
			}

			event.Type = input.TextEvent
			event.Meta["reply"] = msg
			return nil
		}
	}
}

func (c *hipchatConn) Send(event *input.Event) error {
	var channel, to string

	parts := strings.Split(event.To, ":")

	if len(parts) == 2 {
		channel = parts[0]
		to = parts[1]
		message := fmt.Sprintf("%s: %s", to, string(event.Data))
		c.client.Say(channel, to, message)
	} else if len(parts) == 1 {
		channel = parts[0]
		c.client.PM(channel, "", string(event.Data))
	} else {
		return errors.New("could not determine who to send to")
	}

	return nil
}
