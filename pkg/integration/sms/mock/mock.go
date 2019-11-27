package mock

import (
	"context"

	cm "etop.vn/backend/pkg/common"

	"etop.vn/backend/pkg/common/validate"
)

var global *Client

type Client struct {
	willFail bool
}

func GetMock() *Client {
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
	n := len(phone) + len(content) + 1
	b := make([]byte, 0, n)
	b = append(b, content...)
	b = append(b, ' ')
	b = append(b, phone...)
	return cm.UnsafeBytesToString(b), nil
}
