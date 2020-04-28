package convertpb

import (
	"o.o/api/main/catalog"
	exttypes "o.o/api/top/external/types"
	catalogmodel "o.o/backend/com/main/catalog/model"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/etop/api/convertpb"
	"o.o/capi/dot"
	"o.o/capi/util"
)

func PbShopProduct(arg *catalog.ShopProduct) *exttypes.ShopProduct {
	if arg == nil {
		return nil
	}
	if arg.Deleted {
		return &exttypes.ShopProduct{
			Id:      arg.ProductID,
			Deleted: true,
		}
	}
	return &exttypes.ShopProduct{
		ExternalId:   dot.String(arg.ExternalID),
		ExternalCode: dot.String(arg.ExternalCode),
		Id:           arg.ProductID,
		Name:         dot.String(arg.Name),
		Description:  dot.String(arg.Description),
		ShortDesc:    dot.String(arg.ShortDesc),
		ImageUrls:    arg.ImageURLs,
		CategoryId:   arg.CategoryID.Wrap(),
		Note:         dot.String(arg.Note),
		Status:       arg.Status.Wrap(),
		ListPrice:    dot.Int(arg.ListPrice),
		RetailPrice:  dot.Int(arg.RetailPrice),
		CreatedAt:    cmapi.PbTime(arg.CreatedAt),
		UpdatedAt:    cmapi.PbTime(arg.UpdatedAt),
		BrandId:      arg.BrandID.Wrap(),
	}
}

func ToNullStrings(items []string) (res []dot.NullString) {
	for _, item := range items {
		res = append(res, dot.String(item))
	}
	return
}

func PbShopProducts(args []*catalog.ShopProduct) []*exttypes.ShopProduct {
	outs := make([]*exttypes.ShopProduct, len(args))
	for i, arg := range args {
		outs[i] = PbShopProduct(arg)
	}
	return outs
}

func ConvertProductWithVariantsToPbProduct(arg *catalog.ShopProductWithVariants) *exttypes.ShopProduct {
	if arg == nil {
		return nil
	}
	return &exttypes.ShopProduct{
		ExternalId:   dot.String(arg.ExternalID),
		ExternalCode: dot.String(arg.ExternalCode),
		Id:           arg.ProductID,
		Name:         dot.String(arg.Name),
		Description:  dot.String(arg.Description),
		ShortDesc:    dot.String(arg.ShortDesc),
		ImageUrls:    arg.ImageURLs,
		CategoryId:   arg.CategoryID.Wrap(),
		Note:         dot.String(arg.Note),
		Status:       arg.Status.Wrap(),
		ListPrice:    dot.Int(arg.ListPrice),
		RetailPrice:  dot.Int(arg.RetailPrice),
		CreatedAt:    cmapi.PbTime(arg.CreatedAt),
		UpdatedAt:    cmapi.PbTime(arg.UpdatedAt),
		BrandId:      arg.BrandID.Wrap(),
	}
}

func PbShopVariant(arg *catalog.ShopVariant) *exttypes.ShopVariant {
	if arg == nil {
		return nil
	}
	if arg.Deleted {
		return &exttypes.ShopVariant{
			Id:      arg.VariantID,
			Deleted: true,
		}
	}
	return &exttypes.ShopVariant{
		ExternalId:   dot.String(arg.ExternalID),
		ExternalCode: dot.String(arg.ExternalCode),
		Id:           arg.VariantID,
		Code:         dot.String(arg.Code),
		Name:         dot.String(arg.Name),
		Description:  dot.String(arg.Description),
		ShortDesc:    dot.String(arg.ShortDesc),
		ImageUrls:    arg.ImageURLs,
		ListPrice:    dot.Int(arg.ListPrice),
		RetailPrice:  dot.Int(util.CoalesceInt(arg.RetailPrice, arg.ListPrice)),
		Note:         dot.String(arg.Note),
		Status:       arg.Status.Wrap(),
		CostPrice:    dot.Int(arg.CostPrice),
		Attributes:   arg.Attributes,
	}
}

func PbShopVariants(args []*catalog.ShopVariant) []*exttypes.ShopVariant {
	outs := make([]*exttypes.ShopVariant, len(args))
	for i, arg := range args {
		outs[i] = PbShopVariant(arg)
	}
	return outs
}

func PbShopVariantHistory(m catalogmodel.ShopVariantHistory) *exttypes.ShopVariant {
	return &exttypes.ShopVariant{
		ExternalId:   m.ExternalID().String(),
		ExternalCode: m.ExternalCode().String(),
		Id:           m.VariantID().ID().Apply(0),
		Code:         m.Code().String(),
		Name:         m.Name().String(),
		Description:  m.Description().String(),
		ShortDesc:    m.ShortDesc().String(),
		ImageUrls:    nil, // TOOO: fill it
		ListPrice:    m.ListPrice().Int(),
		RetailPrice:  m.RetailPrice().Int(),
		Note:         m.Note().String(),
		Status:       convertpb.Pb3Ptr(m.Status().Int()),
		CostPrice:    m.CostPrice().Int(),
		Attributes:   nil, // TODO: fill it`
	}
}

func PbShopProductHistory(m catalogmodel.ShopProductHistory) *exttypes.ShopProduct {
	return &exttypes.ShopProduct{
		ExternalId:   m.ExternalID().String(),
		ExternalCode: m.ExternalCode().String(),
		Id:           m.ProductID().ID().Apply(0),
		Name:         m.Name().String(),
		Description:  m.Description().String(),
		ShortDesc:    m.ShortDesc().String(),
		ImageUrls:    nil, // TODO: fill it
		CategoryId:   m.CategoryID().ID(),
		Note:         m.Note().String(),
		Status:       convertpb.Pb3Ptr(m.Status().Int()),
		ListPrice:    m.ListPrice().Int(),
		RetailPrice:  m.RetailPrice().Int(),
		Variants:     nil,
		CreatedAt:    cmapi.PbTime(m.CreatedAt().Time()),
		UpdatedAt:    cmapi.PbTime(m.UpdatedAt().Time()),
		BrandId:      m.BrandID().ID(),
	}
}
