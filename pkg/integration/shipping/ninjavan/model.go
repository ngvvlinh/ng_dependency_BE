package ninjavan

import cc "o.o/backend/pkg/common/config"

type ClientType byte

const (
	NinjaVanCodeDefault ClientType = 'N'
)

type WebhookConfig struct {
	cc.HTTP `yaml:",inline"`
}

func DefaultWebhookConfig() WebhookConfig {
	return WebhookConfig{
		HTTP: cc.HTTP{Port: 9062},
	}
}
