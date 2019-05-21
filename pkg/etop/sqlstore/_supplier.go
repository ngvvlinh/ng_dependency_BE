package sqlstore

import (
	"context"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/l"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/backend/pkg/zdeprecated/supplier/modelx"
	modelx2 "etop.vn/backend/pkg/zdeprecated/sync/modelx"
)

func init() {
	bus.AddHandlers("sql",
		GetAllSuppliers,
		GetSupplier,
		GetSupplierExtended,
		GetSupplierExtendeds,
		GetSuppliersWithShipFromAddress,
		GetSuppliers,
		GetSupplierWithPermission,
		SyncUpdateCategories,
		SyncUpdateProducts,
		SyncUpdateProductsQuantity,
		UpdateKiotvietAccessToken,
		UpdatePriceRules,
		UpdateSupplier,
	)
}

func GetSupplier(ctx context.Context, query *modelx.GetSupplierQuery) error {
	supplier := new(model.Supplier)
	if err := x.Table("supplier").ShouldGet(supplier); err != nil {
		return err
	}
	query.Result = supplier
	return nil
}

// kiotviet := new(model.SupplierKiotviet)

// switch {
// case query.SupplierID != 0:
// 	ok, err := x.Table("supplier").
// 		Where("id = ?", query.SupplierID).
// 		Get(supplier)
// 	if err != nil {
// 		return err
// 	}
// 	if !ok {
// 		return model.ErrNotFound
// 	}

// 	ok, err = x.Table("supplier_kiotviet").
// 		Where("supplier_id = ?", supplier.ID).
// 		Get(kiotviet)
// 	if err != nil {
// 		return err
// 	}
// 	if !ok {
// 		ll.Error("Supplier without Kiotviet", l.Int64("SupplierID", query.SupplierID))
// 		kiotviet = nil
// 	}

// case query.KiotvietRetailerID != "":
// 	ok, err := x.Table("supplier_kiotviet").
// 		Where("retailer_id = ?", query.KiotvietRetailerID).
// 		Get(kiotviet)
// 	if err != nil {
// 		return err
// 	}
// 	if !ok {
// 		return model.ErrNotFound
// 	}

// 	ok, err = x.Table("supplier").
// 		Where("id = ?", kiotviet.SupplierID).
// 		Get(supplier)
// 	if err != nil {
// 		return err
// 	}
// 	if !ok {
// 		ll.Error("Kiotviet without supplier", l.Int64("KiotvietRetailerID", query.SupplierID))
// 		return cm.Error(cm.Internal, "Kiotviet not found", nil)
// 	}

// default:
// 	return cm.Error(cm.InvalidArgument, "Missing ID", nil)
// }

// query.Result.Supplier = supplier
// // query.Result.Kiotviet = kiotviet
// return nil

func GetSupplierExtended(ctx context.Context, query *modelx.GetSupplierExtendedQuery) error {
	supplier := new(model.SupplierExtended)
	if err := x.Table("supplier").
		Where("s.id = ?", query.SupplierID).
		ShouldGet(supplier); err != nil {
		return err
	}
	query.Result = supplier
	return nil
}

// switch {
// case query.SupplierID != 0:
// 	ok, err :=
// 	if err != nil {
// 		return err
// 	}
// 	if !ok {
// 		return model.ErrNotFound
// 	}

// ok, err = x.Table("supplier_kiotviet").
// 	Where("supplier_id = ?", supplier.ID).
// 	Get(kiotviet)
// if err != nil {
// 	return err
// }
// if !ok {
// 	ll.Error("Supplier without Kiotviet", l.Int64("SupplierID", query.SupplierID))
// 	kiotviet = nil
// }

// case query.KiotvietRetailerID != "":
// 	ok, err := x.Table("supplier_kiotviet").
// 		Where("retailer_id = ?", query.KiotvietRetailerID).
// 		Get(kiotviet)
// 	if err != nil {
// 		return err
// 	}
// 	if !ok {
// 		return model.ErrNotFound
// 	}

