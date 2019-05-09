package api

import (
	"context"

	cmP "etop.vn/backend/pb/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/integration/kiotviet"
	kiotvietW "etop.vn/backend/wrapper/integration/kiotviet"
)

func init() {
	bus.AddHandler("kiotviet", VersionInfo)
	bus.AddHandler("kiotviet", SyncProductSource)
}

func VersionInfo(ctx context.Context, q *kiotvietW.VersionInfoEndpoint) error {
	q.Result = &cmP.VersionInfoResponse{
		Service: "etop.Kiotviet",
		Version: "0.1",
	}
	return nil
}

func SyncProductSource(ctx context.Context, q *kiotvietW.SyncProductSourceEndpoint) error {
	cmd := &kiotviet.SyncProductSourceCommand{
		SourceID:      q.Id,
		FromBeginning: q.FromBeginning,
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = &cmP.Empty{}
	return nil
}
