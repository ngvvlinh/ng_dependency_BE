package notifier

import (
	"context"

	"o.o/backend/com/handler/notifier/model"
	"o.o/backend/pkg/common/bus"
	cc "o.o/backend/pkg/common/config"
	"o.o/backend/pkg/common/extservice/onesignal"
	"o.o/backend/pkg/common/sql/cmsql"
)

func init() {
	bus.AddHandlers("notification",
		CreateNotification,
	)
}

var (
	onesignalClient *onesignal.Client
	x               *cmsql.Database
)

func Init(db *cmsql.Database, cfg cc.OnesignalConfig) error {
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
