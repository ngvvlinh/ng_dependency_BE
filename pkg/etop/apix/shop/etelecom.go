package xshop

import (
	"context"

	"o.o/api/etelecom"
	"o.o/api/main/identity"
	api "o.o/api/top/external/shop"
	externaltypes "o.o/api/top/external/types"
	"o.o/backend/pkg/etop/apix/convertpb"
	"o.o/backend/pkg/etop/authorize/session"
	"o.o/capi/dot"
)

type EtelecomService struct {
	session.Session
	EtelecomQuery etelecom.QueryBus
	IdentityQuery identity.QueryBus
}

func (s *EtelecomService) Clone() api.EtelecomService { res := *s; return &res }

func (s *EtelecomService) GetExtensionInfo(ctx context.Context, r *externaltypes.GetExtensionInfoRequest) (*externaltypes.ExtensionInfo, error) {
	if err := r.Validate(); err != nil {
		return nil, err
	}
	var userID dot.ID
	if r.Phone != "" || r.Email != "" {
		getUserQuery := &identity.GetUserByPhoneOrEmailQuery{
			Phone: r.Phone,
			Email: r.Email,
		}

		if err := s.IdentityQuery.Dispatch(ctx, getUserQuery); err != nil {
			return nil, err
		}
		userID = getUserQuery.Result.ID
	}

	getExtensionQuery := &etelecom.GetExtensionQuery{
		AccountID:       s.SS.Shop().ID,
		ExtensionNumber: r.ExtensionNumber,
		UserID:          userID,
	}
	if err := s.EtelecomQuery.Dispatch(ctx, getExtensionQuery); err != nil {
		return nil, err
	}
	extension := getExtensionQuery.Result
	return convertpb.Convert_core_Extension_To_api_ExtensionInfo(extension), nil
}
