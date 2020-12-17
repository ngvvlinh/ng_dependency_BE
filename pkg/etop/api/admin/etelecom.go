package admin

import (
	"context"

	"o.o/api/etelecom"
	"o.o/api/top/int/admin"
	shoptypes "o.o/api/top/int/shop/types"
	pbcm "o.o/api/top/types/common"
	shopetelecom "o.o/backend/pkg/etop/api/shop/etelecom"
	"o.o/backend/pkg/etop/authorize/session"
)

type EtelecomService struct {
	session.Session

	EtelecomAggr  etelecom.CommandBus
	EtelecomQuery etelecom.QueryBus
}

func (s *EtelecomService) Clone() admin.EtelecomService {
	res := *s
	return &res
}

func (s *EtelecomService) CreateHotline(ctx context.Context, r *admin.CreateHotlineRequest) (*shoptypes.Hotline, error) {
	cmd := &etelecom.CreateHotlineCommand{
		OwnerID:      r.OwnerID,
		Name:         r.Name,
		Hotline:      r.Hotline,
		Network:      r.Network,
		ConnectionID: r.ConnectionID,
		Description:  r.Description,
		IsFreeCharge: r.IsFreeCharge,
	}
	if err := s.EtelecomAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	res := shopetelecom.Convert_etelecom_Hotline_shoptypes_Hotline(cmd.Result, nil)
	return res, nil
}

func (s *EtelecomService) UpdateHotline(ctx context.Context, r *admin.UpdateHotlineRequest) (*pbcm.UpdatedResponse, error) {
	cmd := &etelecom.UpdateHotlineInfoCommand{
		ID:           r.ID,
		IsFreeCharge: r.IsFreeCharge,
		Name:         r.Name,
		Description:  r.Description,
		Status:       r.Status,
	}
	if err := s.EtelecomAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	return &pbcm.UpdatedResponse{Updated: 1}, nil
}
