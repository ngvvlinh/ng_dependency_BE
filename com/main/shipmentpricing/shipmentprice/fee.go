package shipmentprice

import (
	"math"

	"o.o/api/main/shipmentpricing/shipmentprice"
	"o.o/api/top/types/etc/price_modifier_type"
	"o.o/api/top/types/etc/shipping_fee_type"
	cm "o.o/backend/pkg/common"
)

// acceptedAdditionalFeeTypes
//
// Luôn tính cước phí của các phụ phí trong mảng này
var acceptedAdditionalFeeTypes = []shipping_fee_type.ShippingFeeType{
	shipping_fee_type.Cods,
}

// calcAdditionalFees
//
// link desc: https://www.notion.so/c-t-B-sung-c-i-lo-i-ph-v-o-c-u-h-nh-gi-fa8a56e96bd0445bb46821d7125e9abd
func calcAdditionalFees(args CalcAdditionalFeeArgs, additionalFees []*shipmentprice.AdditionalFee) (res []*shipmentprice.ShippingFee, err error) {
	for _, addFee := range additionalFees {
		feeLine, err := calcAdditionalFee(args, addFee)
		if err != nil {
			return nil, err
		}
		if feeLine != nil {
			res = append(res, feeLine)
		}
	}
	return
}

type CalcAdditionalFeeArgs struct {
	BasketValue        int
	CODAmount          int
	MainFee            int
	AdditionalFeeTypes []shipping_fee_type.ShippingFeeType
}

func calcAdditionalFee(args CalcAdditionalFeeArgs, addFee *shipmentprice.AdditionalFee) (*shipmentprice.ShippingFee, error) {
	if addFee == nil {
		return nil, nil
	}
	if !shipping_fee_type.Contain(args.AdditionalFeeTypes, addFee.FeeType) && !shipping_fee_type.Contain(acceptedAdditionalFeeTypes, addFee.FeeType) {
		return nil, nil
	}
	var fee int
	var err error
	switch addFee.FeeType {
	case shipping_fee_type.Insurance:
		// Phí bảo hiểm
		// Giá trị tính theo % so với giá trị khai giá, có giá trị tối thiểu, thay đổi theo ngưỡng
		fee, err = applyFeeRule(addFee, args.BasketValue)
		if err != nil {
			return nil, err
		}
	case shipping_fee_type.Cods:
		// Phí thu hộ
		// Giá trị tính theo % so với giá trị thu hộ, có giá trị tối thiểu, thay đổi theo ngưỡng
		fee, err = applyFeeRule(addFee, args.CODAmount)
		if err != nil {
			return nil, err
		}
	case shipping_fee_type.Return,
		shipping_fee_type.Redelivery,
		shipping_fee_type.Adjustment:
		// Phí trả hàng, Phí kích hoạt giao lại, phí đổi thông tin
		// Giá trị tính theo % so với phí chính, có giá trị tối thiểu
		fee, err = applyFeeRule(addFee, args.MainFee)
	default:
		return nil, cm.Errorf(cm.FailedPrecondition, nil, "Loại cước phí không hợp lệ.")
	}
	if fee == 0 {
		return nil, nil
	}
	return &shipmentprice.ShippingFee{
		FeeType: addFee.FeeType,
		Price:   fee,
	}, nil
}

func applyFeeRule(addFee *shipmentprice.AdditionalFee, baseValue int) (int, error) {
	applyRule := false
	price := 0
	for _, rule := range addFee.Rules {
		if rule.MinValue > baseValue {
			continue
		}
		if rule.MaxValue != -1 && rule.MaxValue < baseValue {
			continue
		}
		applyRule = true
		switch rule.PriceModifierType {
		case price_modifier_type.Percentage:
			fee := rule.Amount * float64(baseValue) / 100
			price = int(math.Round(fee))
		case price_modifier_type.FixedAmount:
			price = int(rule.Amount)
		}

		if price < rule.MinPrice {
			price = rule.MinPrice
		}
		break
	}
	if !applyRule {
		return 0, cm.Errorf(cm.FailedPrecondition, nil, "Không thể tính phí %v. Cấu hình không hợp lệ.", addFee.FeeType.Name())
	}

	return price, nil
}
