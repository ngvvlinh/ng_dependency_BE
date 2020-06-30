package admin

import (
	"context"

	"o.o/api/top/int/admin"
	"o.o/api/top/int/etop"
	identitymodel "o.o/backend/com/main/identity/model"
	identitymodelx "o.o/backend/com/main/identity/modelx"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/etop/api/convertpb"
	"o.o/backend/pkg/etop/authorize/session"
	"o.o/backend/pkg/etop/sqlstore"
)

type AccountService struct {
	session.Session
}

func (s *AccountService) Clone() admin.AccountService {
	res := *s
	return &res
}

func (s *AccountService) CreatePartner(ctx context.Context, q *admin.CreatePartnerRequest) (*etop.Partner, error) {
	cmd := &identitymodelx.CreatePartnerCommand{
		Partner: convertpb.CreatePartnerRequestToModel(q),
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return nil, err
	}
	result := convertpb.PbPartner(cmd.Result.Partner)
	return result, nil
}

func (s *AccountService) GenerateAPIKey(ctx context.Context, q *admin.GenerateAPIKeyRequest) (*admin.GenerateAPIKeyResponse, error) {
	_, err := sqlstore.AccountAuth(ctx).AccountID(q.AccountId).Get()
	if cm.ErrorCode(err) != cm.NotFound {
		return nil, cm.MapError(err).
			Map(cm.OK, cm.AlreadyExists, "account already has an api_key").
			Throw()
	}

	aa := &identitymodel.AccountAuth{
		AccountID:   q.AccountId,
		Status:      1,
		Roles:       nil,
		Permissions: nil,
	}
	err = sqlstore.AccountAuth(ctx).Create(aa)
	result := &admin.GenerateAPIKeyResponse{
		AccountId: q.AccountId,
		ApiKey:    aa.AuthKey,
	}
	return result, err
}
