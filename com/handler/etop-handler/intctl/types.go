package intctl

import (
	"bytes"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/capi/dot"
)

func Topic(prefix string) string {
	if prefix == "" {
		ll.Fatal("no prefix")
	}
	return prefix + "_intctl"
}

const ChannelReloadWebhook = "webhook/reload"

func NewKey(channel string) string {
	return channel + ":" + cm.IDToDec(cm.NewID())
}

func ParseKey(key []byte) string {
	parts := bytes.Split(key, []byte{':'})
	return string(parts[0])
}

type ReloadWebhook struct {
	AccountID dot.ID `json:"account_id"`
}
