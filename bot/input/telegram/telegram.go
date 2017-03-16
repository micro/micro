package telegram

import (
	"github.com/micro/cli"
	"github.com/micro/micro/bot/input"
	"errors"
	"sync"
	"gopkg.in/telegram-bot-api.v4"
)

type telegramInput struct {
	sync.Mutex

	debug  bool
	token  string
	admins []string

	api    *tgbotapi.BotAPI
}

type ChatType string

const (
	Private ChatType = "private"
	Group ChatType = "group"
	Supergroup ChatType = "supergroup"
)

func init() {
	input.Inputs["telegram"] = &telegramInput{}
}

func (t *telegramInput) Flags() []cli.Flag {
	return []cli.Flag{
		cli.BoolFlag{
			Name:  "telegram_debug",
			Usage: "Telegram debug output",
		},
		cli.StringFlag{
			Name:  "telegram_token",
			Usage: "Telegram token",
		},
		cli.StringSliceFlag{
			Name:  "telegram_admins",
			Usage: "Telegram bot's administrators",
		},
	}
}

func (t *telegramInput) Init(ctx *cli.Context) error {
	t.debug = ctx.Bool("telegram_debug")
	t.token = ctx.String("telegram_token")
	t.admins = ctx.StringSlice("telegram_admins")

	if len(t.token) == 0 {
		return errors.New("missing telegram token")
	}

	return nil
}

func (t *telegramInput) Stream() (input.Conn, error) {
	t.Lock()
	defer t.Unlock()
	return newConn(t)
}

func (t *telegramInput) Start() error {
	t.Lock()
	defer t.Unlock()

	api, err := tgbotapi.NewBotAPI(t.token)
	if err != nil {
		return err
	}

	t.api = api

	api.Debug = t.debug

	return nil
}

func (t *telegramInput) Stop() error {
	return nil
}

func (p *telegramInput) String() string {
	return "telegram"
}
