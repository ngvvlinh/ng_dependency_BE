package wl

import (
	"context"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/apifw/whitelabel"
	"etop.vn/backend/pkg/common/apifw/whitelabel/drivers"
	"etop.vn/capi/dot"
)

var whiteLabel *whitelabel.WhiteLabel

func Init(env cm.EnvType) *whitelabel.WhiteLabel {
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
