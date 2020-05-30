package query

import (
	"context"

	"o.o/api/shopping"
	"o.o/api/shopping/suppliering"
	com "o.o/backend/com/main"
	"o.o/backend/com/shopping/suppliering/sqlstore"
	"o.o/backend/pkg/common/bus"
)

var _ suppliering.QueryService = &SupplierQuery{}

type SupplierQuery struct {
	store sqlstore.SupplierStoreFactory
}

func NewSupplierQuery(db com.MainDB) *SupplierQuery {
	return &SupplierQuery{
		store: sqlstore.NewSupplierStore(db),
	}
}

func SupplierQueryMessageBus(q *SupplierQuery) suppliering.QueryBus {
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
