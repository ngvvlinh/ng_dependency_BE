package wl

import (
	"context"

	"etop.vn/backend/pkg/common/apifw/whitelabel"
	"etop.vn/backend/pkg/common/apifw/whitelabel/drivers"
	"etop.vn/backend/pkg/common/cmenv"
	"etop.vn/capi/dot"
)

var whiteLabel *whitelabel.WhiteLabel

func Init(env cmenv.EnvType) *whitelabel.WhiteLabel {
	if whiteLabel != nil {
		panic("already init")
	}
	whiteLabel = whitelabel.New(drivers.Drivers(env))
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

func GetWLPartnerID(ctx context.Context) dot.ID {
	wlPartner := X(ctx)
	if wlPartner.IsWhiteLabel() {
		return wlPartner.ID
	}
	return 0
}
