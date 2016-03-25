package slack

import (
	"errors"
	"sync"

	"github.com/micro/cli"
	"github.com/micro/micro/bot/input"
	"github.com/nlopes/slack"
)

type slackInput struct {
	debug    bool
	token    string
	username string

	sync.Mutex
	running bool
	exit    chan bool

	api *slack.Client
}

func init() {
	input.Inputs["slack"] = new(slackInput)
}

func (p *slackInput) Flags() []cli.Flag {
	return []cli.Flag{
		cli.BoolFlag{
			Name:  "slack_debug",
			Usage: "Slack debug output",
		},
		cli.StringFlag{
			Name:  "slack_token",
			Usage: "Slack token",
		},
	}
}

func (p *slackInput) Init(ctx *cli.Context) error {
	debug := ctx.Bool("slack_debug")
	token := ctx.String("slack_token")

	if len(token) == 0 {
		return errors.New("missing slack token")
	}

	p.debug = debug
	p.token = token

	return nil
}

func (p *slackInput) Connect() (input.Conn, error) {
	p.Lock()
	defer p.Unlock()

	if !p.running {
		return nil, errors.New("not running")
	}

	// test auth
	auth, err := p.api.AuthTest()
	if err != nil {
		return nil, err
	}

	rtm := p.api.NewRTM()
	exit := make(chan bool)

	go rtm.ManageConnection()

	go func() {
		select {
		case <-p.exit:
			close(exit)
		case <-exit:
		}

		rtm.Disconnect()
	}()

	return &slackConn{
		auth: auth,
		rtm:  rtm,
		exit: exit,
	}, nil
}

func (p *slackInput) Start() error {
	if len(p.token) == 0 {
		return errors.New("missing slack token")
	}

	p.Lock()
	defer p.Unlock()

	if p.running {
		return nil
	}

	api := slack.New(p.token)
	api.SetDebug(p.debug)

	// test auth
	_, err := api.AuthTest()
	if err != nil {
		return err
	}

	p.api = api
	p.exit = make(chan bool)
	p.running = true
	return nil
}

func (p *slackInput) Stop() error {
	p.Lock()
	defer p.Unlock()

	if !p.running {
		return nil
	}

	close(p.exit)
	p.running = false
	return nil
}

func (p *slackInput) String() string {
	return "slack"
}
