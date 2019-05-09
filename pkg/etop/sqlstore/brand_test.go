package sqlstore

// import (
// 	"context"
// 	"testing"
// 	"time"

// 	"etop.vn/backend/pkg/common"
// 	"etop.vn/backend/pkg/common/bus"
// 	. "etop.vn/backend/pkg/common/testing"
// 	"etop.vn/backend/pkg/etop/model"

// 	. "github.com/smartystreets/goconvey/convey"
// 	"github.com/stretchr/testify/assert"
// )

// func TestBrand(t *testing.T) {
// 	ctx := context.Background()
// 	var supplierID int64 = 1234567890
// 	now := time.Now().In(time.UTC).Truncate(time.Second)

// 	Convey("Create", t, func() {
// 		Convey("Create (panic)", func() {
// 			So(func() {
// 				cmd := &model.CreateSupplierBrandCommand{}
// 				bus.Dispatch(ctx, cmd)
// 			}, ShouldPanic)
// 		})
// 		Convey("Create: Empty SupplierID (error)", func() {
// 			cmd := &model.CreateSupplierBrandCommand{
// 				Brand: &model.SupplierBrand{
// 					Name: "Sample",
// 				},
// 			}
// 			err := bus.Dispatch(ctx, cmd)
// 			So(err, ShouldCMError, cm.InvalidArgument, "Missing SupplierID")
// 		})

// 		brands := []*model.SupplierBrand{
// 			{
// 				SupplierID:  supplierID,
// 				Name:        "Sample 1",
// 				Description: "Sample description",
// 				CreatedAt:   now,
// 				UpdatedAt:   now,
// 			}, {
// 				SupplierID:  supplierID,
// 				Name:        "Sample 2",
// 				Description: "Sample description",
// 				CreatedAt:   now,
// 				UpdatedAt:   now,
// 			}, {
// 				SupplierID:  supplierID + 1,
// 				Name:        "Sample 3",
// 				Description: "Sample description",
// 				CreatedAt:   now,
// 				UpdatedAt:   now,
// 			},
// 		}
// 		for _, brand := range brands {
// 			cmd := &model.CreateSupplierBrandCommand{Brand: brand}
// 			err := bus.Dispatch(ctx, cmd)
// 			So(err, ShouldBeNil)

// 			result := cmd.Result
// 			So(result.ID, ShouldNotBeEmpty)
// 			So(result.SupplierID, ShouldEqual, brand.SupplierID)
// 			assert.WithinDuration(t, now, result.CreatedAt, time.Second)
// 			assert.WithinDuration(t, now, result.UpdatedAt, time.Second)

// 			result.CreatedAt = result.CreatedAt.Truncate(time.Second)
// 			result.UpdatedAt = result.UpdatedAt.Truncate(time.Second)
// 		}

// 		Reset(func() {
// 			MustExec(`TRUNCATE supplier_brand`)
// 		})

// 		brand := &model.SupplierBrandExtended{
// 			SupplierBrand: brands[1],
// 			Supplier: &model.Supplier{
// 				Name: "",
// 			},
// 		}
// 		So(brand.ID, ShouldNotEqual, 0)

