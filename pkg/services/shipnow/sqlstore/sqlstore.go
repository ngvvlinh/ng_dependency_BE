package sqlstore

import (
	"context"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/services/shipnow/model"
)

type ShipnowStore struct {
	ctx context.Context
	db  cmsql.Database
}

func NewShipnowStore(db cmsql.Database) *ShipnowStore {
	return &ShipnowStore{db: db, ctx: context.Background()}
}

type GetByIDArgs struct {
	ID int64

	ShopID    int64
	PartnerID int64
}

func (s *ShipnowStore) WithContext(ctx context.Context) *ShipnowStore {
	return &ShipnowStore{
		ctx: ctx,
		db:  s.db,
	}
}

func (s *ShipnowStore) GetByID(args GetByIDArgs) (result *model.ShipnowFulfillment, err error) {
	if args.ID == 0 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "missing id")
	}

	q := s.db.WithContext(s.ctx).Where("id = ?")
	if args.ShopID != 0 {
		q = q.Where("shop_id = ?", args.ShopID)
	}
	if args.PartnerID != 0 {
		q = q.Where("partner_id = ?", args.PartnerID)
	}

	result = &model.ShipnowFulfillment{}
	err = q.ShouldGet(result)
	return result, err
}
