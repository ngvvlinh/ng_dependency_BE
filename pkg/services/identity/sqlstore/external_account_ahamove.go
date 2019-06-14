package sqlstore

import (
	"context"
	"time"

	"etop.vn/api/main/identity"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/services/identity/convert"

	"etop.vn/backend/pkg/common/cmsql"
	identitymodel "etop.vn/backend/pkg/services/identity/model"
)

type XAccountAhamoveStoreFactory func(context.Context) *XAccountAhamoveStore

func NewXAccountAhamoveStore(db cmsql.Database) XAccountAhamoveStoreFactory {
	return func(ctx context.Context) *XAccountAhamoveStore {
		return &XAccountAhamoveStore{
			query: func() cmsql.QueryInterface {
				return cmsql.GetTxOrNewQuery(ctx, db)
			},
		}
	}
}

type XAccountAhamoveStore struct {
	query func() cmsql.QueryInterface
	ft    ExternalAccountAhamoveFilters
	preds []interface{}
}

func (s *XAccountAhamoveStore) ID(id int64) *XAccountAhamoveStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *XAccountAhamoveStore) Phone(phone string) *XAccountAhamoveStore {
	s.preds = append(s.preds, s.ft.ByPhone(phone))
	return s
}

func (s *XAccountAhamoveStore) OwnerID(id int64) *XAccountAhamoveStore {
	s.preds = append(s.preds, s.ft.ByOwnerID(id))
	return s
}

func (s *XAccountAhamoveStore) GetXAccountAhamoveDB() (*identitymodel.ExternalAccountAhamove, error) {
	var account identitymodel.ExternalAccountAhamove
	err := s.query().Where(s.preds).ShouldGet(&account)
	return &account, err
}

func (s *XAccountAhamoveStore) GetXAccountAhamove() (*identity.ExternalAccountAhamove, error) {
	account, err := s.GetXAccountAhamoveDB()
	if err != nil {
		return nil, err
	}
	return convert.XAccountAhamove(account), nil
}

type CreateXAccountAhamoveArgs struct {
	ID      int64
	OwnerID int64
	Phone   string
	Name    string
}

func (s *XAccountAhamoveStore) CreateXAccountAhamove(args *CreateXAccountAhamoveArgs) (*identity.ExternalAccountAhamove, error) {
	if args.ID == 0 {
		args.ID = cm.NewID()
	}
	if args.OwnerID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "CreateXAccountAhamove: Missing owner_id")
	}
	account := &identitymodel.ExternalAccountAhamove{
		ID:      args.ID,
		OwnerID: args.OwnerID,
		Phone:   args.Phone,
		Name:    args.Name,
	}
	if err := s.query().ShouldInsert(account); err != nil {
		return nil, err
	}

	return s.ID(args.ID).GetXAccountAhamove()
}

type UpdateXAccountAhamoveInfoArgs struct {
	ID                int64
	ExternalID        string
	ExternalVerified  bool
	ExternalToken     string
	ExternalCreatedAt time.Time
}

func (s *XAccountAhamoveStore) UpdateXAccountAhamove(args *UpdateXAccountAhamoveInfoArgs) (*identity.ExternalAccountAhamove, error) {
	if args.ID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "UpdateXAccountAhamove: Missing ID")
	}
	account := &identitymodel.ExternalAccountAhamove{
		ExternalID:        args.ExternalID,
		ExternalVerified:  args.ExternalVerified,
		ExternalCreatedAt: args.ExternalCreatedAt,
		ExternalToken:     args.ExternalToken,
	}
	if err := s.query().Where("id = ?", args.ID).ShouldUpdate(account); err != nil {
		return nil, err
	}
	return s.ID(args.ID).GetXAccountAhamove()
}

type UpdateXAccountAhamoveVerifiedInfoArgs struct {
	ID                 int64
	ExternalTickerID   string
	LastSendVerifiedAt time.Time
	ExternalVerified   bool
}

func (s *XAccountAhamoveStore) UpdateXAccountAhamoveVerifiedInfo(args *UpdateXAccountAhamoveVerifiedInfoArgs) (*identity.ExternalAccountAhamove, error) {
	account := &identitymodel.ExternalAccountAhamove{
		ExternalTicketID:   args.ExternalTickerID,
		LastSendVerifiedAt: args.LastSendVerifiedAt,
		ExternalVerified:   args.ExternalVerified,
	}
	if err := s.query().Where("id = ?", args.ID).ShouldUpdate(account); err != nil {
		return nil, err
	}
	return s.ID(args.ID).GetXAccountAhamove()
}

type UpdateXAccountAhamoveVerificationImageArgs struct {
	ID             int64
	IDCardFrontImg string
	IDCardBackImg  string
	PortraitImg    string
}

func (s *XAccountAhamoveStore) UpdateVerificationImages(args *UpdateXAccountAhamoveVerificationImageArgs) (*identity.ExternalAccountAhamove, error) {
	account := &identitymodel.ExternalAccountAhamove{
		IDCardFrontImg: args.IDCardFrontImg,
		IDCardBackImg:  args.IDCardBackImg,
		PortraitImg:    args.PortraitImg,
		UploadedAt:     time.Now(),
	}
	if err := s.query().Where("id = ?", args.ID).ShouldUpdate(account); err != nil {
		return nil, err
	}

	return s.ID(args.ID).GetXAccountAhamove()
}
