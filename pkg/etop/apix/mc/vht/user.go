package vht

import (
	"context"

	"o.o/api/main/identity"
	"o.o/api/top/external/mc/vht"
	pbcm "o.o/api/top/types/common"
	cm "o.o/backend/pkg/common"
	wldriver "o.o/backend/pkg/common/apifw/whitelabel/drivers"
	"o.o/backend/pkg/etop/authorize/session"
)

var _ vht.UserService = &VHTUserService{}

type VHTUserService struct {
	session.Session

	IdentityAggr  identity.CommandBus
	IdentityQuery identity.QueryBus
}

func (s *VHTUserService) Clone() vht.UserService {
	res := *s
	return &res
}

func (s *VHTUserService) RegisterUser(ctx context.Context, r *vht.VHTRegisterUser) (*pbcm.Empty, error) {
	partner := s.SS.Partner()
	if partner.ID != wldriver.VNPostID {
		return nil, cm.Errorf(cm.PermissionDenied, nil, "")
	}

	if r.Phone == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing phone")
	}
	if r.FullName == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing fullname")
	}

	query := &identity.GetUserByPhoneQuery{
		Phone: r.Phone,
	}
	err := s.IdentityQuery.Dispatch(ctx, query)
	if err == nil {
		// user đã tồn tại
		return &pbcm.Empty{}, nil
	}

	if err != nil && cm.ErrorCode(err) != cm.NotFound {
		return nil, err
	}

	// Tạo user mới
	cmd := &identity.RegisterSimplifyCommand{
		Phone:    r.Phone,
		FullName: r.FullName,
	}
	if err := s.IdentityAggr.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	return &pbcm.Empty{}, nil
}
