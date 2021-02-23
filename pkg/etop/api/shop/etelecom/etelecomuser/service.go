package etelecomuser

import (
	"context"

	"o.o/api/etelecom/usersetting"
	etelecomapi "o.o/api/top/int/etelecom"
	pbcm "o.o/api/top/types/common"
	"o.o/backend/pkg/etop/api/shop/etelecom"
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

func (s *EtelecomUserService) UpdateUserSetting(ctx context.Context, r *etelecomapi.UpdateUserSettingRequest) (*pbcm.UpdatedResponse, error) {
	cmd := &usersetting.UpdateUserSettingCommand{
		UserID:              s.SS.Shop().OwnerID,
		ExtensionChargeType: r.ExtensionChargeType,
	}
	if err := s.UserSettingAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	return &pbcm.UpdatedResponse{Updated: 1}, nil
}

func (s *EtelecomUserService) GetUserSetting(ctx context.Context, empty *pbcm.Empty) (*etelecomapi.EtelecomUserSetting, error) {
	query := &usersetting.GetUserSettingQuery{
		UserID: s.SS.Shop().OwnerID,
	}
	if err := s.UserSettingQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	res := etelecom.Convert_usersetting_UserSetting_api_UserSetting(query.Result)
	return res, nil
}
