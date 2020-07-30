package telegram

import (
	"context"
	"fmt"

	"o.o/api/meta"
	"o.o/backend/pkg/common/validate"
	"o.o/common/l"
)

var global *Client
var ll = l.New()

type Client struct {
	willFail bool
}

func GetTelegram() *Client {
	if global == nil {
		global = &Client{}
	}
	return global
}

func (c *Client) WillFail() {
	c.willFail = true
}

func (c *Client) SendSMS(ctx context.Context, phone string, content string) (string, error) {
	if phone == "" {
		return "", nil
	}
	tempPhone, isPhone := validate.NormalizePhone(phone)
	if !isPhone {
		return "", nil
	}
	if tempPhone == "" {
		return "", nil
	}
	ll.WithChannel(meta.ChannelSMS).SendMessage(fmt.Sprintf("Send SMS to %v:\n%v", phone, content))
	return content + " " + phone, nil
}