// 	ok, err = x.Table("supplier").
// 		Where("id = ?", kiotviet.SupplierID).
// 		Get(supplier)
// 	if err != nil {
// 		return err
// 	}
// 	if !ok {
// 		ll.Error("Kiotviet without supplier", l.Int64("KiotvietRetailerID", query.SupplierID))
// 		return cm.Error(cm.Internal, "Kiotviet not found", nil)
// 	}

// default:
// 	return cm.Error(cm.InvalidArgument, "Missing ID", nil)
// }

func GetSupplierWithPermission(ctx context.Context, query *modelx.GetSupplierWithPermissionQuery) error {
	if query.SupplierID == 0 || query.UserID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing required params", nil)
	}

	supplier := new(model.Supplier)
	if err := x.Where("id = ?", query.SupplierID).
		ShouldGet(supplier); err != nil {
		return err
	}
	query.Result.Supplier = supplier

	accUser := new(model.AccountUser)
	if err := x.
		Where("account_id = ? AND user_id = ?", query.SupplierID, query.UserID).
		ShouldGet(accUser); err != nil {
		return err
	}
	query.Result.Permission = accUser.Permission
	return nil
}

func GetSuppliers(ctx context.Context, query *modelx.GetSuppliersQuery) error {
	s := x.Table("supplier")
	if query.Status != nil {
		s = s.Where("status = ?", query.Status)
	}

	{
		s2 := s.Clone()
		s2, err := LimitSort(s2, query.Paging, Ms{"id": "", "created_at": "", "updated_at": "", "name": ""})
		if err != nil {
			return err
		}
		if query.IDs != nil {
			s2 = s2.In("id", query.IDs)
		}
		if err := s2.Find((*model.Suppliers)(&query.Result.Suppliers)); err != nil {
			return err
		}
	}
	{
		total, err := s.Count(&model.Supplier{})
		if err != nil {
			return err
		}
		query.Result.Total = int(total)
	}
	return nil
}

func GetSupplierExtendeds(ctx context.Context, query *modelx.GetSupplierExtendedsQuery) error {
	s := x.Table("supplier")
	if query.Status != nil {
		s = s.Where("status = ?", query.Status)
	}

	{
		s2 := s.Clone()
		s2, err := LimitSort(s2, query.Paging, Ms{"id": "", "created_at": "", "updated_at": "", "name": ""})
		if err != nil {
			return err
		}
		if query.IDs != nil {
			s2 = s2.In("s.id", query.IDs)
		}
		if err := s2.Find((*model.SupplierExtendeds)(&query.Result.Suppliers)); err != nil {
			return err
		}
	}
	{
		total, err := s.Count(&model.Supplier{})
		if err != nil {
			return err
		}
		query.Result.Total = int(total)
	}
	return nil
}

func GetSuppliersWithShipFromAddress(ctx context.Context, query *modelx.GetSuppliersWithShipFromAddressQuery) error {
	return x.In("s.id", query.IDs).
		Find((*model.SupplierShipFromAddresses)(&query.Result.Suppliers))
}

func GetAllSuppliers(ctx context.Context, query *modelx.GetAllSuppliersQuery) error {
	err := x.Find((*model.Suppliers)(&query.Result))
	return err
}

func UpdatePriceRules(ctx context.Context, cmd *modelx.UpdatePriceRulesCommand) error {
	if cmd.SupplierID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing SupplierID", nil)
	}

	m := &model.Supplier{
		Rules: cmd.PriceRules,
	}
	if err := x.Where("id = ?", cmd.SupplierID).
		ShouldUpdate(m); err != nil {
		return cm.ErrorTrace(cm.Unknown, "Can not update price rules", err)
	}
	return nil
}

func UpdateKiotvietAccessToken(ctx context.Context, cmd *modelx.UpdateKiotvietAccessTokenCommand) error {
	if cmd.ProductSourceID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing ProductSourceID", nil)
	}

	m := &model.ProductSourceInternalAccessToken{
		AccessToken: cmd.ClientToken,
		ExpiresAt:   cmd.ExpiresAt,
	}
	if err := x.
		Where("id = ?", cmd.ProductSourceID).
		Where("type = ?", model.TypeKiotviet).
		ShouldUpdate(m); err != nil {
		return err
	}
	return nil
}

