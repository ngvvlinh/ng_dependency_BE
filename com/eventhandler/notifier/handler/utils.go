package handler

import (
	"context"
	"encoding/json"

	notifiermodel "o.o/backend/com/eventhandler/notifier/model"
	fbpage "o.o/backend/com/fabo/main/fbpage/model"
	cm "o.o/backend/pkg/common"
	"o.o/capi/dot"
	"o.o/common/l"
)

func getPage(fbExternalID string) (*fbpage.FbExternalPage, error) {
	var page fbpage.FbExternalPage
	err := x.Table("fb_external_page").
		Where("external_id = ?", fbExternalID).ShouldGet(&page)
	if err != nil {
		return nil, err
	}
	return &page, nil
}

type buildNotifyCmdArgs struct {
	UserIDs    []dot.ID
	ShopID     dot.ID
	Title      string
	Message    string
	SendNotify bool
	Entity     notifiermodel.NotiEntity
	EntityID   dot.ID
	Meta       interface{}
	TopicType  string
}

func buildNotifyCmds(
	args *buildNotifyCmdArgs,
) []*notifiermodel.CreateNotificationArgs {
	var res []*notifiermodel.CreateNotificationArgs
	var _meta []byte
	if args.Meta != nil {
		_meta, _ = json.Marshal(args.Meta)
	}

	for _, userID := range args.UserIDs {
		cmd := &notifiermodel.CreateNotificationArgs{
			AccountID:        args.ShopID,
			UserID:           userID,
			Title:            args.Title,
			Message:          args.Message,
			EntityID:         args.EntityID,
			Entity:           args.Entity,
			SendNotification: args.SendNotify,
			MetaData:         _meta,
			TopicType:        args.TopicType,
		}
		res = append(res, cmd)
	}
	return res
}

func createNotifications(_ context.Context, cmds []*notifiermodel.CreateNotificationArgs) error {
	if len(cmds) == 0 {
		return nil
	}

	chErr := make(chan error, len(cmds))
	for _, cmd := range cmds {
		go func(_cmd *notifiermodel.CreateNotificationArgs) (_err error) {
			defer func() {
				chErr <- _err
			}()
			defer cm.RecoverAndLog()
			_, _err = notifyStore.CreateNotification(_cmd)
			if _err != nil {
				ll.Debug("err", l.Error(_err))
			}
			return
		}(cmd)
	}
	var created, errors int
	for i, n := 0, len(cmds); i < n; i++ {
		err := <-chErr
		if err == nil {
			created++
		} else {
			errors++
		}
	}
	ll.S.Infof("Create notifications: success %v/%v, errors %v/%v", created, len(cmds), errors, len(cmds))
	return nil
}
