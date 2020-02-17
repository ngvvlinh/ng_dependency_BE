package sms

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_SendSMS(t *testing.T) {
	tests := []struct {
		name         string
		phoneInput   string
		contentInput string
		want         string
	}{
		{
			name:         "test phone sms",
			phoneInput:   "0364226983",
			contentInput: "Chúc bạn may mắn lần sau",
			want:         "Chúc bạn may mắn lần sau 0364226983",
		},
		{
			name:         "test phone nil",
			phoneInput:   "",
			contentInput: "Chúc bạn may mắn lần sau",
			want:         "",
		},
	}
	cfg := Config{
		Enabled: true,
		Mock:    true,
	}
	cfg.MustLoadEnv()
	client := New(cfg, nil, nil)
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out, err := client.inner.SendSMS(context.Background(), tc.phoneInput, tc.contentInput)
			assert.NoError(t, err, "test fail")
			assert.Equal(t, tc.want, out)
		})
	}
}
