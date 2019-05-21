package admin

import (
	"context"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/l"
)

func calcProductsPriceAndUpdate(ctx context.Context, products []*model.VariantExtended) error {
	updatePrices, err := calcProductsPriceWithDiff(ctx, products, true)
	if err != nil {
		return err
	}

	var errs cm.ErrorCollector
	for id, p := range updatePrices {
		cmd := &modelx3.UpdateVariantPriceCommand{
			VariantID: id,
			PriceDef:  p,
		}
		err := bus.Dispatch(ctx, cmd)
		errs.Collect(err)
		if err != nil {
			ll.Error("Unable to update price", l.Int64("id", id), l.Error(err))
			continue
		}
	}
	return errs.All()
}

func calcProductsPrice(ctx context.Context, variants []*model.VariantExtended) error {
	_, err := calcProductsPriceWithDiff(ctx, variants, false)
	return err
}

func getSupplierRules(ctx context.Context, products []*model.VariantExtended) (map[int64]*model.SupplierPriceRules, error) {
	if len(products) == 0 {
		return nil, nil
	}

	fallBack :=
		func(ids []int64) ([]*cache.SupplierPriceRulesWithID, error) {
			query := &modelx4.GetSuppliersQuery{IDs: ids}
			if err := bus.Dispatch(ctx, query); err != nil {
				return nil, err
			}
			res := make([]*cache.SupplierPriceRulesWithID, len(query.Result.Suppliers))
			for i, supplier := range query.Result.Suppliers {
				res[i] = &cache.SupplierPriceRulesWithID{
					ID:    supplier.ID,
					Rules: supplier.Rules,
				}
			}
			return res, nil
		}

	supplierIDs := make([]int64, len(products))
	for i, product := range products {
		supplierIDs[i] = product.SupplierID
	}

	query := &cache.GetSuppliersRulesQuery{
		SupplierIDs: supplierIDs,
		Fallback:    fallBack,
	}
	err := cache.GetSuppliersRules(ctx, query)
	if err != nil {
		ll.Error("Error with GetSuppliersRules in cache", l.Error(err))
		return nil, err
	}
	return query.Result.SupplierRules, nil
}

func calcProductsPriceWithDiff(ctx context.Context, products []*model.VariantExtended, withDiff bool) (map[int64]*model.PriceDef, error) {
	ruleMap, err := getSupplierRules(ctx, products)
	if err != nil {
		return nil, err
	}
	if ruleMap == nil {
		return nil, cm.Error(cm.InvalidArgument, "Can not get rules price", nil)
	}

	var updatePrices map[int64]*model.PriceDef
	for _, product := range products {
		// Only caclulate price when etop_status != 1
		if product.EtopStatus == model.StatusActive {
			continue
		}
		rules := ruleMap[product.SupplierID]
		if rules == nil {
			continue
		}

		newRule, err := pricing.NewSupplierPriceRules(rules)
		if err != nil {
			ll.Error("Error with price rule from supplier", l.Int64("supplier", product.SupplierID), l.Error(err))
			continue
		}
		p, diff := newRule.Apply(product.Variant, product.VariantExternal)
		p.ApplyTo(product.Variant)

		// prices need to update
		if withDiff && diff {
			if updatePrices == nil {
				updatePrices = make(map[int64]*model.PriceDef)
			}
			updatePrices[product.ID] = p
		}
	}
	return updatePrices, nil
}
