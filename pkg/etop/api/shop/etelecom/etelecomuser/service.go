package etelecomuser

import (
	"context"

	"o.o/api/etelecom/usersetting"
	etelecomapi "o.o/api/top/int/etelecom"
	etelecomtypes "o.o/api/top/int/etelecom/types"
	pbcm "o.o/api/top/types/common"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/etop/api/convertpb"
	"o.o/backend/pkg/etop/authorize/session"
)

type EtelecomUserService struct {
	session.Session

	UserSettingAggr  usersetting.CommandBus
	UserSettingQuery usersetting.QueryBus
}

func (s *EtelecomUserService) Clone() etelecomapi.EtelecomUserService {
	res := *s
	return &res
}

func (s *EtelecomUserService) GetUserSetting(ctx context.Context, empty *pbcm.Empty) (*etelecomtypes.EtelecomUserSetting, error) {
	userID := s.SS.Shop().OwnerID
	query := &usersetting.GetUserSettingQuery{
		UserID: userID,
	}
	err := s.UserSettingQuery.Dispatch(ctx, query)
	userSetting := &usersetting.UserSetting{}
	switch cm.ErrorCode(err) {
	case cm.NoError:
		userSetting = query.Result
	case cm.NotFound:
		userSetting = &usersetting.UserSetting{ID: userID}
	default:
		return nil, err
	}
	res := convertpb.Convert_usersetting_UserSetting_api_UserSetting(userSetting)
	return res, nil
}
