package pricing

import (
	"math"

	"etop.vn/backend/pkg/common/l"
	"etop.vn/backend/pkg/etop/model"
)

var ll = l.New()

const etopMakeupRatio = 0.1

type SupplierPriceRules struct {
	General *model.SupplierPriceRule
	MapXCat map[string]*model.SupplierPriceRule
}

func DefaultSupplierPriceRules() SupplierPriceRules {
	rule := model.DefaultSupplierPriceRule()
	return SupplierPriceRules{
		General: rule,
		MapXCat: make(map[string]*model.SupplierPriceRule),
	}
}

func NewSupplierPriceRules(rules *model.SupplierPriceRules) (SupplierPriceRules, error) {
	if err := rules.Validate(); err != nil {
		return DefaultSupplierPriceRules(), err
	}

	xcat := make(map[string]*model.SupplierPriceRule)
	for _, r := range rules.Rules {
		if r.ExternalCategoryID != "" {
			xcat[r.ExternalCategoryID] = r
		} else {
			ll.Error("Dropped rule", l.Object("r", r))
		}
	}

	return SupplierPriceRules{
		General: rules.General,
		MapXCat: xcat,
	}, nil
}

func (rules SupplierPriceRules) ApplyTo(v *model.Variant, vx *model.VariantExternal) {
	p, _ := rules.Apply(v, vx)
	p.ApplyTo(v)
}

// Apply calculate rule for given product
//
// Work with both real rules and empty rules
func (rules SupplierPriceRules) Apply(v *model.Variant, vx *model.VariantExternal) (*model.PriceDef, bool) {
	if rules.General == nil {
		return ApplySupplierRule(nil, v, vx)
	}

	rule := rules.General
	if r := rules.MapXCat[vx.ExternalCategoryID]; r != nil {
		rule = r
	}
	return ApplySupplierRule(rule, v, vx)
}

func ApplySupplierRule(rule *model.SupplierPriceRule, v *model.Variant, vx *model.VariantExternal) (*model.PriceDef, bool) {
	p := new(model.PriceDef)

	// p.ListPrice = calcPrice(
	// 	prod.ExternalPrice,
	// 	prod.SupplierListPrice,
	// 	rule.GetListPriceA(),
	// 	rule.GetListPriceB(),
	// )
	p.ListPrice = calcPrice(
		vx.ExternalPrice,
		v.EdListPrice,
		rule.GetListPriceA(),
		rule.GetListPriceB(),
	)

	// p.WholesalePrice0 = calcPrice(
	// 	prod.ExternalPrice,
	// 	prod.SupplierWholesalePrice,
	// 	rule.GetWholesalePriceA(),
	// 	rule.GetWholesalePriceB(),
	// )
	p.WholesalePrice0 = calcPrice(
		vx.ExternalPrice,
		v.EdWholesalePrice,
		rule.GetWholesalePriceA(),
		rule.GetWholesalePriceB(),
	)

	// Make up etop price = wholesale_price_0 + 10% * list_price
	p.WholesalePrice = calcMakeupPrice(p.WholesalePrice0, p.ListPrice)

	p.RetailPriceMin = calcPrice(
		p.ListPrice,
		// prod.SupplierRetailPriceMin,
		v.EdRetailPriceMin,
		rule.GetRetailPriceMinA(),
		rule.GetRetailPriceMinB(),
	)
	p.RetailPriceMax = calcPrice(
		p.ListPrice,
		// prod.SupplierRetailPriceMax,
		v.EdRetailPriceMax,
		rule.GetRetailPriceMaxA(),
		rule.GetRetailPriceMaxB(),
	)

	diff := p.ListPrice != v.ListPrice ||
		p.WholesalePrice0 != v.WholesalePrice0 ||
		p.WholesalePrice != v.WholesalePrice ||
		p.RetailPriceMin != v.RetailPriceMin ||
		p.RetailPriceMax != v.RetailPriceMax
	return p, diff
}

func calcPrice(base, overwrite int, a, b float64) int {
	var price float64
	if overwrite != 0 {
		price = float64(overwrite)
	} else {
		price = a*float64(base) + b
	}
	return round(price)
}

/*
	(0, 500k)     : +10%
	[500k, 1000k) : +7%
	[1000k, ∞)	  : +5%
*/
func calcMakeupPrice(wholesale, listprice int) int {
	switch {
	case listprice < 0:
		return 0

	case listprice < 500000:
		price := float64(wholesale) + float64(listprice)*0.1
		return round(price)

	case listprice < 1000000:
		price := float64(wholesale) + float64(listprice)*0.07
		return round(price)

	default:
		price := float64(wholesale) + float64(listprice)*0.05
		return round(price)
	}
}

func round(price float64) int {
	return int(math.Round(price/1000) * 1000)
}
