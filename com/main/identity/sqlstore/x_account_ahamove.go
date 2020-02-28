package sqlstore

import (
	"context"
	"encoding/json"
	"time"

	"etop.vn/api/main/identity"
	"etop.vn/backend/com/main/identity/convert"
	"etop.vn/backend/com/main/identity/model"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/sql/cmsql"
	"etop.vn/capi/dot"
)

type XAccountAhamoveStoreFactory func(context.Context) *XAccountAhamoveStore

func NewXAccountAhamoveStore(db *cmsql.Database) XAccountAhamoveStoreFactory {
	model.SQLVerifySchema(db)
	return func(ctx context.Context) *XAccountAhamoveStore {
		return &XAccountAhamoveStore{
			query: cmsql.NewQueryFactory(ctx, db),
		}
	}
}

type XAccountAhamoveStore struct {
	query cmsql.QueryFactory
	ft    ExternalAccountAhamoveFilters
	preds []interface{}
}

func (s *XAccountAhamoveStore) ID(id dot.ID) *XAccountAhamoveStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *XAccountAhamoveStore) Phone(phone string) *XAccountAhamoveStore {
	s.preds = append(s.preds, s.ft.ByPhone(phone))
	return s
}

func (s *XAccountAhamoveStore) ExternalID(externalID string) *XAccountAhamoveStore {
	s.preds = append(s.preds, s.ft.ByExternalID(externalID))
	return s
}

func (s *XAccountAhamoveStore) OwnerID(id dot.ID) *XAccountAhamoveStore {
	s.preds = append(s.preds, s.ft.ByOwnerID(id))
	return s
}

func (s *XAccountAhamoveStore) GetXAccountAhamoveDB() (*model.ExternalAccountAhamove, error) {
	var account model.ExternalAccountAhamove
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
	ID      dot.ID
	OwnerID dot.ID
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
	account := &model.ExternalAccountAhamove{
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
	ID                dot.ID
	ExternalID        string
	ExternalVerified  bool
	ExternalToken     string
	ExternalCreatedAt time.Time
}

func (s *XAccountAhamoveStore) UpdateXAccountAhamove(args *UpdateXAccountAhamoveInfoArgs) (*identity.ExternalAccountAhamove, error) {
	if args.ID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "UpdateXAccountAhamove: Missing ID")
	}
	account := &model.ExternalAccountAhamove{
		ExternalID:        args.ExternalID,
		ExternalVerified:  args.ExternalVerified,
		ExternalCreatedAt: args.ExternalCreatedAt,
		ExternalToken:     args.ExternalToken,
	}
	if err := s.query().Where(s.ft.ByID(args.ID)).ShouldUpdate(account); err != nil {
		return nil, err
	}
	return s.ID(args.ID).GetXAccountAhamove()
}

type UpdateXAccountAhamoveVerifiedInfoArgs struct {
	ID                   dot.ID
	ExternalTickerID     string
	LastSendVerifiedAt   time.Time
	ExternalVerified     bool
	ExternalDataVerified json.RawMessage
}

func (s *XAccountAhamoveStore) UpdateXAccountAhamoveVerifiedInfo(args *UpdateXAccountAhamoveVerifiedInfoArgs) (*identity.ExternalAccountAhamove, error) {
	account := &model.ExternalAccountAhamove{
		ExternalTicketID:     args.ExternalTickerID,
		LastSendVerifiedAt:   args.LastSendVerifiedAt,
		ExternalVerified:     args.ExternalVerified,
		ExternalDataVerified: args.ExternalDataVerified,
	}
	if err := s.query().Where(s.ft.ByID(args.ID)).ShouldUpdate(account); err != nil {
		return nil, err
	}
	return s.ID(args.ID).GetXAccountAhamove()
}

type UpdateXAccountAhamoveVerificationImageArgs struct {
	ID                  dot.ID
	IDCardFrontImg      string
	IDCardBackImg       string
	PortraitImg         string
	WebsiteURL          string
	FanpageURL          string
	CompanyImgs         []string
	BusinessLicenseImgs []string
}

func (s *XAccountAhamoveStore) UpdateVerificationImages(args *UpdateXAccountAhamoveVerificationImageArgs) (*identity.ExternalAccountAhamove, error) {
	account := &model.ExternalAccountAhamove{
		IDCardFrontImg:      args.IDCardFrontImg,
		IDCardBackImg:       args.IDCardBackImg,
		PortraitImg:         args.PortraitImg,
		WebsiteURL:          args.WebsiteURL,
		FanpageURL:          args.FanpageURL,
		CompanyImgs:         args.CompanyImgs,
		BusinessLicenseImgs: args.BusinessLicenseImgs,
		UploadedAt:          time.Now(),
	}
	if err := s.query().Where(s.ft.ByID(args.ID)).ShouldUpdate(account); err != nil {
		return nil, err
	}

	return s.ID(args.ID).GetXAccountAhamove()
}