func SyncUpdateCategories(ctx context.Context, cmd *modelx2.SyncUpdateCategoriesCommand) error {
	if cmd.SourceID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing SourceID", nil)
	}

	var source model.ProductSourceFtSupplier
	if has, err := x.
		Where("ps.id = ?", cmd.SourceID).
		Get(&source); err != nil {
		return err
	} else if !has {
		return cm.Error(cm.Internal, "ProductSource not found", err)
	}

	var updateOK, insertOK, updateErr, insertErr, skip, deleted int
	var errs cm.ErrorCollector
	for _, cx := range cmd.Data {
		if cx.ExternalID == "" {
			ll.Error("Missing ExternalID", l.Int64("source", cmd.SourceID), l.Object("cx", cx))
			continue
		}

		var currentCategoryExternal model.ProductSourceCategoryExternal
		if exists, err := x.
			Where("product_source_id = ?", cmd.SourceID).
			Where("product_source_type = ?", model.TypeKiotviet).
			Where("external_id = ?", cx.ExternalID).
			Get(&currentCategoryExternal); err != nil {
			errs.Collect(err)
			continue

		} else if exists {
			// Due to a bug from Kiotviet, where modifiedDate may be zero,
			// We ignore products without update_at / modifiedDate.
			if cx.ExternalUpdatedAt.IsZero() {
				skip++
				continue
			}

			cx.ID = 0
			cx.LastSyncAt = cmd.LastSyncAt

			err := x.Where("id = ?", currentCategoryExternal.ID).
				ShouldUpdate(cx)
			errs.Collect(err)
			if err != nil {
				updateErr++
			} else {
				updateOK++
			}

		} else {
			// cat.ID = cm.NewID()
			// cat.SupplierID = cmd.SupplierID
			// cat.ExternalCategory = c
			// cat.LastSyncAt = cmd.LastSyncAt

			err := x.InTransaction(func(x Qx) error {
				id := cm.NewID()
				c := &model.ProductSourceCategory{
					ID:                id,
					ProductSourceID:   cmd.SourceID,
					ProductSourceType: source.Type,
					SupplierID:        source.Supplier.ID,
					Name:              cx.ExternalName,
				}

				cx.ID = id
				cx.LastSyncAt = cmd.LastSyncAt
				cx.ProductSourceID = cmd.SourceID
				cx.ProductSourceType = source.Type
				return x.ShouldInsert(c, cx)
			})
			if err != nil {
				insertErr++
			} else {
				insertOK++
			}
			errs.Collect(err)
		}
	}

	/*
		if len(cmd.DeletedIDs) > 0 {
			del, err := x.
				Where("supplier_id = ? AND external_status != ?", cmd.SupplierID, model.StatusDisabled).
				In("external_id", cmd.DeletedIDs).
				Update(M{"external_status": model.StatusDisabled})
			if err != nil {
				ll.Error("Error deleting categories", l.Int64("supplier", cmd.SupplierID), l.Error(err))
			} else {
				deleted = int(del)
			}
		}
	*/

	if updateErr > 0 || insertErr > 0 {
		ll.S.Errorf("Sync categories to DB: inserted %v/%v (%v error), updated %v/%v (%v error)\n, skip %v, deleted %v - %v",
			insertOK, insertOK+insertErr, insertErr,
			updateOK, updateOK+updateErr, updateErr,
			skip, deleted, errs.Any())

	} else {
		ll.S.Infof("Sync categories to DB: inserted %v, updated %v, skip %v, deleted %v", insertOK, updateOK, skip, deleted)
	}

	if insertOK+updateOK > 0 {
		ps := &model.ProductSource{
			SyncStateCategories: cmd.SyncState,
		}
		if err := x.
			Where("id = ?", cmd.SourceID).
			ShouldUpdate(ps); err != nil {
			ll.Error("Error updating SyncState", l.Error(err))
		}
	}
	return errs.All()
}

