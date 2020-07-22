package _telebot

import cc "o.o/backend/pkg/common/config"

func DefaultConfig() cc.TelegramBot {
	return cc.TelegramBot{
		Chats: map[string]int64{
			"default": 0,
			"webhook": 0,
			"import":  0,
			"sms":     0,
			"high":    0,
			"deploy":  0,
		},
	}
}
