package sms

import (
	"context"
	"fmt"

	smsing "etop.vn/api/etc/logging/smslog"
	"etop.vn/api/pb/etop/etc/status3"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	cc "etop.vn/backend/pkg/common/config"
	"etop.vn/backend/pkg/common/telebot"
	"etop.vn/backend/pkg/common/validate"
	"etop.vn/backend/pkg/integration/sms/mock"
	"etop.vn/backend/pkg/integration/sms/vietguys"
	"etop.vn/common/l"
)

var ll = l.New()
var smsAggr smsing.CommandBus

func Init(smsCommandBus smsing.CommandBus) {
	smsAggr = smsCommandBus
}

type SendSMSCommand struct {
	Phone   string
	Content string

	Result struct {
		SMSID string
	}
}

type Config struct {
	Enabled  bool            `yaml:"enabled"`
	Mock     bool            `yaml:"mock"`
	Vietguys vietguys.Config `yaml:"vietguys"`
}

func (c *Config) MustLoadEnv(prefix ...string) {
	p := "ET_SMS"
	if len(prefix) > 0 {
		p = prefix[0]
	}
	cc.EnvMap{
		p + "_ENABLED":             &c.Enabled,
		p + "_MOCK":                &c.Mock,
		p + "_VIETGUYS_USERNAME":   &c.Vietguys.Username,
		p + "_VIETGUYS_API_KEY":    &c.Vietguys.APIKey,
		p + "_VIETGUYS_BRAND_NAME": &c.Vietguys.BrandName,
	}.MustLoad()
}

type Client struct {
	inner Driver
	bot   *telebot.Channel
}

func New(cfg Config, bot *telebot.Channel) Client {
	c := Client{
		bot: bot,
	}
	if cfg.Mock {
		c.inner = mock.GetMock()
	} else {
		c.inner = vietguys.New(cfg.Vietguys)
	}
	return c

}

func (c Client) Register(bus bus.Bus) Client {
	bus.AddHandlers(c.SendSMS)
	return c
}

func (c Client) SendSMS(ctx context.Context, cmd *SendSMSCommand) (_err error) {
	phone, _, ok := validate.TrimTest(cmd.Phone)
	if cm.IsDevOrStag() && !ok {
		return cm.Errorf(cm.FailedPrecondition, nil, "Chỉ có thể gửi tin nhắn đến địa chỉ test trên dev!")
	}

	var smsError error
	defer func() {
		createSms := &smsing.CreateSmsLogCommand{
			Content:  cmd.Content,
			Phone:    cmd.Phone,
			Status:   status3.Status_P,
			Provider: "Vietguys",
		}
		if smsError != nil {
			createSms.Status = status3.Status_Z
			createSms.Error = smsError.Error()
		}
		if err := smsAggr.Dispatch(ctx, createSms); err != nil {
			if _err == nil {
				_err = err
			}
		}
	}()

	resp, err := c.inner.SendSMS(ctx, phone, cmd.Content)
	if err != nil {
		smsError = err
		c.bot.SendMessage(fmt.Sprintf("Vietguys: %v", err))
		return cm.Errorf(cm.ExternalServiceError, nil, "Không thể gửi tin nhắn")
	}
	cmd.Result.SMSID = resp
	return nil
}
