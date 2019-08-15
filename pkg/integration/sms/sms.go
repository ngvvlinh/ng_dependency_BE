package sms

import (
	"context"
	"fmt"

	cm "etop.vn/backend/pkg/common"
	cc "etop.vn/backend/pkg/common/config"
	"etop.vn/backend/pkg/common/telebot"
	"etop.vn/backend/pkg/common/validate"
	"etop.vn/backend/pkg/integration/sms/vietguys"
	"etop.vn/common/bus"
	"etop.vn/common/l"
)

var ll = l.New()

type SendSMSCommand struct {
	Phone   string
	Content string

	Result struct {
		SMSID string
	}
}

type Config struct {
	Enabled  bool            `yaml:"enabled"`
	Vietguys vietguys.Config `yaml:"vietguys"`
}

func (c *Config) MustLoadEnv(prefix ...string) {
	p := "ET_SMS"
	if len(prefix) > 0 {
		p = prefix[0]
	}
	cc.EnvMap{
		p + "_ENABLED":             &c.Enabled,
		p + "_VIETGUYS_USERNAME":   &c.Vietguys.Username,
		p + "_VIETGUYS_API_KEY":    &c.Vietguys.APIKey,
		p + "_VIETGUYS_BRAND_NAME": &c.Vietguys.BrandName,
	}.MustLoad()
}

type Client struct {
	inner *vietguys.Client
	bot   *telebot.Channel
}

func New(cfg Config, bot *telebot.Channel) Client {
	return Client{
		inner: vietguys.New(cfg.Vietguys),
		bot:   bot,
	}
}

func (c Client) Register(bus bus.Bus) Client {
	bus.AddHandlers(c.SendSMS)
	return c
}

func (c Client) SendSMS(ctx context.Context, cmd *SendSMSCommand) error {
	phone, _, ok := validate.TrimTest(cmd.Phone)
	if cm.IsDevOrStag() && !ok {
		return cm.Errorf(cm.FailedPrecondition, nil, "Chỉ có thể gửi tin nhắn đến địa chỉ test trên dev!")
	}

	resp, err := c.inner.SendSMS(ctx, phone, cmd.Content)
	if err != nil {
		c.bot.SendMessage(fmt.Sprintf("Vietguys: %v", err))
		return cm.Errorf(cm.ExternalServiceError, nil, "Không thể gửi tin nhắn")
	}
	cmd.Result.SMSID = resp
	return nil
}
