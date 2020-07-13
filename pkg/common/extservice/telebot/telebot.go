package telebot

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"o.o/common/l"
)

var (
	ll = l.New()

	client = &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
)

// MaxMessageLength defines maximum length for a telegram message
//
// https://core.telegram.org/method/messages.sendMessage
const MaxMessageLength = 4096

type Bot struct {
	token   string
	baseURL string
}

func NewBot(token string) (*Bot, error) {
	if token == "" {
		panic("Empty token")
	}
	_, err := tg.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	return &Bot{
		token:   token,
		baseURL: "https://api.telegram.org/bot" + token + "/",
	}, nil
}

func (b *Bot) Channel(chatID int64) *Channel {
	if chatID == 0 {
		panic("Empty chat")
	}
	return &Channel{
		bot:    *b,
		chatID: chatID,
	}
}

type Channel struct {
	bot        Bot
	chatID     int64
	urlSendMsg string
	console    l.MockMessenger
	ch         chan string
}

func NewChannel(ctx context.Context, name, token string, chatID int64) (*Channel, error) {
	bot, err := NewBot(token)
	if err != nil {
		return nil, err
	}
	ch := &Channel{
		bot:        *bot,
		chatID:     chatID,
		urlSendMsg: bot.baseURL + "sendMessage",
		console:    l.MockMessenger{Name: name},
		ch:         make(chan string, 256),
	}
	go ch.ListenAndSend(ctx)
	return ch, nil
}

func (c *Channel) ListenAndSend(ctx context.Context) {
	select {
	case <-ctx.Done():
		return

	case msg := <-c.ch:
		c.DoSendMessage(msg)
	}
}

func (c *Channel) SendMessage(msg string) {
	select {
	case c.ch <- msg:
		// ok
	default:
		ll.Warn("dropped message", l.String("msg", msg))
	}
}

func (c *Channel) MaySendMessagef(msg string, args ...interface{}) {
	if c != nil {
		if len(args) > 0 {
			msg = fmt.Sprintf(msg, args...)
		}
		c.DoSendMessage(msg)
	}
}

func (c *Channel) MaySendMessage(msg string) {
	if c != nil {
		c.DoSendMessage(msg)
	}
}

func (c *Channel) DoSendMessage(msg string) {
	defer func() {
		r := recover()
		if r != nil {
			ll.Error("panic", l.Any("recovered", r), l.Stack())
		}
	}()

	err := c.sendMessage("", msg)
	if err != nil {
		ll.Error("Telegram Bot: Unable to send message", l.Error(err))
		errMsg := "Telegram Bot: Unable to send message: " + err.Error()
		_ = c.sendMessage("", errMsg)
	}
}

func (c *Channel) SendMarkdownMessage(msg string) {
	err := c.sendMessage("markdown", msg)
	if err != nil {
		ll.Error("Telegram Bot: Unable to send message", l.Error(err))
		errMsg := "Telegram Bot: Unable to send message: " + err.Error()
		_ = c.sendMessage("", errMsg)
	}
}

func (c *Channel) sendMessage(mode string, msg string) error {
	c.console.SendMessage(msg)
	if len(msg) >= MaxMessageLength {
		msg = msg[:MaxMessageLength]
	}
	m := ChatMessage{
		ChatID:    c.chatID,
		Text:      msg,
		ParseMode: mode,
	}
	var buf bytes.Buffer
	_ = json.NewEncoder(&buf).Encode(m)
	resp, err := client.Post(c.urlSendMsg, "application/json", &buf)
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		errorResp := make(map[string]interface{})
		json.NewDecoder(resp.Body).Decode(&errorResp)
		errorMsg, _ := errorResp["description"].(string)
		if errorMsg == "" {
			errorMsg = "Unknown error from Telegram"
		}
		return errors.New(errorMsg)
	}
	return nil
}

type ChatMessage struct {
	ChatID    int64  `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode,omitempty"`
}
