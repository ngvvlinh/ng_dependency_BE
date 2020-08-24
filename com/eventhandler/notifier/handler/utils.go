package handler

import (
	"context"
	"encoding/json"

	notifiermodel "o.o/backend/com/eventhandler/notifier/model"
	fbpage "o.o/backend/com/fabo/main/fbpage/model"
	"o.o/capi/dot"
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

func getUserIDsWithShopID(ctx context.Context, shopID dot.ID) ([]dot.ID, error) {
	accUsers, err := accountUserStore(ctx).ByAccountID(shopID).ListAccountUserDBs()
	if err != nil {
		return nil, err
	}

	var userIDs []dot.ID
	for _, accUser := range accUsers {
		userIDs = append(userIDs, accUser.UserID)
	}
	return userIDs, nil
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
}

func buildNotifyCmds(
	args *buildNotifyCmdArgs,
) []*notifiermodel.CreateNotificationArgs {
	var res []*notifiermodel.CreateNotificationArgs
	for _, userID := range args.UserIDs {
		_meta, _ := json.Marshal(args.Meta)
		cmd := &notifiermodel.CreateNotificationArgs{
			AccountID:        args.ShopID,
			UserID:           userID,
			Title:            args.Title,
			Message:          args.Message,
			EntityID:         args.EntityID,
			Entity:           args.Entity,
			SendNotification: args.SendNotify,
			MetaData:         _meta,
		}
		res = append(res, cmd)
	}
	return res
}
