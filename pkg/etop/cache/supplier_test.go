package cache

import (
	"context"
	"fmt"
	"testing"

	"etop.vn/backend/cmd/etop-server/config"
	"etop.vn/backend/pkg/common/redis"
	. "etop.vn/backend/pkg/common/testing"
	"etop.vn/backend/pkg/etop/model"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
)

func TestSupplierRules(t *testing.T) {
	ctx := context.Background()
	cfg := config.DefaultTest()
	store := redis.Connect(cfg.Redis.ConnectionString())
	ch, err := newRedisCache1(store, cfg.Postgres, true)
	if err != nil {
		panic(err)
	}

	reset := func() {
		err := redisStore.Del("srules:1000", "srules:1001", "srules:1002")
		assert.NoError(t, err)
	}
	reset()
	Convey("GetSuppliersRules", t, func() {
		ids := []int64{1000, 1001, 1002}
		called := make(map[int64]bool)
		fallback := func(ids []int64) ([]*SupplierPriceRulesWithID, error) {
			res := make([]*SupplierPriceRulesWithID, len(ids))
			for i, id := range ids {
				called[id] = true
				res[i] = &SupplierPriceRulesWithID{
					ID: id,
					Rules: &model.SupplierPriceRules{
						General: &model.SupplierPriceRule{
							Tag: fmt.Sprintf("rule %v", id),
						},
					},
				}
			}
			return res, nil
		}

		Convey("Get: The first time", func() {
			Reset(reset)

			query := &GetSuppliersRulesQuery{
				SupplierIDs: ids,
				Fallback:    fallback,
			}
			err := GetSuppliersRules(ctx, query)
			So(err, ShouldBeNil)
			So(called, ShouldDeepEqual,
				map[int64]bool{1000: true, 1001: true, 1002: true})
			res := query.Result.SupplierRules
			So(
				[]string{
					res[1000].General.Tag,
					res[1001].General.Tag,
					res[1002].General.Tag,
				}, ShouldResembleSlice,
				[]string{"rule 1000", "rule 1001", "rule 1002"},
			)

			Convey("Get: The second time", func() {
				called = make(map[int64]bool)
				query := &GetSuppliersRulesQuery{
					SupplierIDs: ids,
					Fallback:    fallback,
				}
				err := GetSuppliersRules(ctx, query)
				So(err, ShouldBeNil)
				So(called, ShouldDeepEqual, map[int64]bool{})
				So(
					[]string{
						res[1000].General.Tag,
						res[1001].General.Tag,
						res[1002].General.Tag,
					}, ShouldResembleSlice,
					[]string{"rule 1000", "rule 1001", "rule 1002"},
				)
			})

			Convey("Update rule 1001", func() {
				sql := fmt.Sprintf(
					`NOTIFY %v, '%v'`,
					EventSupplierRulesUpdate, ids[1])
				x.MustExec(sql)

				// Wait for the delete event to be acknowledged
				<-ch

				called = make(map[int64]bool)
				query := &GetSuppliersRulesQuery{
					SupplierIDs: ids,
					Fallback:    fallback,
				}
				err = GetSuppliersRules(ctx, query)
				So(err, ShouldBeNil)
				So(called, ShouldDeepEqual, map[int64]bool{1001: true})
				So(
					[]string{
						res[1000].General.Tag,
						res[1001].General.Tag,
						res[1002].General.Tag,
					}, ShouldResembleSlice,
					[]string{"rule 1000", "rule 1001", "rule 1002"},
				)
			})
		})
		Convey("get rule return nil", func() {
			Reset(reset)

			called2 := make(map[int64]bool)
			fallback2 := func(ids []int64) ([]*SupplierPriceRulesWithID, error) {
				res := make([]*SupplierPriceRulesWithID, len(ids))
				for i, id := range ids {
					if i == 1 {
						res[i] = nil
						continue
					}
					called2[id] = true

					res[i] = &SupplierPriceRulesWithID{
						ID: id,
						Rules: &model.SupplierPriceRules{
							General: &model.SupplierPriceRule{
								Tag: fmt.Sprintf("rule %v", id),
							},
						},
					}
				}
				return res, nil
			}
			query := &GetSuppliersRulesQuery{
				SupplierIDs: ids,
				Fallback:    fallback2,
			}
			err := GetSuppliersRules(ctx, query)
			So(err, ShouldBeNil)
			So(called2, ShouldDeepEqual,
				map[int64]bool{1000: true, 1002: true})
			res := query.Result.SupplierRules
			So(
				[]string{
					res[1000].General.Tag,
					res[1002].General.Tag,
				}, ShouldResembleSlice,
				[]string{"rule 1000", "rule 1002"},
			)
		})
	})
}
