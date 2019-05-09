package admin

import (
	pbadmin "etop.vn/backend/pb/etop/admin"
	pbsupplier "etop.vn/backend/pb/etop/supplier"
	"etop.vn/backend/pkg/etop/model"
)

func PbCreateCategoryToModel(pb *pbadmin.CreateCategoryRequest) *model.EtopCategory {
	return &model.EtopCategory{
		Name:     pb.Name,
		ParentID: pb.ParentId,
	}
}

func PbVariantsWithSupplier(items []*model.VariantExtended) []*pbadmin.VariantWithSupplier {
	res := make([]*pbadmin.VariantWithSupplier, len(items))
	for i, item := range items {
		res[i] = PbVariantWithSupplier(item)
	}
	return res
}

func PbVariantWithSupplier(m *model.VariantExtended) *pbadmin.VariantWithSupplier {
	return &pbadmin.VariantWithSupplier{
		Variant: pbsupplier.PbVariant(m),
	}
}

func PbProductWithSupplier(m *model.ProductFtVariant) *pbadmin.ProductWithSupplier {
	return &pbadmin.ProductWithSupplier{
		Product: pbsupplier.PbProduct(m),
	}
}

func VExternalExtendedToVExtended(vxs []*model.VariantExternalExtended) []*model.VariantExtended {
	if len(vxs) == 0 {
		return nil
	}
	variants := make([]*model.VariantExtended, len(vxs))
	for i, v := range vxs {
		variants[i] = &model.VariantExtended{
			Variant:         v.Variant,
			VariantExternal: v.VariantExternal,
		}
	}
	return variants
}

func VExtendedToVExternalExtended(vs []*model.VariantExtended) []*model.VariantExternalExtended {
	if len(vs) == 0 {
		return nil
	}
	variantExternals := make([]*model.VariantExternalExtended, len(vs))
	for i, v := range vs {
		variantExternals[i] = &model.VariantExternalExtended{
			Variant:         v.Variant,
			VariantExternal: v.VariantExternal,
		}
	}
	return variantExternals
}

func PbUpdateProductToModel(p *pbadmin.UpdateProductRequest) *model.Product {
	res := &model.Product{
		ID:            p.Id,
		EdName:        p.Name,
		EdShortDesc:   p.ShortDesc,
		EdDescription: p.Description,
		EdDescHTML:    p.DescHtml,
	}
	return res
}

func PbUpdateVariantToModel(p *pbadmin.UpdateVariantRequest) *model.Variant {
	res := &model.Variant{
		ID:            p.Id,
		EdName:        p.Name,
		EdShortDesc:   p.ShortDesc,
		EdDescription: p.Description,
		EdDescHTML:    p.DescHtml,
	}
	return res
}
