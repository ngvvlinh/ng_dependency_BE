package sqlstore

import (
	"context"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/cmsql"
	catalogmodel "etop.vn/backend/pkg/services/catalog/model"
	catalogmodelx "etop.vn/backend/pkg/services/catalog/modelx"
)

func init() {
	bus.AddHandler("sql", GetAllProductSources)
	bus.AddHandler("sql", GetProductSource)
}

func GetAllProductSources(ctx context.Context, query *catalogmodelx.GetAllProductSourcesQuery) error {
	s := x.NewQuery()
	if query.External != nil {
		s = s.Exists("external_key", *query.External)
	}
	return s.Find((*catalogmodel.ProductSources)(&query.Result.Sources))
}

func GetProductSource(ctx context.Context, query *catalogmodelx.GetProductSourceQuery) error {
	s, err := getProductSourceProps(&query.GetProductSourceProps, "")
	if err != nil {
		return err
	}

	var productSource catalogmodel.ProductSource
	err = s.ShouldGet(&productSource)
	query.Result = &productSource
	return err
}

func getProductSourceProps(query *catalogmodelx.GetProductSourceProps, prefix string) (cmsql.Query, error) {
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
