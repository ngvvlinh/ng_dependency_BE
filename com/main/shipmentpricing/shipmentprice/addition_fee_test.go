package shipmentprice

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"o.o/api/main/shipmentpricing/shipmentprice"
	"o.o/api/top/types/etc/additional_fee_base_value"
	"o.o/api/top/types/etc/calculation_method"
	"o.o/api/top/types/etc/price_modifier_type"
	"o.o/api/top/types/etc/shipping_fee_type"
)

func TestCalcInsuranceFeeByFirstSatisfy(t *testing.T) {
	args := CalcAdditionalFeeArgs{
		AdditionalFeeTypes: []shipping_fee_type.ShippingFeeType{
			shipping_fee_type.Insurance,
		},
	}
	addFee := &shipmentprice.AdditionalFee{
		FeeType:           shipping_fee_type.Insurance,
		CalculationMethod: calculation_method.FirstSatisfy,
		BaseValueType:     additional_fee_base_value.BasketValue,
		Rules: []*shipmentprice.AdditionalFeeRule{
			{
				MinValue:          0,
				MaxValue:          50000,
				PriceModifierType: price_modifier_type.Percentage,
				Amount:            10.15,
				MinPrice:          3000,
			},
			{
				MinValue:          100001,
				MaxValue:          -1,
				PriceModifierType: price_modifier_type.Percentage,
				Amount:            20,
				MinPrice:          0,
				StartValue:        20000,
			},
			{
				MinValue:          50001,
				MaxValue:          100000,
				PriceModifierType: price_modifier_type.Percentage,
				Amount:            15,
				MinPrice:          0,
				StartValue:        30000,
			},
		},
	}
	var Cases = []struct {
		description        string
		basketValue        int
		insuranceFeeResult int
	}{
		{
			description:        "Giá trị đơn 20k",
			basketValue:        20000,
			insuranceFeeResult: 3000,
		}, {
			description:        "Giá trị đơn 50k",
			basketValue:        50000,
			insuranceFeeResult: 5075,
		}, {
			description:        "Giá trị đơn 100k",
			basketValue:        100000,
			insuranceFeeResult: 10500,
		}, {
			description:        "Giá trị đơn 150k",
			basketValue:        150000,
			insuranceFeeResult: 26000,
		},
	}
	for _, tt := range Cases {
		t.Run("TestCalcInsuranceFeeByFirstSatisfy", func(t *testing.T) {
			args.BasketValue = tt.basketValue
			feeLine, err := calcAdditionalFee(args, addFee)
			assert.NoError(t, err)
			if !assert.EqualValues(t, tt.insuranceFeeResult, feeLine.Price) {
				t.Fatalf("FAIL %v", tt.description)
			}
			t.Log("PASS", tt.description)
		})
	}
}

func TestCalcInsuranceFeeByFirstSatisfyOutOfRange(t *testing.T) {
	args := CalcAdditionalFeeArgs{
		AdditionalFeeTypes: []shipping_fee_type.ShippingFeeType{
			shipping_fee_type.Insurance,
		},
	}
	addFee := &shipmentprice.AdditionalFee{
		FeeType:           shipping_fee_type.Insurance,
		CalculationMethod: calculation_method.FirstSatisfy,
		BaseValueType:     additional_fee_base_value.BasketValue,
		Rules: []*shipmentprice.AdditionalFeeRule{
			{
				MinValue:          0,
				MaxValue:          50000,
				PriceModifierType: price_modifier_type.Percentage,
				Amount:            10.15,
				MinPrice:          3000,
			},
			{
				MinValue:          55000,
				MaxValue:          100000,
				PriceModifierType: price_modifier_type.Percentage,
				Amount:            15,
				MinPrice:          0,
				StartValue:        30000,
			},
		},
	}
	var OutOfRangeCases = []struct {
		description string
		basketValue int
	}{
		{
			description: "Giá trị đơn 51k",
			basketValue: 51000,
		},
		{
			description: "Giá trị đơn 150k",
			basketValue: 150000,
		},
	}
	for _, tt := range OutOfRangeCases {
		t.Run("TestCalcInsuranceFeeByFirstSatisfyOutOfRange", func(t *testing.T) {
			args.BasketValue = tt.basketValue
			_, err := calcAdditionalFee(args, addFee)
			assert.Error(t, err, "Không thể tính phí insurance. Cấu hình không hợp lệ")
		})
	}
}

