package sqlstore

import (
	"context"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/etop/model"
)

func init() {
	bus.AddHandler("sql", GetAllProductSources)
	bus.AddHandler("sql", GetProductSource)
	bus.AddHandler("sql", GetProductSourceExtended)
	bus.AddHandler("sql", GetOrderSource)
	bus.AddHandler("sql", GetOrderSourceExtended)
}

func GetAllProductSources(ctx context.Context, query *model.GetAllProductSourcesQuery) error {
	s := x.NewQuery()
	if query.External != nil {
		s = s.Exists("external_key", *query.External)
	}
	if query.Supplier != nil {
		s = s.Exists("supplier_id", *query.Supplier)
	}
	return s.Find((*model.ProductSources)(&query.Result.Sources))
}

func GetProductSource(ctx context.Context, query *model.GetProductSourceQuery) error {
	s, err := getProductSourceProps(&query.GetProductSourceProps, "")
	if err != nil {
		return err
	}

	var productSource model.ProductSource
	err = s.ShouldGet(&productSource)
	query.Result = &productSource
	return err
}

func GetProductSourceExtended(ctx context.Context, query *model.GetProductSourceExtendedQuery) error {
	s, err := getProductSourceProps(&query.GetProductSourceProps, "ps.")
	if err != nil {
		return err
	}

	var productSource model.ProductSourceExtended
	err = s.ShouldGet(&productSource)
	query.Result = &productSource
	return err
}

func getProductSourceProps(query *model.GetProductSourceProps, prefix string) (cmsql.Query, error) {
	switch {
	case query.ID != 0:
	case query.Type != "" && query.ExternalKey != "":
	default:
		return cmsql.Query{}, cm.Error(cm.InvalidArgument, "Missing required params", nil)
	}

	s := x.NewQuery()
	if query.ID != 0 {
		s = s.Where(prefix+"id = ?", query.ID)
	}
	if query.Type != "" {
		s = s.Where(prefix+"type = ?", query.Type)
	}
	if query.ExternalKey != "" {
		s = s.Where(prefix+"external_key = ?", query.ExternalKey)
	}
	return s, nil
}

func GetOrderSource(ctx context.Context, query *model.GetOrderSourceQuery) error {
	s, err := getOrderSourceProps(&query.GetOrderSourceProps, "")
	if err != nil {
		return err
	}

	var orderSource model.OrderSource
	err = s.ShouldGet(&orderSource)
	query.Result = &orderSource
	return err
}

func GetOrderSourceExtended(ctx context.Context, query *model.GetOrderSourceExtendedQuery) error {
	s, err := getOrderSourceProps(&query.GetOrderSourceProps, "os.")
	if err != nil {
		return err
	}

	var orderSource model.OrderSourceExtended
	err = s.ShouldGet(&orderSource)
	query.Result = &orderSource
	return err
}

func getOrderSourceProps(query *model.GetOrderSourceProps, prefix string) (cmsql.Query, error) {
	switch {
	case query.ID != 0:
	case query.ShopID != 0:
	case query.Type != "" && query.ExternalKey != "":
	default:
		return cmsql.Query{}, cm.Error(cm.InvalidArgument, "Missing required params", nil)
	}

	s := x.NewQuery()
	if query.ID != 0 {
		s = s.Where(prefix+"id = ?", query.ID)
	}
	if query.Type != "" {
		s = s.Where(prefix+"type = ?", query.Type)
	}
	if query.ExternalKey != "" {
		s = s.Where(prefix+"external_key = ?", query.ExternalKey)
	}
	if query.ShopID != 0 {
		s = s.Where(prefix+"shop_id = ?", query.ShopID)
	}
	return s, nil
}
