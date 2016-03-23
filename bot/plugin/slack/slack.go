package slack

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/micro/cli"
	"github.com/micro/micro/bot/plugin"
	"github.com/nlopes/slack"
)

type slackPlugin struct {
	debug    bool
	token    string
	username string

	sync.Mutex
	running bool
	exit    chan bool
}

func init() {
	plugin.Plugins["slack"] = new(slackPlugin)
}

func (p *slackPlugin) process(auth *slack.AuthTestResponse, rtm *slack.RTM, e slack.RTMEvent) error {
	switch ev := e.Data.(type) {
	case *slack.MessageEvent:
		switch {
		case strings.HasPrefix(ev.Channel, "D"):
		case strings.HasPrefix(ev.Text, auth.User):
		case strings.HasPrefix(ev.Text, fmt.Sprintf("<@%s>", auth.UserID)):
		default:
			return nil
		}
		fmt.Printf("Received message to me %+v", ev)
	case *slack.InvalidAuthEvent:
		return errors.New("invalid credentials")
	}

	return nil
}

func (p *slackPlugin) run(auth *slack.AuthTestResponse, rtm *slack.RTM, exit chan bool) {
	for {
		select {
		case <-exit:
			return
		case e := <-rtm.IncomingEvents:
			if err := p.process(auth, rtm, e); err != nil {
				return
			}
		}
	}
}

func (p *slackPlugin) Flags() []cli.Flag {
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

func (p *slackPlugin) Init(ctx *cli.Context) error {
	debug := ctx.Bool("slack_debug")
	token := ctx.String("slack_token")

	if len(token) == 0 {
		return errors.New("missing slack token")
	}

	p.debug = debug
	p.token = token

	return nil
}

func (p *slackPlugin) Start() error {
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
	auth, err := api.AuthTest()
	if err != nil {
		return err
	}

	rtm := api.NewRTM()
	exit := make(chan bool)

	go rtm.ManageConnection()
	go p.run(auth, rtm, exit)

	p.exit = exit
	p.running = true
	return nil
}

func (p *slackPlugin) Stop() error {
	p.Lock()
	defer p.Unlock()

	if !p.running {
		return nil
	}

	close(p.exit)
	p.running = false
	return nil
}

func (p *slackPlugin) String() string {
	return "slack"
}
