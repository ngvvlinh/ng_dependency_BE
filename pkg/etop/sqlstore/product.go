package sqlstore

import (
	"context"
	"time"

	"etop.vn/backend/com/main/catalog/convert"
	catalogmodel "etop.vn/backend/com/main/catalog/model"
	catalogmodelx "etop.vn/backend/com/main/catalog/modelx"
	catalogsqlstore "etop.vn/backend/com/main/catalog/sqlstore"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/capi/dot"
)

func init() {
	bus.AddHandlers("sql",
		GetShopVariant,
		RemoveShopVariants,
		UpdateShopVariant,

		RemoveShopProducts,
		UpdateShopProduct,
		UpdateShopProductsTags,
	)
}

var shopProductStore catalogsqlstore.ShopProductStoreFactory

func GetShopVariant(ctx context.Context, query *catalogmodelx.GetShopVariantQuery) error {
	if query.ShopID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing AccountID", nil)
	}

	var variant catalogmodel.ShopVariant
	if err := x.Where("shop_id = ? AND variant_id = ?", query.ShopID, query.VariantID).
		ShouldGet(&variant); err != nil {
		return err
	}
	query.Result = convert.ShopVariant(&variant)
	return nil
}

func updateOrInsertShopVariant(sv *catalogmodel.ShopVariant, x Qx) error {
	if sv.VariantID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing VariantID", nil)
	}
	if sv.ShopID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing ShopID", nil)
	}

	var shopVariant catalogmodel.ShopVariant
	ok, err := x.Table("shop_variant").
		Where("variant_id = ? AND shop_id = ?", sv.VariantID, sv.ShopID).
		Get(&shopVariant)
	if err != nil {
		return err
	}

	if !ok {
		return x.Table("shop_variant").ShouldInsert(sv)
	}

	return x.Table("shop_variant").
		Where("variant_id = ? AND shop_id = ?", sv.VariantID, sv.ShopID).
		ShouldUpdate(sv)
}

func UpdateShopVariant(ctx context.Context, cmd *catalogmodelx.UpdateShopVariantCommand) error {
	if cmd.ShopID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing AccountID", nil)
	}
	if cmd.Variant == nil || cmd.Variant.VariantID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing product", nil)
	}

	sv := *cmd.Variant
	if err := updateOrInsertShopVariant(&sv, x); err != nil {
		return err
	}

	query := &catalogmodelx.GetShopVariantQuery{
		ShopID:    cmd.ShopID,
		VariantID: cmd.Variant.VariantID,
	}
	if err := bus.Dispatch(ctx, query); err != nil {
		return cm.Error(cm.Internal, "", err)
	}

	cmd.Result = query.Result
	return nil
}

func RemoveShopVariants(ctx context.Context, cmd *catalogmodelx.RemoveShopVariantsCommand) error {
	if cmd.ShopID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing AccountID", nil)
	}

	updated, err := x.Table("shop_variant").
		Where("shop_id = ?", cmd.ShopID).
		In("variant_id", cmd.IDs).
		UpdateMap(map[string]interface{}{
			"deleted_at": time.Now(),
		})
	if err != nil {
		return err
	}

	cmd.Result.Removed = int(updated)
	return nil
}

func RemoveShopProducts(ctx context.Context, cmd *catalogmodelx.RemoveShopProductsCommand) error {
	if cmd.ShopID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing AccountID", nil)
	}

	if len(cmd.IDs) == 0 {
		return cm.Error(cm.InvalidArgument, "Misssing IDs", nil)
	}

	return inTransaction(func(x Qx) error {
		var productsCount uint64

		deletedCount, err := x.Table("shop_product").
			Where("shop_id = ?", cmd.ShopID).
			In("product_id", cmd.IDs).
			UpdateMap(map[string]interface{}{
				"deleted_at": time.Now(),
			})
		if err != nil {
			return nil
		}

		if _, err := x.Table("shop_variant").
			Where("shop_id = ?", cmd.ShopID).
			In("product_id", cmd.IDs).
			UpdateMap(map[string]interface{}{
				"deleted_at": time.Now(),
			}); err != nil {
			return err
		}

		cmd.Result.Removed = int(productsCount) + int(deletedCount)
		return nil
	})
}

func UpdateOrInsertShopProduct(sp *catalogmodel.ShopProduct) error {
	return updateOrInsertShopProduct(sp, x)
}

func updateOrInsertShopProduct(sp *catalogmodel.ShopProduct, x Qx) error {
	if sp.ProductID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing ProductID", nil)
	}

	if sp.ShopID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing AccountID", nil)
	}

	var shopProduct = new(catalogmodel.ShopProduct)
	if has, err := x.Table("shop_product").
		Where("product_id = ? AND shop_id = ?", sp.ProductID, sp.ShopID).
		Get(shopProduct); err != nil {
		return err
	} else if has {
		if err := x.Table("shop_product").
			Where("product_id = ? AND shop_id = ?", sp.ProductID, sp.ShopID).
			ShouldUpdate(sp); err != nil {
			return err
		}
		return nil
	}

	if err := x.Table("shop_product").ShouldInsert(sp); err != nil {
		return err
	}
	return nil
}

func UpdateShopProduct(ctx context.Context, cmd *catalogmodelx.UpdateShopProductCommand) error {
	if cmd.ShopID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing AccountID", nil)
	}
	if cmd.Product == nil || cmd.Product.ProductID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing product", nil)
	}

	sp := *cmd.Product
	if err := updateOrInsertShopProduct(&sp, x); err != nil {
		return err
	}
	{
		q := shopProductStore(ctx).ShopID(cmd.ShopID).ID(cmd.Product.ProductID)
		product, err := q.GetShopProductWithVariants()
		if err != nil {
			return err
		}
		cmd.Result = product
	}
	return nil
}

func UpdateShopProductsTags(ctx context.Context, cmd *catalogmodelx.UpdateShopProductsTagsCommand) error {
	if cmd.ShopID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing AccountID", nil)
	}
	req := cmd.Update
	if err := req.Verify(); err != nil {
		return err
	}

	var products []*catalogmodel.ShopProduct
	if err := x.Where("shop_id = ?", cmd.ShopID).
		In("product_id", cmd.ProductIDs).
		Find((*catalogmodel.ShopProducts)(&products)); err != nil {
		return err
	}

	productMap := make(map[dot.ID]*catalogmodel.ShopProduct)
	for _, p := range products {
		productMap[p.ProductID] = p
	}

	countUpdated := 0
	var savedError error
	for _, id := range cmd.ProductIDs {
		p := productMap[id]
		var pTag []string
		if p != nil {
			pTag = p.Tags
		}
		tags, tErr := model.PatchTag(pTag, *req)
		if tErr != nil {
			savedError = tErr
			continue
		}
		sp := &catalogmodel.ShopProduct{
			ShopID:    cmd.ShopID,
			ProductID: id,
			Tags:      tags,
		}

		if err := UpdateOrInsertShopProduct(sp); err != nil {
			savedError = err
			continue
		}
		countUpdated++
	}
	if countUpdated > 0 {
		cmd.Result.Updated = countUpdated
		return nil
	}
	if savedError != nil {
		return savedError
	}
	return cm.Error(cm.NotFound, "No product updated", nil)
}
