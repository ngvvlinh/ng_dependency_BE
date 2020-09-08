package notifier

import (
	"context"

	"o.o/backend/com/eventhandler/notifier/model"
	cc "o.o/backend/pkg/common/config"
	"o.o/backend/pkg/common/extservice/onesignal"
)

var NotifyTopics = []string{
	"fulfillment",
	"fb_external_comment",
	"fb_external_message",
	"system",
}

type Notifier struct {
	onesignalClient *onesignal.Client
}

func NewOneSignalNotifier(cfg cc.OnesignalConfig) (*Notifier, error) {
	onesignalClient := onesignal.New(cfg.AppID, cfg.ApiKey)
	if err := onesignalClient.Ping(); err != nil {
		return nil, err
	}
	return &Notifier{
		onesignalClient: onesignalClient,
	}, nil
}

func (n *Notifier) CreateNotification(ctx context.Context, cmd *model.SendNotificationCommand) error {
	var err error
	req := cmd.Request.ToOnesignalModel()
	cmd.Result, err = n.onesignalClient.CreateNotification(ctx, req)
	return err
}
