package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"path"

	identitycore "o.o/api/main/identity"
	"o.o/api/top/types/etc/account_type"
	"o.o/backend/com/main/identity"
	identitymodel "o.o/backend/com/main/identity/model"
	identitymodelx "o.o/backend/com/main/identity/modelx"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/whitelabel/wl"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/cmenv"
	cc "o.o/backend/pkg/common/config"
	"o.o/backend/pkg/common/projectpath"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/etop/sqlstore"
	"o.o/backend/scripts/migration/partner/list"
	"o.o/common/l"
)

var (
	ll               = l.New()
	db               *cmsql.Database
	identityAggr     *identity.Aggregate
	identityQuery    *identity.QueryService
	partnerStore     sqlstore.PartnerStoreFactory
	accountAuthStore sqlstore.AccountAuthStoreFactory
)

func main() {
	cfg := cc.DefaultPostgres()
	cfg.Host = "localhost"
	cfg.Port = 5432
	cfg.Username = "postgres"
	cfg.Password = "postgres"
	cfg.Database = "etopv1.12"

	var err error
	db, err = cmsql.Connect(cfg)
	if err != nil {
		ll.Fatal("Error while connecting database", l.Error(err))
	}
	var ctx = bus.Ctx()
	_ = wl.Init(cmenv.EnvDev, wl.EtopServer)
	evenBus := bus.New()
	identityAggr = identity.NewAggregate(db, evenBus)
	identityQuery = identity.NewQueryService(db)
	partnerStore = sqlstore.NewPartnerStore(db)
	accountAuthStore = sqlstore.NewAccountAuthStore(db)

	err = db.InTransaction(ctx, func(tx cmsql.QueryInterface) error {
		for _, partner := range list.PartnerList {
			if err := createPartner(ctx, tx, partner); err != nil {
				ll.Error("", l.Object("partner", partner))
				return err
			}

			// create api key
			accountAuth := &identitymodel.AccountAuth{
				AccountID: partner.AccountID,
				AuthKey:   partner.AuthToken,
				Status:    1,
			}
			err2 := accountAuthStore(ctx).Create(accountAuth)
			if err2 != nil {
				return err2
			}

			partner.AuthToken = accountAuth.AuthKey
		}
		return nil
	})
	if err != nil {
		ll.Fatal("Create partner failed!", l.Error(err))
	}

	data, _ := json.MarshalIndent(list.PartnerList, "", " ")

	dir := path.Join(projectpath.GetPath(), "scripts", "migration", "partner")
	fileName := path.Join(dir, "partners.json")
	_ = ioutil.WriteFile(fileName, data, 0644)
	ll.Info("export file success")
}

func createPartner(ctx context.Context, tx cmsql.QueryInterface, partner *list.PartnerInfo) error {
	if err := partner.Validate(); err != nil {
		return err
	}
	if err := createUser(ctx, tx, partner); err != nil {
		return err
	}

	switch partner.AccountType {
	case account_type.Partner:
		cmd := &identitymodelx.CreatePartnerCommand{
			Partner: &identitymodel.Partner{
				ID:         partner.AccountID,
				OwnerID:    partner.OwnerID,
				Status:     1,
				Name:       partner.Fullname,
				PublicName: partner.Fullname,
				Phone:      partner.Phone,
				Email:      partner.Email,
			},
		}
		return partnerStore(ctx).CreatePartner(ctx, cmd)
	case account_type.Shop:
		cmd := &identitycore.CreateShopArgs{
			ID:      partner.AccountID,
			Name:    partner.Fullname,
			OwnerID: partner.OwnerID,
			Phone:   partner.Phone,
			Email:   partner.Email,
		}
		_, err := identityAggr.CreateShop(ctx, cmd)
		if err != nil {
			return err
		}
		return nil
	default:
		return cm.Errorf(cm.InvalidArgument, nil, "Account type does not valid")
	}
}

func createUser(ctx context.Context, tx cmsql.QueryInterface, partner *list.PartnerInfo) error {
	if partner.OwnerID != 0 {
		query := &identitycore.GetUserByIDQueryArgs{
			UserID: partner.OwnerID,
		}
		_, err := identityQuery.GetUserByID(ctx, query)
		if err == nil {
			// user existed
			// skip create
			return nil
		}
		if err != nil && cm.ErrorCode(err) != cm.NotFound {
			return err
		}
	}

	password := partner.Password
	if password == "" {
		password = "123456789"
	}
	args := &identitycore.CreateUserArgs{
		UserID:   partner.OwnerID,
		FullName: partner.Fullname,
		Email:    partner.Email,
		Phone:    partner.Phone,
		Password: password,
		Status:   1,
	}
	_, err := identityAggr.CreateUser(ctx, args)
	if err != nil {
		return err
	}
	return nil
}
