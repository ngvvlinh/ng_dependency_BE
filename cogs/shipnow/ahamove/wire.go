package _ahamove

import (
	"github.com/google/wire"

	"o.o/backend/pkg/integration/shipnow/ahamove"
	"o.o/backend/pkg/integration/shipnow/ahamove/client"
	"o.o/backend/pkg/integration/shipnow/ahamove/server"
	"o.o/backend/pkg/integration/shipnow/ahamove/webhook"
)

var WireSet = wire.NewSet(
	webhook.New,
	client.New,
	ahamove.NewCarrierAccount,
	ahamove.New,
	server.NewAhamoveWebhookServer,
	server.NewAhamoveVerificationFileServer,
)
