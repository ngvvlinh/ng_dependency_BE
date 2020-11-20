package fabo

import (
	"context"
	"time"

	cm "o.o/backend/pkg/common"
	ghnclient "o.o/backend/pkg/integration/shipping/ghn/clientv2"
	"o.o/backend/pkg/integration/shipping/ghn/driverv2"
	"o.o/capi/dot"
	"o.o/capi/filter"
)

var _ driverv2.SupportedGHNDriver = &FaboSupportedGHNDriver{}

type FaboSupportedGHNDriver struct {
	client *ghnclient.Client
}

func NewFaboSupportedGHNDriver(env string, cfg ghnclient.GHNAccountCfg) *FaboSupportedGHNDriver {
	return &FaboSupportedGHNDriver{client: ghnclient.New(env, cfg)}
}

func (f *FaboSupportedGHNDriver) AddClientContract(ctx context.Context, clientID int) error {
	return f.client.AddClientContract(ctx, &ghnclient.AddClientContractRequest{ClientID: clientID})
}

type PromotionCoupon struct {
	Code        string
	Date        filter.Date
	Location    PromotionCouponLocation
	Description string
}

type PromotionCouponLocation struct {
	FromProvinceCodes []string
	FromDistrictCodes []string
}

var promotionCoupons = []PromotionCoupon{
	{
		Code: "FABO1357924681GHN",
		Date: filter.Date{
			From: dot.Time(time.Date(2020, 11, 24, 0, 0, 0, 0, time.Local)),
			To:   dot.Time(time.Date(2020, 11, 29, 23, 59, 0, 0, time.Local)),
		},
		Location: PromotionCouponLocation{
			FromProvinceCodes: []string{
				"79", // HCM
			},
		},
		Description: "Tất cả các KH đều áp dụng khi tạo đơn qua faboshop và chỉ áp dung cho đơn hàng từ HCM. Thời gian chạy: 24 - hết ngày 29/11 (6 ngày).",
	},
}

func (f *FaboSupportedGHNDriver) GetPromotionCoupon(args *driverv2.GetPromotionCouponArgs) (string, error) {
	current := args.CurrentTime
	for _, c := range promotionCoupons {
		if c.Date.From.ToTime().After(current) ||
			c.Date.To.ToTime().Before(current) {
			return "", cm.Errorf(cm.FailedPrecondition, nil, "Đơn không thuộc khoản thời gian áp dụng coupon")
			continue
		}
		allowFromProvinceCodes := c.Location.FromProvinceCodes
		if len(allowFromProvinceCodes) > 0 {
			if !cm.StringsContain(allowFromProvinceCodes, args.FromProvinceCode) {
				return "", cm.Errorf(cm.FailedPrecondition, nil, "Đơn không thuộc khu vực cho phép áp dụng coupon")
				continue
			}
		}
		return c.Code, nil
	}
	return "", nil
}
