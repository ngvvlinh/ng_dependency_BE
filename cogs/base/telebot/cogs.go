package _telebot

import (
	"o.o/api/meta"
	cc "o.o/backend/pkg/common/config"
)

func DefaultConfig() cc.TelegramBot {
	return cc.TelegramBot{
		Chats: map[string]int64{
			meta.ChannelDefault:         0,
			meta.ChannelWebhook:         0,
			meta.ChannelImport:          0,
			meta.ChannelSMS:             0,
			meta.ChannelHigh:            0,
			meta.ChannelDeploy:          0,
			meta.ChannelShipmentCarrier: 0,
			meta.ChannelTelecomProvider: 0,
		},
	}
}
