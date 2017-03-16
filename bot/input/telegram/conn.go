package telegram

import (
	"gopkg.in/telegram-bot-api.v4"
	"github.com/micro/micro/bot/input"
	"errors"
	"strings"
	"sync"
	"github.com/forestgiant/sliceutil"
)

type telegramConn struct {
	api      *tgbotapi.BotAPI
	admins   []string

	recv     <-chan tgbotapi.Update
	exit     chan bool

	syncCond *sync.Cond
	mutex    sync.Mutex
}

func (c *telegramConn) run() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := c.api.GetUpdatesChan(u)
	if err != nil {
		return
	}

	c.recv = updates
	c.syncCond.Signal()

	for {
		select {
		case <-c.exit:
			return
		}
	}
}

func newConn(input *telegramInput) (*telegramConn, error) {
	conn := &telegramConn{
		api: input.api,
		admins: input.admins,
	}

	conn.syncCond = sync.NewCond(&conn.mutex)

	go conn.run()

	return conn, nil
}

func (c *telegramConn) Close() error {
	select {
	case <-c.exit:
		return nil
	default:
		close(c.exit)
	}
	return nil
}

func (c *telegramConn) Recv(event *input.Event) error {
	if event == nil {
		return errors.New("event cannot be nil")
	}

	for {
		if c.recv == nil {
			c.mutex.Lock()
			c.syncCond.Wait()
		}

		select {
		case <-c.exit:
			return errors.New("connection closed")
		case update := <-c.recv:
			if update.Message == nil || !sliceutil.Contains(c.admins, update.Message.From.UserName) {
				continue
			}

			if event.Meta == nil {
				event.Meta = make(map[string]interface{})
			}

			event.Type = input.TextEvent
			event.From = update.Message.From.UserName
			event.Data = []byte(update.Message.Text)
			event.Meta["chatId"] = update.Message.Chat.ID
			event.Meta["chatType"] = update.Message.Chat.Type
			event.Meta["messageId"] = update.Message.MessageID

			return nil
		}
	}
}

func (c *telegramConn) Send(event *input.Event) (err error) {
	messageText := strings.TrimSpace(string(event.Data))

	chatId := event.Meta["chatId"].(int64)
	chatType := ChatType(event.Meta["chatType"].(string))

	msgConfig := tgbotapi.NewMessage(chatId, "<pre>" + messageText + "</pre>")
	msgConfig.ParseMode = "html"

	if sliceutil.Contains([]ChatType{Group, Supergroup}, chatType) {
		msgConfig.ReplyToMessageID = event.Meta["messageId"].(int)
	}

	_, err = c.api.Send(msgConfig)

	return
}
