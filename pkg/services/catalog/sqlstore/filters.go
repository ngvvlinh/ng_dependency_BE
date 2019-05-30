package sqlstore

import sq "etop.vn/backend/pkg/common/sq"

func (ft ProductFilters) NotDeleted() sq.WriterTo {
	return ft.Filter("$.deleted_at IS NULL")
}

func (ft VariantFilters) NotDeleted() sq.WriterTo {
	return ft.Filter("$.deleted_at IS NULL")
}

func (ft ShopProductFilters) NotDeleted() sq.WriterTo {
	// shop_product does not use deleted_at
	return nil
}

func (ft ShopVariantFilters) NotDeleted() sq.WriterTo {
	// shop_variant does not use deleted_at
	return nil
}
