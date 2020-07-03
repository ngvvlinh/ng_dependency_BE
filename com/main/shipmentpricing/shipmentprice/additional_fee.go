package shipmentprice

import (
	"math"

	"o.o/api/main/shipmentpricing/shipmentprice"
	"o.o/api/top/types/etc/additional_fee_base_value"
	"o.o/api/top/types/etc/calculation_method"
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

func (args *CalcAdditionalFeeArgs) GetBaseValue(baseValueType additional_fee_base_value.BaseValueType) (int, error) {
	switch baseValueType {
	case additional_fee_base_value.CODAmount:
		return args.CODAmount, nil
	case additional_fee_base_value.MainFee:
		return args.MainFee, nil
	case additional_fee_base_value.BasketValue:
		return args.BasketValue, nil
	default:
		return 0, cm.Errorf(cm.InvalidArgument, nil, "base_value_type không hợp lệ")
	}
}

func calcAdditionalFee(args CalcAdditionalFeeArgs, addFee *shipmentprice.AdditionalFee) (*shipmentprice.ShippingFee, error) {
	if addFee == nil {
		return nil, nil
	}
	if !shipping_fee_type.Contain(args.AdditionalFeeTypes, addFee.FeeType) && !shipping_fee_type.Contain(acceptedAdditionalFeeTypes, addFee.FeeType) {
		return nil, nil
	}
	baseValue, err := args.GetBaseValue(addFee.BaseValueType)
	if err != nil {
		return nil, err
	}
	fee, err := applyFeeRule(addFee, baseValue)
	if err != nil {
		return nil, err
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
	switch addFee.CalculationMethod {
	case calculation_method.Cumulation:
		return calcFeeByCumulativeMethod(addFee.Rules, baseValue, addFee.FeeType.Name())
	case calculation_method.FirstSatisfy:
		return calcFeeByFirstSatisfyMethod(addFee.Rules, baseValue, addFee.FeeType.Name())
	default:
		return 0, cm.Errorf(cm.InvalidArgument, nil, "Phương pháp tính phí (calculation_method) không hợp lệ")
	}
}

func calcFeeByRule(rule *shipmentprice.AdditionalFeeRule, baseValue int, calculationMethod calculation_method.CalculationMethodType) (int, bool) {
	if rule.MinValue > baseValue {
		return 0, false
	}
	switch calculationMethod {
	case calculation_method.FirstSatisfy:
		if rule.MaxValue != -1 && rule.MaxValue < baseValue {
			return 0, false
		}
	case calculation_method.Cumulation:
		if rule.MaxValue != -1 && rule.MaxValue < baseValue {
			baseValue = rule.MaxValue
		}
	}

	amount := baseValue - rule.StartValue
	price := 0
	switch rule.PriceModifierType {
	case price_modifier_type.Percentage:
		fee := rule.Amount * float64(amount) / 100
		price = int(math.Round(fee))
	case price_modifier_type.FixedAmount:
		price = int(rule.Amount)
	}

	if price < rule.MinPrice {
		price = rule.MinPrice
	}
	return price, true
}

func calcFeeByCumulativeMethod(rules []*shipmentprice.AdditionalFeeRule, baseValue int, feeType string) (int, error) {
	if !checkValidFeeRules(baseValue, rules) {
		return 0, cm.Errorf(cm.FailedPrecondition, nil, "Không thể tính phí %v. Cấu hình không hợp lệ.", feeType)
	}

	price := 0
	applyRule := false
	for _, rule := range rules {
		subPrice, ok := calcFeeByRule(rule, baseValue, calculation_method.Cumulation)
		if !ok {
			continue
		}
		applyRule = true
		price += subPrice
	}
	if !applyRule {
		return 0, cm.Errorf(cm.FailedPrecondition, nil, "Không thể tính phí %v. Cấu hình không hợp lệ.", feeType)
	}

	return price, nil
}

func checkValidFeeRules(baseValue int, rules []*shipmentprice.AdditionalFeeRule) bool {
	maxValue := 0
	for _, rule := range rules {
		if rule.MaxValue == shipmentprice.MaximumValue {
			return true
		}
		if maxValue < rule.MaxValue {
			maxValue = rule.MaxValue
		}
	}
	if baseValue < maxValue {
		return true
	}
	return false
}

func calcFeeByFirstSatisfyMethod(rules []*shipmentprice.AdditionalFeeRule, baseValue int, feeType string) (int, error) {
	for _, rule := range rules {
		price, ok := calcFeeByRule(rule, baseValue, calculation_method.FirstSatisfy)
		if ok {
			return price, nil
		}
	}
	return 0, cm.Errorf(cm.FailedPrecondition, nil, "Không thể tính phí %v. Cấu hình không hợp lệ.", feeType)
}