func SyncUpdateProducts(ctx context.Context, cmd *modelx2.SyncUpdateProductsCommand) error {
	if cmd.SourceID == 0 {
		return cm.Error(cm.InvalidArgument, "Missing SourceID", nil)
	}

	var source model.ProductSourceFtSupplier
	if has, err := x.
		Where("ps.id = ?", cmd.SourceID).
		Get(&source); err != nil {
		return err
	} else if !has {
		return cm.Error(cm.Internal, "ProductSource not found", err)
	}

	var updateOK, insertOK, updateErr, insertErr, skip, deleted int
	var errs cm.ErrorCollector
	for _, item := range cmd.Data {
		vx := item.Variant
		if vx.ExternalID == "" || vx.ExternalProductID == "" {
			ll.Error("Missing ExternalID or ExternalProductID", l.Int64("source", cmd.SourceID), l.Object("vx", vx))
			errs.Collect(cm.Error(cm.InvalidArgument, "Missing ExternalID", nil))
			insertErr++
			continue
		}

		var currentVariantExternal model.VariantExternal
		if exists, err := x.
			Where("product_source_id = ?", cmd.SourceID).
			Where("product_source_type = ?", model.TypeKiotviet).
			Where("external_id = ?", vx.ExternalID).
			Get(&currentVariantExternal); err != nil {
			errs.Collect(err)
			continue

		} else if exists {
			// Due to a bug from Kiotviet, where modifiedDate may be zero,
			// We ignore products without update_at / modifiedDate.
			if vx.ExternalUpdatedAt.IsZero() {
				ll.Warn("The product is skipped: UpdatedAt is zero", l.Int64("id", vx.ID), l.String("x_id", vx.ExternalID))
				skip++
				continue
			}

			// Just to make sure we don't accidently overwrite the id
			vx.ID = 0
			vx.LastSyncAt = cmd.LastSyncAt

			err := x.InTransaction(func(x Qx) error {
				v := &model.VariantQuantity{
					QuantityOnHand:    item.QuantityOnHand,
					QuantityReserved:  item.QuantityReserved,
					QuantityAvailable: calcAvailable(item.QuantityOnHand, item.QuantityReserved),
				}
				if err := x.Where("id = ?", currentVariantExternal.ID).
					UpdateAll().
					ShouldUpdate(v); err != nil {
					return err
				}
				return x.Where("id = ?", currentVariantExternal.ID).
					ShouldUpdate(vx)

				// TODO(qv): Also broadcast to update shop products.
			})
			errs.Collect(err)
			if err != nil {
				updateErr++
			} else {
				updateOK++
			}

		} else {
			err := x.InTransaction(func(x Qx) error {
				prodCmd := &modelx2.SyncGetOrCreateProductCommand{
					SourceID:   cmd.SourceID,
					Variant:    vx,
					LastSyncAt: cmd.LastSyncAt,
				}
				if err := syncGetOrCreateProduct(ctx, x, prodCmd); err != nil {
					return err
				}

				id := cm.NewID()
				productID := prodCmd.Result.ProductID

				v := &model.Variant{
					ID:              id,
					ProductID:       productID,
					ProductSourceID: cmd.SourceID,
					// ProductSourceType: source.Type,
					SupplierID: source.Supplier.ID,
					// Name:        vx.ExternalName,
					ShortDesc:   vx.ExternalDescription,
					Description: vx.ExternalDescription,
					DescHTML:    "", // TODO(qv): Ignored
					Code:        vx.ExternalCode,

					QuantityOnHand:    item.QuantityOnHand,
					QuantityReserved:  item.QuantityReserved,
					QuantityAvailable: calcAvailable(item.QuantityOnHand, item.QuantityReserved),
				}
				if err := v.BeforeInsert(); err != nil {
					return err
				}

				vx.ID = id
				vx.LastSyncAt = cmd.LastSyncAt
				vx.ProductSourceID = cmd.SourceID
				vx.ProductSourceType = source.Type
				return x.ShouldInsert(v, vx)
			})
			errs.Collect(err)
			if err != nil {
				insertErr++
			} else {
				insertOK++
			}
		}
	}

	/*
		if len(cmd.DeletedIDs) > 0 {
			del, err := x.
				Where("supplier_id = ? AND external_status != ?", cmd.SupplierID, model.StatusDisabled).
				In("external_id", cmd.DeletedIDs).
				Update(M{"external_status": model.StatusDisabled})
			if err != nil {
				ll.Error("Error deleting products", l.Int64("supplier", cmd.SupplierID), l.Error(err))
			} else {
				deleted = int(del)
			}
		}
	*/

	if updateErr > 0 || insertErr > 0 {
		ll.S.Errorf("Source %v Sync products to DB: inserted %v/%v (%v error), updated %v/%v (%v error), skip %v, deleted %v - %v",
			cmd.SourceID,
			insertOK, insertOK+insertErr, insertErr,
			updateOK, updateOK+updateErr, updateErr,
			skip, deleted, errs.Any())
	} else {
		ll.S.Infof("Sync products to DB: inserted %v, updated %v, skip %v, deleted %v", insertOK, updateOK, skip, deleted)
	}

	if insertOK+updateOK > 0 {
		ps := &model.ProductSource{
			SyncStateProducts: cmd.SyncState,
		}
		err := x.Where("id = ?", cmd.SourceID).ShouldUpdate(ps)
		if err != nil {
			ll.Error("Error updating SyncState", l.Error(err))
		}
	}
	return errs.All()
}

