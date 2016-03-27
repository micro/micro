package hipchat

import (
	"errors"
	"sync"

	"github.com/micro/cli"
	"github.com/micro/hipchat"
	"github.com/micro/micro/bot/input"
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
}

func init() {
	input.Inputs["hipchat"] = NewInput()
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

func (h *hipchatInput) Stream() (input.Conn, error) {
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

	return newConn(c), nil
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
	return &hipchatInput{}
}
