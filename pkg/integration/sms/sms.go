package sms

import (
	"context"

	"etop.vn/backend/pkg/common/l"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	cc "etop.vn/backend/pkg/common/config"
	"etop.vn/backend/pkg/common/validate"
	"etop.vn/backend/pkg/integration/sms/esms"
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
	Enabled bool        `yaml:"enabled"`
	ESMS    esms.Config `yaml:"esms"`
}

func (c *Config) MustLoadEnv(prefix ...string) {
	p := "ET_SMS"
	if len(prefix) > 0 {
		p = prefix[0]
	}
	cc.EnvMap{
		p + "_ENABLED":         &c.Enabled,
		p + "_ESMS_BASE_URL":   &c.ESMS.BaseURL,
		p + "_ESMS_API_KEY":    &c.ESMS.APIKey,
		p + "_ESMS_SECRET_KEY": &c.ESMS.SecretKey,
		p + "_ESMS_BRAND_NAME": &c.ESMS.BrandName,
	}.MustLoad()
}

type Client struct {
	esms *esms.ESMS
}

func New(cfg Config) Client {
	return Client{esms.New(cfg.ESMS)}
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

	resp, err := c.esms.SendSMS(ctx, esms.SMSTypeBrandCustomer, phone, cmd.Content)
	if err != nil {
		ll.Error("can not send sms", l.Error(err))
		return cm.Errorf(cm.ExternalServiceError, nil, "Không thể gửi tin nhắn")
	}
	cmd.Result.SMSID = resp.SMSID
	return nil
}