func syncGetOrCreateProduct(ctx context.Context, x Qx, cmd *modelx2.SyncGetOrCreateProductCommand) error {
	return cm.ErrTODO
}

func SyncUpdateProductsQuantity(ctx context.Context, cmd *modelx2.SyncUpdateProductsQuantityCommand) error {
	ll.Info("Update quantity", l.Int64("source", cmd.SourceID), l.Any("updates", cmd.Updates))

	var productSource model.ProductSource
	if err := x.Table("product_source").
		Where("id = ?", cmd.SourceID).
		ShouldGet(&productSource); err != nil {
		return err
	}

	var defaultBranchID string
	if productSource.ExtraInfo == nil ||
		productSource.ExtraInfo.DefaultBranchID == "" {
		return cm.Error(cm.Internal, "Kiotviet has no default branch", nil)
	}

	skipped := 0
	var errs cm.ErrorCollector
	for _, update := range cmd.Updates {
		if update.ExternalProductID == "" || update.BranchID == "" {
			err := cm.Error(cm.InvalidArgument, "Missing ExternalID or Branch", nil)
			errs.Collect(err)
			continue
		}
		if update.BranchID != defaultBranchID {
			skipped++
			continue
		}

		variantQuantity := &model.VariantQuantity{
			QuantityAvailable: calcAvailable(update.QuantityOnHand, update.QuantityReserved),
			QuantityOnHand:    update.QuantityOnHand,
			QuantityReserved:  update.QuantityReserved,
		}

		if updated, err := x.
			Table("variant").
			Where("product_source_id = ?", cmd.SourceID).
			Where("external_id = ?", update.ExternalProductID).
			Update(variantQuantity); err != nil {
			errs.Collect(err)
		} else if updated == 0 {
			err := cm.Error(cm.NotFound, "", nil)
			errs.Collect(err)
		} else {
			// Increase total
			errs.Collect(nil)
		}
	}

	updated := errs.N() - errs.NErrors()
	ll.S.Infof("Updated quantity for product source %v: %v/%v (error %v, skip %v)",
		cmd.SourceID, updated, len(cmd.Updates), errs.NErrors(), skipped)
	return errs.All()
}

func calcAvailable(onHand, reserved int) int {
	return onHand - reserved - 3
}

func UpdateSupplier(ctx context.Context, cmd *modelx.UpdateSupplierCommand) error {
	supplier := cmd.Supplier
	if supplier.ID == 0 {
		cm.Error(cm.InvalidArgument, "Missing SupplierID", nil)
	}

	if err := x.Table("supplier").
		Where("id = ?", supplier.ID).
		ShouldUpdate(supplier); err != nil {
		return err
	}

	_supplier := new(model.SupplierExtended)
	s := x.
		Table("supplier").
		Where("s.id = ?", supplier.ID)

	has, err := s.Get(_supplier)
	if err != nil {
		return err
	}
	if !has {
		return cm.Error(cm.NotFound, "", nil)
	}

	cmd.Result = _supplier
	return nil
}
