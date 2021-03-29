package etelecomuser

import (
	"context"

	"o.o/api/etelecom/usersetting"
	etelecomapi "o.o/api/top/int/etelecom"
	etelecomtypes "o.o/api/top/int/etelecom/types"
	pbcm "o.o/api/top/types/common"
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
	query := &usersetting.GetUserSettingQuery{
		UserID: s.SS.Shop().OwnerID,
	}
	if err := s.UserSettingQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	res := convertpb.Convert_usersetting_UserSetting_api_UserSetting(query.Result)
	return res, nil
}