// 		Convey("Get: Missing SupplierID (error)", func() {
// 			query := &model.GetSupplierBrandQuery{}
// 			err := bus.Dispatch(ctx, query)
// 			So(err, ShouldCMError, cm.InvalidArgument, "Missing ID")
// 		})
// 		Convey("Get: Missing ID (error)", func() {
// 			query := &model.GetSupplierBrandQuery{
// 				SupplierID: supplierID,
// 			}
// 			err := bus.Dispatch(ctx, query)
// 			So(err, ShouldCMError, cm.InvalidArgument, "Missing ID")
// 		})
// 		Convey("Get", func() {
// 			query := &model.GetSupplierBrandQuery{
// 				ID:         brand.ID,
// 				SupplierID: supplierID,
// 			}
// 			err := bus.Dispatch(ctx, query)
// 			So(err, ShouldBeNil)
// 			So(query.Result, ShouldResemble, brand)
// 		})
// 		Convey("Get: SupplierID Not Found (error)", func() {
// 			query := &model.GetSupplierBrandQuery{
// 				ID:         brand.ID,
// 				SupplierID: supplierID + 1,
// 			}
// 			err := bus.Dispatch(ctx, query)
// 			So(err, ShouldCMError, cm.NotFound, "")
// 		})
// 		Convey("List", func() {
// 			query := &model.GetSupplierBrandsQuery{
// 				SupplierID: supplierID,
// 			}
// 			err := bus.Dispatch(ctx, query)
// 			So(err, ShouldBeNil)
// 			So(query.Result.Brands, ShouldResembleSlice,
// 				[]*model.SupplierBrandExtended{
// 					&model.SupplierBrandExtended{brands[0], &model.Supplier{}},
// 					&model.SupplierBrandExtended{brands[1], &model.Supplier{}},
// 				})
// 		})
// 		Convey("Update: Empty Info (panic)", func() {
// 			So(func() {
// 				cmd := &model.UpdateSupplierBrandCommand{}
// 				bus.Dispatch(ctx, cmd)
// 			}, ShouldPanic)
// 		})
// 		Convey("Update: Missing BrandID (error)", func() {
// 			cmd := &model.UpdateSupplierBrandCommand{
// 				Brand: &model.SupplierBrand{
// 					SupplierID: supplierID,
// 				},
// 			}
// 			err := bus.Dispatch(ctx, cmd)
// 			So(err, ShouldCMError, cm.InvalidArgument, "Missing BrandID")
// 		})
// 		Convey("Update", func() {
// 			cmd := &model.UpdateSupplierBrandCommand{
// 				Brand: &model.SupplierBrand{
// 					ID:          brand.ID,
// 					SupplierID:  supplierID,
// 					Description: "test",
// 				},
// 			}
// 			err := bus.Dispatch(ctx, cmd)
// 			So(err, ShouldBeNil)

// 			Convey("Get again", func() {
// 				query := &model.GetSupplierBrandQuery{
// 					ID:         brand.ID,
// 					SupplierID: supplierID,
// 				}
// 				err := bus.Dispatch(ctx, query)
// 				So(err, ShouldBeNil)

// 				brand.Description = "test"
// 				brand.UpdatedAt = brand.UpdatedAt.Truncate(time.Second)
// 				So(query.Result, ShouldResemble, brand)
// 			})
// 		})
// 		Convey("Update: ID Not Found (error)", func() {
// 			cmd := &model.UpdateSupplierBrandCommand{
// 				Brand: &model.SupplierBrand{
// 					ID:          brand.ID + 1,
// 					SupplierID:  supplierID,
// 					Description: "test",
// 				},
// 			}
// 			err := bus.Dispatch(ctx, cmd)
// 			So(err, ShouldCMError, cm.NotFound, "")
// 		})
// 		Convey("Update: SupplierID Not Found (error)", func() {
// 			cmd := &model.UpdateSupplierBrandCommand{
// 				Brand: &model.SupplierBrand{
// 					ID:          brand.ID,
// 					SupplierID:  supplierID + 1,
// 					Description: "test",
// 				},
// 			}
// 			err := bus.Dispatch(ctx, cmd)
// 			So(err, ShouldCMError, cm.NotFound, "")
// 		})
// 		Convey("Delete: Missing SupplierID (error)", func() {
// 			cmd := &model.DeleteSupplierBrandCommand{
// 				ID: brand.ID,
// 			}
// 			err := bus.Dispatch(ctx, cmd)
// 			So(err, ShouldCMError, cm.InvalidArgument, "Missing SupplierID")
// 		})
// 		Convey("Delete", func() {
// 			cmd := &model.DeleteSupplierBrandCommand{
// 				ID:         brand.ID,
// 				SupplierID: supplierID,
// 			}
// 			err := bus.Dispatch(ctx, cmd)
// 			So(err, ShouldBeNil)

// 			Convey("Get again", func() {
// 				query := &model.GetSupplierBrandQuery{
// 					ID:         brand.ID,
// 					SupplierID: supplierID,
// 				}
// 				err := bus.Dispatch(ctx, query)
// 				So(err, ShouldCMError, cm.NotFound, "")
// 			})
// 		})
// 		Convey("Delete: Not Found (error)", func() {
// 			cmd := &model.DeleteSupplierBrandCommand{
// 				ID:         123456,
// 				SupplierID: 123123,
// 			}
// 			err := bus.Dispatch(ctx, cmd)
// 			So(err, ShouldCMError, cm.NotFound, "")
// 		})
// 	})
// }
