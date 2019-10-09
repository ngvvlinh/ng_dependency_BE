package query

import (
	"context"

	"etop.vn/api/shopping"
	"etop.vn/api/shopping/vendoring"
	"etop.vn/backend/com/shopping/vendoring/sqlstore"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/cmsql"
)

var _ vendoring.QueryService = &VendorQuery{}

type VendorQuery struct {
	store sqlstore.VendorStoreFactory
}

func NewVendorQuery(db cmsql.Database) *VendorQuery {
	return &VendorQuery{
		store: sqlstore.NewVendorStore(db),
	}
}

func (q *VendorQuery) MessageBus() vendoring.QueryBus {
	b := bus.New()
	return vendoring.NewQueryServiceHandler(q).RegisterHandlers(b)
}

func (q *VendorQuery) GetVendorByID(
	ctx context.Context, args *shopping.IDQueryShopArg,
) (*vendoring.ShopVendor, error) {
	return q.store(ctx).ID(args.ID).OptionalShopID(args.ShopID).GetVendor()
}

func (q *VendorQuery) ListVendors(
	ctx context.Context, args *shopping.ListQueryShopArgs,
) (*vendoring.VendorsResponse, error) {
	query := q.store(ctx).ShopID(args.ShopID).Paging(args.Paging).Filters(args.Filters)
	vendors, err := query.ListVendors()
	if err != nil {
		return nil, err
	}
	count, err := query.Count()
	if err != nil {
		return nil, err
	}
	return &vendoring.VendorsResponse{
		Vendors: vendors,
		Count:   int32(count),
		Paging:  query.GetPaging(),
	}, nil
}

func (q *VendorQuery) ListVendorsByIDs(
	ctx context.Context, args *shopping.IDsQueryShopArgs,
) (*vendoring.VendorsResponse, error) {
	vendors, err := q.store(ctx).ShopID(args.ShopID).IDs(args.IDs...).ListVendors()
	if err != nil {
		return nil, err
	}
	return &vendoring.VendorsResponse{Vendors: vendors}, nil
}
