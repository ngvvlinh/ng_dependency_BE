package sms

import "context"

type Driver interface {
	SendSMS(ctx context.Context, phone string, content string) (string, error)
}
