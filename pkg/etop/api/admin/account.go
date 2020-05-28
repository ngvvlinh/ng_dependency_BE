package admin

import (
	"context"

	"o.o/api/top/int/admin"
	identitymodel "o.o/backend/com/main/identity/model"
	identitymodelx "o.o/backend/com/main/identity/modelx"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/etop/api/convertpb"
	"o.o/backend/pkg/etop/sqlstore"
)

type AccountService struct{}

func (s *AccountService) Clone() *AccountService {
	res := *s
	return &res
}

func (s *AccountService) CreatePartner(ctx context.Context, q *CreatePartnerEndpoint) error {
	cmd := &identitymodelx.CreatePartnerCommand{
		Partner: convertpb.CreatePartnerRequestToModel(q.CreatePartnerRequest),
	}
	if err := bus.Dispatch(ctx, cmd); err != nil {
		return err
	}
	q.Result = convertpb.PbPartner(cmd.Result.Partner)
	return nil
}

func (s *AccountService) GenerateAPIKey(ctx context.Context, q *GenerateAPIKeyEndpoint) error {
	_, err := sqlstore.AccountAuth(ctx).AccountID(q.AccountId).Get()
	if cm.ErrorCode(err) != cm.NotFound {
		return cm.MapError(err).
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
	q.Result = &admin.GenerateAPIKeyResponse{
		AccountId: q.AccountId,
		ApiKey:    aa.AuthKey,
	}
	return err
}
