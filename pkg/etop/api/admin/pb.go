package admin

import (
	pbadmin "etop.vn/backend/pb/etop/admin"
	"etop.vn/backend/pkg/etop/api/convertpb"
	"etop.vn/backend/pkg/etop/model"
	catalogmodel "etop.vn/backend/pkg/services/catalog/model"
)

func PbCreateCategoryToModel(pb *pbadmin.CreateCategoryRequest) *model.EtopCategory {
	return &model.EtopCategory{
		Name:     pb.Name,
		ParentID: pb.ParentId,
	}
}

func PbVariantWithSupplier(m *catalogmodel.VariantExtended) *pbadmin.VariantWithSupplier {
	return &pbadmin.VariantWithSupplier{
		Variant: convertpb.PbVariant(m),
	}
}

func PbProductWithSupplier(m *catalogmodel.ProductFtVariant) *pbadmin.ProductWithSupplier {
	return &pbadmin.ProductWithSupplier{
		Product: convertpb.PbProduct(m),
	}
}

func VExternalExtendedToVExtended(vxs []*catalogmodel.VariantExternalExtended) []*catalogmodel.VariantExtended {
	if len(vxs) == 0 {
		return nil
	}
	variants := make([]*catalogmodel.VariantExtended, len(vxs))
	for i, v := range vxs {
		variants[i] = &catalogmodel.VariantExtended{
			Variant:         v.Variant,
			VariantExternal: v.VariantExternal,
		}
	}
	return variants
}

func VExtendedToVExternalExtended(vs []*catalogmodel.VariantExtended) []*catalogmodel.VariantExternalExtended {
	if len(vs) == 0 {
		return nil
	}
	variantExternals := make([]*catalogmodel.VariantExternalExtended, len(vs))
	for i, v := range vs {
		variantExternals[i] = &catalogmodel.VariantExternalExtended{
			Variant:         v.Variant,
			VariantExternal: v.VariantExternal,
		}
	}
	return variantExternals
}

func PbUpdateProductToModel(p *pbadmin.UpdateProductRequest) *catalogmodel.Product {
	res := &catalogmodel.Product{
		ID:            p.Id,
		EdName:        p.Name,
		EdShortDesc:   p.ShortDesc,
		EdDescription: p.Description,
		EdDescHTML:    p.DescHtml,
	}
	return res
}

func PbUpdateVariantToModel(p *pbadmin.UpdateVariantRequest) *catalogmodel.Variant {
	res := &catalogmodel.Variant{
		ID:            p.Id,
		EdName:        p.Name,
		EdShortDesc:   p.ShortDesc,
		EdDescription: p.Description,
		EdDescHTML:    p.DescHtml,
	}
	return res
}
