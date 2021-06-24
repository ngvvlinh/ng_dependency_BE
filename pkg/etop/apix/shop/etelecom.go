package xshop

import (
	"context"

	"o.o/api/etelecom"
	api "o.o/api/top/external/shop"
	externaltypes "o.o/api/top/external/types"
	"o.o/backend/pkg/etop/apix/convertpb"
	"o.o/backend/pkg/etop/authorize/session"
)

type EtelecomService struct {
	session.Session
	EtelecomQuery etelecom.QueryBus
}

func (s *EtelecomService) Clone() api.EtelecomService { res := *s; return &res }

func (s *EtelecomService) GetExtensionInfo(ctx context.Context, r *externaltypes.GetExtensionInfoRequest) (*externaltypes.ExtensionInfo, error) {
	if err := r.Validate(); err != nil {
		return nil, err
	}
	getExtensionQuery := &etelecom.GetExtensionQuery{
		AccountID:       s.SS.Shop().ID,
		ExtensionNumber: r.ExtensionNumber,
	}
	if err := s.EtelecomQuery.Dispatch(ctx, getExtensionQuery); err != nil {
		return nil, err
	}
	extension := getExtensionQuery.Result
	return convertpb.Convert_core_Extension_To_api_ExtensionInfo(extension), nil
}
