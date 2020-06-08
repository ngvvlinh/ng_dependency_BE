package sms

import (
	"context"
	"fmt"

	smslog "o.o/api/etc/logging/smslog"
	"o.o/api/top/types/etc/status3"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/whitelabel/wl"
	"o.o/backend/pkg/common/cmenv"
	cc "o.o/backend/pkg/common/config"
	"o.o/backend/pkg/common/validate"
	"o.o/backend/pkg/integration/sms/vietguys"
	"o.o/common/l"
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
	drivers []DriverConfig
	smsLog  smslog.CommandBus
}

type DriverConfig struct {
	Code   string // check for whitelabel, etc.
	Driver Driver
}

func New(cfg Config, drivers []DriverConfig, smsLog smslog.CommandBus) *Client {
	c := Client{
		drivers: drivers,
		smsLog:  smsLog,
	}
	return &c
}

func (c Client) SendSMS(ctx context.Context, cmd *SendSMSCommand) (_err error) {
	client := c.drivers[0].Driver
	key := wl.X(ctx).Key
	for _, d := range c.drivers {
		if key == d.Code {
			client = d.Driver
		}
	}

	phone, _, ok := validate.TrimTest(cmd.Phone)
	if cmenv.IsDevOrStag() && !ok {
		return cm.Errorf(cm.FailedPrecondition, nil, "Chỉ có thể gửi tin nhắn đến địa chỉ test trên dev!")
	}

	resp, err := client.SendSMS(ctx, phone, cmd.Content)
	defer func() {
		createSms := &smslog.CreateSmsLogCommand{
			Content:    cmd.Content,
			Phone:      cmd.Phone,
			Status:     status3.P,
			Provider:   "Vietguys",
			ExternalID: resp,
		}
		if err != nil {
			createSms.Status = status3.Z
			createSms.Error = err.Error()
		}
		if logErr := c.smsLog.Dispatch(ctx, createSms); logErr != nil {
			if _err == nil {
				_err = logErr
			}
		}
	}()

	if err != nil {
		ll.SendMessage(fmt.Sprintf("sms: %v", err))
		return cm.Errorf(cm.ExternalServiceError, nil, "Không thể gửi tin nhắn")
	}
	cmd.Result.SMSID = resp
	return nil
}
