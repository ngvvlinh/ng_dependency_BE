package wl

import (
	"context"

	"o.o/backend/pkg/common/apifw/whitelabel"
	etopdrivers "o.o/backend/pkg/common/apifw/whitelabel/drivers"
	fabodrivers "o.o/backend/pkg/common/apifw/whitelabel/fabo"
	"o.o/backend/pkg/common/cmenv"
	"o.o/capi/dot"
)

type ServerName string

const (
	EtopServer ServerName = "etop"
	FaboServer ServerName = "fabo"
)

var whiteLabel *whitelabel.WhiteLabel

func Init(env cmenv.EnvType, serverName ServerName) *whitelabel.WhiteLabel {
	if whiteLabel != nil {
		panic("already init")
	}
	switch serverName {
	case EtopServer:
		whiteLabel = whitelabel.New(etopdrivers.Drivers(env))
	case FaboServer:
		whiteLabel = whitelabel.New(fabodrivers.Drivers(env))
	default:
		// TODO define later
	}

	return whiteLabel
}

func X(ctx context.Context) *whitelabel.WL {
	if whiteLabel == nil {
		panic("whitelabel has not been initialized")
	}
	return whiteLabel.ByContext(ctx)
}

func WrapContext(ctx context.Context, partnerID dot.ID) context.Context {
	return whiteLabel.WrapContext(ctx, partnerID)
}

func WrapContextByPartnerID(ctx context.Context, partnerID dot.ID) context.Context {
	return whiteLabel.WrapContextByPartnerID(ctx, partnerID)
}

func GetWLPartnerID(ctx context.Context) dot.ID {
	wlPartner := X(ctx)
	if wlPartner.IsWhiteLabel() {
		return wlPartner.ID
	}
	return 0
}

func GetWLPartnerByKey(key string) *whitelabel.WL {
	if whiteLabel == nil {
		panic("whitelabel has not been initialized")
	}
	return whiteLabel.ByPartnerKey(key)
}
