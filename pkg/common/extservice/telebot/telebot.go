package telebot

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	tg "github.com/go-telegram-bot-api/telegram-bot-api"

	"o.o/backend/pkg/common/bus"
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

// MaxMessageLength defines maximum length for telegram messages
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
}

func NewChannel(token string, chatID int64) (*Channel, error) {
	bot, err := NewBot(token)
	if err != nil {
		return nil, err
	}
	return &Channel{
		bot:        *bot,
		chatID:     chatID,
		urlSendMsg: bot.baseURL + "sendMessage",
	}, nil
}

func (c *Channel) MaySendMessagef(msg string, args ...interface{}) {
	if c != nil {
		if len(args) > 0 {
			msg = fmt.Sprintf(msg, args...)
		}
		c.SendMessage(msg)
	}
}

func (c *Channel) MaySendMessage(msg string) {
	if c != nil {
		c.SendMessage(msg)
	}
}

func (c *Channel) SendMessage(msg string) {
	err := c.sendMessage("", msg)
	if err != nil {
		ll.Error("Telegram Bot: Unable to send message", l.Error(err))
		msg := "Telegram Bot: Unable to send message: " + err.Error()
		c.sendMessage("", msg)
	}
}

func (c *Channel) SendMarkdownMessage(msg string) {
	err := c.sendMessage("markdown", msg)
	if err != nil {
		ll.Error("Telegram Bot: Unable to send message", l.Error(err))
		msg := "Telegram Bot: Unable to send message: " + err.Error()
		c.sendMessage("", msg)
	}
}

func (c *Channel) sendMessage(mode string, msg string) error {
	if len(msg) >= MaxMessageLength {
		msg = msg[:MaxMessageLength]
	}
	m := ChatMessage{
		ChatID:    c.chatID,
		Text:      msg,
		ParseMode: mode,
	}
	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(m)
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

var (
	chans   = make(map[string]*Channel)
	defChan *Channel
	m       sync.RWMutex
)

func init() {
	bus.AddHandler("telebot", SendMessage)
}

type SendMessageCommand struct {
	Channel string
	Message string
}

func RegisterChannel(name string, ch *Channel) {
	if name == "" {
		defChan = ch
	} else {
		m.Lock()
		chans[name] = ch
		m.Unlock()
	}
}

func SendMessage(ctx context.Context, cmd *SendMessageCommand) error {
	var ch *Channel
	if cmd.Channel == "" {
		ch = defChan
	} else {
		m.RLock()
		channel, ok := chans[cmd.Channel]
		m.RUnlock()
		if !ok {
			panic("No channel available")
		}
		ch = channel
	}
	if ch != nil {
		go ch.SendMessage(cmd.Message)
	}
	return nil
}
