package admin

import (
	"context"
	"testing"

	"etop.vn/backend/pkg/zdeprecated/supplier/modelx"

	"etop.vn/backend/cmd/etop-server/config"
	"etop.vn/backend/pkg/common/bus"
	"etop.vn/backend/pkg/common/redis"
	"etop.vn/backend/pkg/etop/cache"
	"etop.vn/backend/pkg/etop/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCalcProducts(t *testing.T) {
	ctx := context.Background()
	cfg := config.DefaultTest()
	store := redis.Connect(cfg.Redis.ConnectionString())
	_, connStr := cfg.Postgres.ConnectionString()
	cache.NewRedisCache(store, connStr)

	Convey("Test calculate price", t, func() {
		variants := []*model.VariantExtended{
			&model.VariantExtended{
				Product: &model.Product{
					Name:        "product_0",
					Description: "Do not calculate active products",
				},
				Variant: &model.Variant{
					SupplierID: 1001,
					EtopStatus: model.StatusActive,
				},
				VariantExternal: &model.VariantExternal{
					ExternalPrice: 10000,
				},
			},
			&model.VariantExtended{
				Product: &model.Product{
					Name:        "product_1",
					Description: "Only calculate unactive products",
				},
				Variant: &model.Variant{
					SupplierID: 1001,
					EtopStatus: model.StatusCreated,
				},
				VariantExternal: &model.VariantExternal{
					ExternalPrice: 10000,
				},
			},
			&model.VariantExtended{
				Product: &model.Product{
					Name:        "product_2",
					Description: "And disabled products",
				},
				Variant: &model.Variant{
					SupplierID: 1002,
					EtopStatus: model.StatusDisabled,
				},
				VariantExternal: &model.VariantExternal{
					ExternalPrice: 20000,
				},
			},
			&model.VariantExtended{
				Product: &model.Product{
					Name:        "product_3",
					Description: "Do not panic if there is no rule",
				},
				Variant: &model.Variant{
					SupplierID: 1003,
					EtopStatus: model.StatusDisabled,
				},
				VariantExternal: &model.VariantExternal{
					ExternalPrice: 20000,
				},
			},
			&model.VariantExtended{
				Product: &model.Product{
					Name:        "product_4",
					Description: "Do not panic if the rule is empty",
				},
				Variant: &model.Variant{
					SupplierID: 1004,
					EtopStatus: model.StatusDisabled,
				},
				VariantExternal: &model.VariantExternal{
					ExternalPrice: 20000,
				},
			},
			&model.VariantExtended{
				Product: &model.Product{
					Name:        "product_5",
					Description: "Do not panic if there is no supplier",
				},
				Variant: &model.Variant{
					SupplierID: 1005,
					EtopStatus: model.StatusDisabled,
				},
				VariantExternal: &model.VariantExternal{
					ExternalPrice: 20000,
				},
			},
		}

		bus.MockHandler(func(query *modelx.GetSuppliersQuery) error {
			suppliers := []*model.Supplier{
				{
					ID: 1001,
					Rules: &model.SupplierPriceRules{
						General: model.DefaultSupplierPriceRule(),
					},
				},
				{
					ID: 1002,
					Rules: &model.SupplierPriceRules{
						General: model.DefaultSupplierPriceRule(),
					},
				},
				{
					// Supplier 1003 has no rule
					ID:    1003,
					Rules: nil,
				},
				{
					// Supplier 1004 has empty rule
					ID:    1004,
					Rules: &model.SupplierPriceRules{},
				},
				// Supplier 1005 is not returned
			}
			suppliers[1].Rules.General.ListPriceA = 1.5
			query.Result.Suppliers = suppliers
			return nil
		})

		err := calcProductsPrice(ctx, variants)
		So(err, ShouldBeNil)
		So(variants[0].ListPrice, ShouldEqual, 0)
		So(variants[1].ListPrice, ShouldEqual, 10000)
		So(variants[2].ListPrice, ShouldEqual, 30000) // 20000x1.5
		So(variants[3].ListPrice, ShouldEqual, 0)
		So(variants[4].ListPrice, ShouldEqual, 0)
		So(variants[5].ListPrice, ShouldEqual, 0)
	})
}
