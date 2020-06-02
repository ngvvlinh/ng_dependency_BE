package notifier

import (
	"context"

	"o.o/backend/com/eventhandler/notifier/model"
	"o.o/backend/pkg/common/bus"
	cc "o.o/backend/pkg/common/config"
	"o.o/backend/pkg/common/extservice/onesignal"
)

func init() {
	bus.AddHandlers("notification",
		CreateNotification,
	)
}

var onesignalClient *onesignal.Client

func Init(cfg cc.OnesignalConfig) error {
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
