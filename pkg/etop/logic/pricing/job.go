package pricing

import (
	"context"
	"fmt"
	"runtime/debug"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/l"
	"etop.vn/backend/pkg/etop/model"
)

func DoUpdatingPrice() (_err error) {
	ctx := context.Background()
	defer func() {
		e := recover()
		if e != nil {
			_err = cm.Error(cm.Internal, cm.F("%v", e), nil)
			debug.PrintStack()
			bus.PrintAllStack(ctx, true)
			return
		}
		if _err != nil {
			fmt.Printf("%+v", _err)
			if cm.IsTrace(_err) || cm.ErrorCode(_err) == cm.Internal {
				bus.PrintAllStack(ctx, true)
			}
		}
	}()

	suppliers, err := getMapSuppliers(ctx)
	if err != nil {
		ll.Error("Unable to query suppliers", l.Error(err))
		return err
	}

	rulesMap := make(map[int64]SupplierPriceRules)
	for id, s := range suppliers {
		if s.Rules == nil {
			ll.S.Warnf("No rules from supplier %v (id=%v)", s.Name, s.ID)
			continue
		}

		rulesMap[id], err = NewSupplierPriceRules(s.Rules)
		if err != nil {
			ll.S.Errorf("Error with price rule from supplier %v (id=%v): %v", s.Name, s.ID, err)
		}
	}

	query := &model.ScanVariantExternalsQuery{}
	if err := bus.Dispatch(ctx, query); err != nil {
		ll.Error("Unable to query products", l.Error(err))
		return err
	}

	scanned := 0
	updated := 0
	unchanged := 0
	errors := 0
	invalid := 0
	norule := 0
	for res := range query.Result {
		variants := res.Variants
		if len(variants) == 0 {
			ll.Error("No variants")
			continue
		}
		ll.S.Infof("Update variants from %v to %v (%v variants)", variants[0].Variant.ID, variants[len(variants)-1].Variant.ID, len(variants))

		for _, variant := range variants {
			v := variant.Variant
			vx := variant.VariantExternal

			scanned++

			// Ignore products with unknown suppliers
			if suppliers[v.SupplierID] == nil {
				ll.Error("Product with unknown supplier", l.Int64("id", v.ID))
				errors++
				continue
			}

			rules, ok := rulesMap[v.SupplierID]
			if !ok {
				norule++
			}

			// Work with both real rules and empty rules
			price, diff := rules.Apply(v, vx)
			if !price.IsValid() {
				invalid++
			}
			if !diff {
				unchanged++
				continue
			}

			cmd := &model.UpdateVariantPriceCommand{
				VariantID: v.ID,
				PriceDef:  price,
			}
			if err := bus.Dispatch(ctx, cmd); err != nil {
				ll.Error("Unable to update price", l.Int64("id", v.ID), l.Error(err))
				errors++
				continue
			}
			updated++
		}
	}

	ll.S.Infof("Updated %v products (%v scanned, %v unchanged, %v errors, %v invalid, %v norule)", updated, scanned, unchanged, errors, invalid, norule)
	return nil
}

func getMapSuppliers(ctx context.Context) (map[int64]*model.Supplier, error) {
	query := &model.GetAllSuppliersQuery{}
	if err := bus.Dispatch(ctx, query); err != nil {
		return nil, err
	}

	m := make(map[int64]*model.Supplier)
	for _, s := range query.Result {
		m[s.ID] = s
	}
	return m, nil
}
