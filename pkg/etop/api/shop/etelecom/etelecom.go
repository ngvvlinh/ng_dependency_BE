package etelecom

import (
	"context"

	"o.o/api/etelecom"
	api "o.o/api/top/int/shop"
	pbcm "o.o/api/top/types/common"
	"o.o/backend/pkg/etop/authorize/session"
)

type ExtensionService struct {
	session.Session

	EtelecomAggr  etelecom.CommandBus
	EtelecomQuery etelecom.QueryBus
}

func (s *ExtensionService) Clone() api.ExtensionService {
	res := *s
	return &res
}

func (s *ExtensionService) GetExtensions(ctx context.Context, _ *pbcm.Empty) (*api.GetExtensionResponse, error) {
	query := &etelecom.ListExtensionsQuery{
		UserID: s.SS.Claim().UserID,
	}
	if err := s.EtelecomQuery.Dispatch(ctx, query); err != nil {
		return nil, err
	}
	res := Convert_etelecom_Extensions_shop_Extensions(query.Result)
	return &api.GetExtensionResponse{Extensions: res}, nil
}

func (s *ExtensionService) CreateExtension(ctx context.Context, r *api.CreateExtensionRequest) (*api.Extension, error) {
	cmd := &etelecom.CreateExtensionCommand{
		UserID:    r.UserID,
		AccountID: s.SS.Shop().ID,
	}
	if err := s.EtelecomAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	var res api.Extension
	Convert_etelecom_Extension_shop_Extension(cmd.Result, &res)
	return &res, nil
}