func TestCalcInsuranceFeeByCummulation(t *testing.T) {
	args := CalcAdditionalFeeArgs{
		AdditionalFeeTypes: []shipping_fee_type.ShippingFeeType{
			shipping_fee_type.Insurance,
		},
	}
	addFee := &shipmentprice.AdditionalFee{
		FeeType:           shipping_fee_type.Insurance,
		CalculationMethod: calculation_method.Cumulation,
		BaseValueType:     additional_fee_base_value.BasketValue,
		Rules: []*shipmentprice.AdditionalFeeRule{
			{
				MinValue:          0,
				MaxValue:          50000,
				PriceModifierType: price_modifier_type.Percentage,
				Amount:            10.15,
				MinPrice:          3000,
			},
			{
				MinValue:          100001,
				MaxValue:          -1,
				PriceModifierType: price_modifier_type.Percentage,
				Amount:            20,
				MinPrice:          0,
				StartValue:        20000,
			},
			{
				MinValue:          50001,
				MaxValue:          100000,
				PriceModifierType: price_modifier_type.Percentage,
				Amount:            15,
				MinPrice:          0,
				StartValue:        30000,
			},
		},
	}
	var Cases = []struct {
		description        string
		basketValue        int
		insuranceFeeResult int
	}{
		{
			description:        "Giá trị đơn 20k",
			basketValue:        20000,
			insuranceFeeResult: 3000,
		}, {
			description:        "Giá trị đơn 50k",
			basketValue:        50000,
			insuranceFeeResult: 5075,
		}, {
			description:        "Giá trị đơn 100k",
			basketValue:        100000,
			insuranceFeeResult: 15575,
		}, {
			description:        "Giá trị đơn 150k",
			basketValue:        150000,
			insuranceFeeResult: 41575,
		},
	}
	for _, tt := range Cases {
		t.Run("TestCalcInsuranceFeeByCummulation", func(t *testing.T) {
			args.BasketValue = tt.basketValue
			feeLine, err := calcAdditionalFee(args, addFee)
			assert.NoError(t, err)
			if !assert.EqualValues(t, tt.insuranceFeeResult, feeLine.Price) {
				t.Fatalf("FAIL %v", tt.description)
			}
			t.Log("PASS", tt.description)
		})
	}
}

func TestCalcInsuranceFeeByCummulationOutOfRange(t *testing.T) {
	args := CalcAdditionalFeeArgs{
		AdditionalFeeTypes: []shipping_fee_type.ShippingFeeType{
			shipping_fee_type.Insurance,
		},
	}
	addFee := &shipmentprice.AdditionalFee{
		FeeType:           shipping_fee_type.Insurance,
		CalculationMethod: calculation_method.Cumulation,
		BaseValueType:     additional_fee_base_value.BasketValue,
		Rules: []*shipmentprice.AdditionalFeeRule{
			{
				MinValue:          0,
				MaxValue:          50000,
				PriceModifierType: price_modifier_type.Percentage,
				Amount:            10.15,
				MinPrice:          3000,
			},
			{
				MinValue:          55000,
				MaxValue:          100000,
				PriceModifierType: price_modifier_type.Percentage,
				Amount:            15,
				MinPrice:          0,
				StartValue:        30000,
			},
		},
	}
	var OutOfRangeCases = []struct {
		description string
		basketValue int
	}{
		{
			description: "Giá trị đơn 150k",
			basketValue: 150000,
		},
	}
	for _, tt := range OutOfRangeCases {
		t.Run("TestCalcInsuranceFeeByCummulationOutOfRange", func(t *testing.T) {
			args.BasketValue = tt.basketValue
			_, err := calcAdditionalFee(args, addFee)
			assert.Error(t, err, "Không thể tính phí insurance. Cấu hình không hợp lệ")
		})
	}
}
