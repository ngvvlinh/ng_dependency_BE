package query

import (
	"context"

	"etop.vn/api/shopping"
	"etop.vn/api/shopping/suppliering"
	"etop.vn/backend/com/shopping/suppliering/sqlstore"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/sql/cmsql"
)

var _ suppliering.QueryService = &SupplierQuery{}

type SupplierQuery struct {
	store sqlstore.SupplierStoreFactory
}

func NewSupplierQuery(db *cmsql.Database) *SupplierQuery {
	return &SupplierQuery{
		store: sqlstore.NewSupplierStore(db),
	}
}

func (q *SupplierQuery) MessageBus() suppliering.QueryBus {
	b := bus.New()
	return suppliering.NewQueryServiceHandler(q).RegisterHandlers(b)
}

func (q *SupplierQuery) GetSupplierByID(
	ctx context.Context, args *shopping.IDQueryShopArg,
) (*suppliering.ShopSupplier, error) {
	return q.store(ctx).ID(args.ID).OptionalShopID(args.ShopID).GetSupplier()
}

func (q *SupplierQuery) ListSuppliers(
	ctx context.Context, args *shopping.ListQueryShopArgs,
) (*suppliering.SuppliersResponse, error) {
	query := q.store(ctx).ShopID(args.ShopID).WithPaging(args.Paging).Filters(args.Filters)
	suppliers, err := query.ListSuppliers()
	if err != nil {
		return nil, err
	}
	return &suppliering.SuppliersResponse{
		Suppliers: suppliers,
		Paging:    query.GetPaging(),
	}, nil
}

func (q *SupplierQuery) ListSuppliersByIDs(
	ctx context.Context, args *shopping.IDsQueryShopArgs,
) (*suppliering.SuppliersResponse, error) {
	suppliers, err := q.store(ctx).ShopID(args.ShopID).IDs(args.IDs...).ListSuppliers()
	if err != nil {
		return nil, err
	}
	return &suppliering.SuppliersResponse{Suppliers: suppliers}, nil
}
