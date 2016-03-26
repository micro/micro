package hipchat

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/micro/cli"
	"github.com/micro/micro/bot/command"
	"github.com/micro/micro/bot/input"

	"github.com/micro/hipchat"
)

type hipchatInput struct {
	username string
	password string
	server   string
	debug    bool

	sync.Mutex
	running bool
	exit    chan bool

	client *hipchat.Client

	ctx  sync.RWMutex
	cmds map[string]command.Command
}

func init() {
	input.Inputs["hipchat"] = NewInput()
}

func (h *hipchatInput) exec(args []string) (string, error) {
	h.ctx.RLock()
	defer h.ctx.RUnlock()

	for _, cmd := range h.cmds {
		if args[0] != cmd.Name() {
			continue
		}

		rsp, err := cmd.Exec(args...)
		if err != nil {
			return "", err
		}

		return string(rsp), nil
	}

	return "", nil
}

func (h *hipchatInput) run() {
	t := time.NewTicker(time.Minute)
	defer t.Stop()

	h.client.RequestUsers()

	names := make(map[string]*hipchat.User)
	messages := h.client.Messages()
	users := h.client.Users()

	fn := func(users []*hipchat.User) map[string]*hipchat.User {
		names := make(map[string]*hipchat.User)
		for _, user := range users {
			names[user.Id] = user
		}
		return names
	}

	// get users
	names = fn(<-users)

	// now we chat
	h.client.Status("chat")
	me := names[h.client.Id].MentionName

	// join rooms
	h.client.RequestRooms()
	for _, room := range <-h.client.Rooms() {
		h.client.Join(room.Id, strings.Title(me))
	}

	for {
		select {
		case <-h.exit:
			return
		case <-t.C:
			h.client.RequestUsers()
		case u := <-users:
			names = fn(u)
		case msg := <-messages:
			args := strings.Split(msg.Body, " ")

			switch msg.Type {
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

				// from user
				var user string
				if len(parts) > 1 {
					for _, u := range names {
						if u.Name == parts[len(parts)-1] {
							user = "@" + u.MentionName
						}
					}
				}

				// execute command
				rsp, err := h.exec(args)
				if err != nil {
					text := fmt.Sprintf("error executing command: %v", err)
					if len(user) > 0 {
						text = fmt.Sprintf("%s: %s", user, text)
					}
					h.client.Say(from, user, text)
					continue
				}

				// don't respond if no response
				if len(rsp) == 0 {
					continue
				}

				// format response
				text := fmt.Sprintf("%s: %s", user, rsp)
				if len(user) == 0 {
					text = rsp
				}

				// send response
				h.client.Say(from, user, text)
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

				// execute command
				rsp, err := h.exec(args)
				if err != nil {
					text := fmt.Sprintf("error executing command: %v", err)
					h.client.Say(from, "", text)
					continue
				}

				// if there's no output, don't respond
				if len(rsp) == 0 {
					continue
				}

				// send private message
				h.client.PM(from, "", rsp)
			// this is an Invite
			case "join_groupchat":
				h.client.Join(msg.From, strings.Title(me))
				// TODO: save rooms we're in
			// we just don't know
			default:
				if h.debug {
					fmt.Printf("[bot][hipchat] unknown message received %+v\n", msg)
				}
			}
		}
	}
}

func (h *hipchatInput) Flags() []cli.Flag {
	return []cli.Flag{
		cli.BoolFlag{
			Name:  "hipchat_debug",
			Usage: "Hipchat debug output",
		},
		cli.StringFlag{
			Name:  "hipchat_username",
			Usage: "Hipchat XMPP username",
		},
		cli.StringFlag{
			Name:  "hipchat_password",
			Usage: "Hipchat XMPP password",
		},
		cli.StringFlag{
			Name:  "hipchat_server",
			Usage: "Hipchat XMPP server",
			Value: "chat.hipchat.com",
		},
	}
}

func (h *hipchatInput) Init(ctx *cli.Context) error {
	username := ctx.String("hipchat_username")
	password := ctx.String("hipchat_password")
	server := ctx.String("hipchat_server")
	debug := ctx.Bool("hipchat_debug")

	if len(username) == 0 {
		return errors.New("require username")
	}

	if len(password) == 0 {
		return errors.New("require password")
	}

	if len(server) == 0 {
		return errors.New("require server")
	}

	h.username = username
	h.password = password
	h.server = server
	h.debug = debug

	return nil
}

func (h *hipchatInput) Connect() (input.Conn, error) {
	h.Lock()
	defer h.Unlock()

	if !h.running {
		return nil, errors.New("not running")
	}

	// TODO: return conn
	// TODO: use server url
	c, err := hipchat.NewClient(h.username, h.password, "bot")
	if err != nil {
		return nil, err
	}

	exit := make(chan bool)

	go func() {
		select {
		case <-h.exit:
			select {
			case <-exit:
				return
			default:
				close(exit)
			}
		case <-exit:
			return
		}

		c.Close()
	}()

	return &hipchatConn{
		exit:   exit,
		client: c,
	}, nil
}

func (h *hipchatInput) Process(cmd command.Command) error {
	h.ctx.Lock()
	defer h.ctx.Unlock()

	if _, ok := h.cmds[cmd.Name()]; ok {
		return errors.New("Command with name " + cmd.Name() + " already exists")
	}

	h.cmds[cmd.Name()] = cmd
	return nil
}

func (h *hipchatInput) Start() error {
	if len(h.username) == 0 || len(h.password) == 0 || len(h.server) == 0 {
		return errors.New("missing hipchat configuration")
	}

	h.Lock()
	defer h.Unlock()

	if h.running {
		return nil
	}

	// TODO: use server url
	c, err := hipchat.NewClient(h.username, h.password, "bot")
	if err != nil {
		return err
	}

	h.client = c
	h.exit = make(chan bool)
	h.running = true
	go h.run()

	return nil
}

func (h *hipchatInput) Stop() error {
	h.Lock()
	defer h.Unlock()

	if !h.running {
		return nil
	}

	close(h.exit)
	h.running = false
	return nil
}

func (h *hipchatInput) String() string {
	return "hipchat"
}

func NewInput() input.Input {
	return &hipchatInput{
		cmds: make(map[string]command.Command),
	}
}
