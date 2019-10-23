package sqlstore

import (
	"context"
	"time"

	"etop.vn/api/external/haravan/identity"
	"etop.vn/backend/com/external/haravan/identity/convert"
	cm "etop.vn/backend/pkg/common"

	identitymodel "etop.vn/backend/com/external/haravan/identity/model"
	"etop.vn/backend/pkg/common/cmsql"
)

type XAccountHaravanStoreFactory func(context.Context) *XAccountHaravanStore

func NewXAccountHaravanStore(db *cmsql.Database) XAccountHaravanStoreFactory {
	identitymodel.SQLVerifySchema(db)
	return func(ctx context.Context) *XAccountHaravanStore {
		return &XAccountHaravanStore{
			query: func() cmsql.QueryInterface {
				return cmsql.GetTxOrNewQuery(ctx, *db)
			},
		}
	}
}

type XAccountHaravanStore struct {
	query func() cmsql.QueryInterface
	ft    ExternalAccountHaravanFilters
	preds []interface{}
}

func (s *XAccountHaravanStore) ID(id int64) *XAccountHaravanStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *XAccountHaravanStore) ShopID(id int64) *XAccountHaravanStore {
	s.preds = append(s.preds, s.ft.ByShopID(id))
	return s
}

func (s *XAccountHaravanStore) ExternalShopID(id int) *XAccountHaravanStore {
	s.preds = append(s.preds, s.ft.ByExternalShopID(id))
	return s
}

func (s *XAccountHaravanStore) GetXAccountHaravanDB() (*identitymodel.ExternalAccountHaravan, error) {
	var account identitymodel.ExternalAccountHaravan
	err := s.query().Where(s.preds).ShouldGet(&account)
	return &account, err
}

func (s *XAccountHaravanStore) GetXAccountHaravan() (*identity.ExternalAccountHaravan, error) {
	account, err := s.GetXAccountHaravanDB()
	if err != nil {
		return nil, err
	}
	return convert.XAccountHaravan(account), nil
}

type CreateXAccountHaravanArgs struct {
	ID             int64
	ShopID         int64
	Subdomain      string
	AccessToken    string
	ExternalShopID int
	ExpiresAt      time.Time
}

func (s *XAccountHaravanStore) CreateXAccountHaravan(args *CreateXAccountHaravanArgs) (*identity.ExternalAccountHaravan, error) {
	if args.ID == 0 {
		args.ID = cm.NewID()
	}
	if args.ShopID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing shop_id")
	}
	if args.AccessToken == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing access token")
	}
	if args.Subdomain == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing subdomain")
	}
	account := &identitymodel.ExternalAccountHaravan{
		ID:             args.ID,
		ShopID:         args.ShopID,
		Subdomain:      args.Subdomain,
		AccessToken:    args.AccessToken,
		ExpiresAt:      args.ExpiresAt,
		ExternalShopID: args.ExternalShopID,
	}
	if err := s.query().ShouldInsert(account); err != nil {
		return nil, err
	}

	return s.ID(args.ID).GetXAccountHaravan()
}

type UpdateXAccountHaravanInfoArgs struct {
	ShopID         int64
	Subdomain      string
	AccessToken    string
	ExternalShopID int
	ExpiresAt      time.Time
}

func (s *XAccountHaravanStore) UpdateXAccountHaravan(args *UpdateXAccountHaravanInfoArgs) (*identity.ExternalAccountHaravan, error) {
	if args.ShopID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing Shop ID")
	}
	account := &identitymodel.ExternalAccountHaravan{
		Subdomain:      args.Subdomain,
		AccessToken:    args.AccessToken,
		ExpiresAt:      args.ExpiresAt,
		ExternalShopID: args.ExternalShopID,
	}
	if err := s.query().Where(s.ft.ByShopID(args.ShopID)).ShouldUpdate(account); err != nil {
		return nil, err
	}
	return s.ShopID(args.ShopID).GetXAccountHaravan()
}

type UpdateXShopIDAccountHaravanArgs struct {
	ShopID         int64
	ExternalShopID int
}

func (s *XAccountHaravanStore) UpdateXShopIDAccountHaravan(args *UpdateXShopIDAccountHaravanArgs) (*identity.ExternalAccountHaravan, error) {
	if args.ShopID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing Shop ID")
	}
	if args.ExternalShopID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing External Shop ID Haravan")
	}
	account := &identitymodel.ExternalAccountHaravan{
		ExternalShopID: args.ExternalShopID,
	}
	if err := s.query().Where(s.ft.ByShopID(args.ShopID)).ShouldUpdate(account); err != nil {
		return nil, err
	}
	return s.ShopID(args.ShopID).GetXAccountHaravan()
}

type UpdateXCarrierServiceInfoArgs struct {
	ShopID                            int64
	ExternalCarrierServiceID          int
	ExternalConnectedCarrierServiceAt time.Time
}

func (s *XAccountHaravanStore) UpdateXCarrierServiceInfo(args *UpdateXCarrierServiceInfoArgs) (*identity.ExternalAccountHaravan, error) {
	if args.ShopID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing Shop ID")
	}
	account := &identitymodel.ExternalAccountHaravan{
		ExternalCarrierServiceID:          args.ExternalCarrierServiceID,
		ExternalConnectedCarrierServiceAt: args.ExternalConnectedCarrierServiceAt,
	}
	if err := s.query().Where(s.ft.ByShopID(args.ShopID)).ShouldUpdate(account); err != nil {
		return nil, err
	}
	return s.ShopID(args.ShopID).GetXAccountHaravan()
}

type UpdateDeleteConnectedXCarrierSeriveArgs struct {
	ShopID    int64
	SubDomain string
}

func (s *XAccountHaravanStore) UpdateDeleteConnectedXCarrierService(args *UpdateDeleteConnectedXCarrierSeriveArgs) error {
	if args.ShopID == 0 && args.SubDomain == "" {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing Shop ID or SubDomain")
	}
	s1 := s.query().Table("external_account_haravan")
	if args.ShopID != 0 {
		s1 = s1.Where(s.ft.ByShopID(args.ShopID))
	}
	if args.SubDomain != "" {
		s1 = s1.Where(s.ft.BySubdomain(args.SubDomain))
	}
	if err := s1.ShouldUpdateMap(map[string]interface{}{
		"external_carrier_service_id":           nil,
		"external_connected_carrier_service_at": nil,
	}); err != nil {
		return err
	}
	return nil
}
