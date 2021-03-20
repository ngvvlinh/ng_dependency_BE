package fabo

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"o.o/backend/pkg/integration/shipping/ghn/driverv2"
)

func TestGetPromotionCoupon(t *testing.T) {
	driver := &FaboSupportedGHNDriver{}
	current := time.Date(2020, 11, 25, 0, 0, 0, 0, time.Local)

	validArgs := &driverv2.GetPromotionCouponArgs{
		FromProvinceCode: "79",
		CurrentTime:      current,
	}

	t.Run("Case valid", func(t *testing.T) {
		coupon, err := driver.GetPromotionCoupon(validArgs)
		assert.Nil(t, err)
		assert.EqualValues(t, coupon, promotionCoupons[0].Code)
	})

	invalidDateArgs1 := &driverv2.GetPromotionCouponArgs{
		FromProvinceCode: "79",
		CurrentTime:      time.Date(2020, 11, 19, 23, 15, 0, 0, time.Local),
	}
	invalidDateArgs2 := &driverv2.GetPromotionCouponArgs{
		FromProvinceCode: "79",
		CurrentTime:      time.Date(2020, 11, 30, 12, 15, 0, 0, time.Local),
	}

	t.Run("Case invalid date", func(t *testing.T) {
		_, err := driver.GetPromotionCoupon(invalidDateArgs1)
		assert.Errorf(t, err, "Đơn không thuộc khoản thời gian áp dụng coupon")

		_, err = driver.GetPromotionCoupon(invalidDateArgs2)
		assert.Errorf(t, err, "Đơn không thuộc khoản thời gian áp dụng coupon")
	})

	invalidLocationArgs := &driverv2.GetPromotionCouponArgs{
		FromProvinceCode: "01",
		CurrentTime:      current,
	}
	t.Run("Case invalid location", func(t *testing.T) {
		_, err := driver.GetPromotionCoupon(invalidLocationArgs)
		assert.Errorf(t, err, "Đơn không thuộc khu vực cho phép áp dụng coupon")
	})
}
