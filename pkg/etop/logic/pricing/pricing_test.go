package pricing

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"etop.vn/backend/pkg/etop/model"
)

func TestApplySupplierRule(t *testing.T) {
	Convey("Price Rule", t, func() {
		rule := &model.SupplierPriceRule{
			ListPriceA:      1.1,
			ListPriceB:      10000,
			WholesalePriceA: 0.8,
			WholesalePriceB: 20000,
			RetailPriceMinA: 0.7,
			RetailPriceMinB: 5000,
			RetailPriceMaxA: 1.2,
			RetailPriceMaxB: 15000,
		}

		Convey("Imported product", func() {
			// prod := &model.Product{}
			v := &model.Variant{}
			vx := &model.VariantExternal{
				ExternalPrice: 100000,
			}
			p, _ := ApplySupplierRule(rule, v, vx)
			var listprice = 100000*1.1 + 10000
			So(p.ListPrice, ShouldEqual, listprice)

			wholesale := 100000*0.8 + 20000 + 0.1*listprice
			So(p.WholesalePrice, ShouldEqual, wholesale)

			So(p.RetailPriceMin, ShouldEqual, listprice*0.7+5000)
			So(p.RetailPriceMax, ShouldEqual, listprice*1.2+15000)
		})

		Convey("Zero external price", func() {
			v := &model.Variant{}
			vx := &model.VariantExternal{}
			p, _ := ApplySupplierRule(rule, v, vx)

			var listprice float64 = 10000
			So(p.ListPrice, ShouldEqual, listprice)

			var wholesale = 20000 + 0.1*listprice
			So(p.WholesalePrice, ShouldEqual, wholesale)

			So(p.RetailPriceMin, ShouldEqual, listprice*0.7+5000)
			So(p.RetailPriceMax, ShouldEqual, listprice*1.2+15000)
		})

		Convey("Supplier provided price", func() {
			v := &model.Variant{
				EdListPrice:      120000,
				EdWholesalePrice: 60000,
				EdRetailPriceMin: 80000,
				EdRetailPriceMax: 130000,
			}
			vx := &model.VariantExternal{
				ExternalPrice: 100000,
			}
			p, _ := ApplySupplierRule(rule, v, vx)

			var listprice float64 = 120000
			So(p.ListPrice, ShouldEqual, listprice)

			var wholesale = 60000 + 0.1*120000
			So(p.WholesalePrice, ShouldEqual, wholesale)

			So(p.RetailPriceMin, ShouldEqual, 80000)
			So(p.RetailPriceMax, ShouldEqual, 130000)
		})
	})

}
