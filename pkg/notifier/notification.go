package notifier

import (
	"context"

	"etop.vn/backend/pkg/common/cmsql"
	cc "etop.vn/backend/pkg/common/config"
	"etop.vn/backend/pkg/common/onesignal"
	"etop.vn/backend/pkg/notifier/model"
	"etop.vn/common/bus"
)

func init() {
	bus.AddHandlers("notification",
		CreateNotification,
	)
}

var (
	onesignalClient *onesignal.Client
	x               cmsql.Database
)

func Init(db cmsql.Database, cfg cc.OnesignalConfig) error {
	x = db
	onesignalClient = onesignal.New(cfg.AppID, cfg.ApiKey)
	if err := onesignalClient.Ping(); err != nil {
		return err
	}
	return nil
}

func CreateNotification(ctx context.Context, cmd *model.SendNotificationCommand) error {
	var err error
	req := cmd.Request.ToOnesignalModel()
	cmd.Result, err = onesignalClient.CreateNotification(ctx, req)
	return err
}
