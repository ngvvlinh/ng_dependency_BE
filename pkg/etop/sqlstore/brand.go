package sqlstore

import (
	"context"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/etop/model"
)

func init() {
	bus.AddHandler("sql", GetProductBrand)
	bus.AddHandler("sql", GetProductBrands)
	bus.AddHandler("sql", CreateProductBrand)
	bus.AddHandler("sql", UpdateProductBrand)
	bus.AddHandler("sql", DeleteProductBrand)
}

func GetProductBrand(ctx context.Context, query *model.GetProductBrandQuery) error {
	if query.ID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing ID", nil)
	}
	p := new(model.ProductBrandExtended)
	s := x.
		Table("product_brand").
		Where("pb.id = ?", query.ID)
	if query.SupplierID != 0 {
		s = s.Where("pb.supplier_id = ?", query.SupplierID)
	}
	has, err := s.Get(p)
	if err != nil {
		return err
	}
	if !has {
		return cm.Error(cm.NotFound, "", nil)
	}

	query.Result = p
	return nil
}

func GetProductBrands(ctx context.Context, query *model.GetProductBrandsQuery) error {
	s := x.Table("product_brand")
	if query.SupplierID != 0 {
		s = s.Where("pb.supplier_id = ?", query.SupplierID)
	}
	if query.Ids != nil {
		s = s.In("pb.id", query.Ids)
	}
	err := s.Find((*model.ProductBrandExtendeds)(&query.Result.Brands))
	return err
}

func CreateProductBrand(ctx context.Context, cmd *model.CreateProductBrandCommand) error {
	brand := cmd.Brand
	if brand.SupplierID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing SupplierID", nil)
	}

	brand.ID = cm.NewID()
	if _, err := x.Table("product_brand").Insert(brand); err != nil {
		return err
	}
	cmd.Result = brand
	return nil
}

func UpdateProductBrand(ctx context.Context, cmd *model.UpdateProductBrandCommand) error {
	brand := cmd.Brand
	if brand.ID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing BrandID", nil)
	}
	if brand.SupplierID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing SupplierID", nil)
	}

	if err := x.Table("product_brand").
		Where("id = ?", brand.ID).
		Where("supplier_id = ?", brand.SupplierID).
		ShouldUpdate(brand); err != nil {
		return err
	}
	cmd.Result = brand
	return nil
}

func DeleteProductBrand(ctx context.Context, cmd *model.DeleteProductBrandCommand) error {
	return inTransaction(func(s Qx) error {
		if cmd.SupplierID == 0 {
			return cm.Error(cm.InvalidArgument, "Missing SupplierID", nil)
		}
		if cmd.ID == 0 {
			return cm.Error(cm.InvalidArgument, "Missing BrandID", nil)
		}
		{
			s2 := s.Table("supplier_brand").Where("id = ?", cmd.ID).Where("supplier_id = ?", cmd.SupplierID)

			if deleted, err := s2.Delete(&model.ProductBrand{}); err != nil {
				return err
			} else if deleted == 0 {
				return cm.Error(cm.NotFound, "", nil)
			}
		}
		if _, err := s.Table("product").
			Where("product_brand_id = ?", cmd.ID).
			UpdateMap(M{"product_brand_id": nil}); err != nil {
			return err
		}
		return nil
	})
}
