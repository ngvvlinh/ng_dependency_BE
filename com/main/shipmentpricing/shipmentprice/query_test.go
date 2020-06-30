package shipmentprice

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"o.o/api/main/shipmentpricing/shipmentprice"
	"o.o/api/top/types/etc/price_modifier_type"
	"o.o/api/top/types/etc/shipping_fee_type"
)

func TestGetOverweightPrice(t *testing.T) {
	pricingDetailOV1 := &shipmentprice.PricingDetailOverweight{
		MinWeight:  100,
		MaxWeight:  500,
		WeightStep: 50,
		PriceStep:  2000,
	}
	pricingDetailOV2 := &shipmentprice.PricingDetailOverweight{
		MinWeight:  100,
		MaxWeight:  -1,
		WeightStep: 50,
		PriceStep:  2000,
	}
	var Cases = []struct {
		description string
		weight      int
		want1       int
		want2       int
	}{
		{
			description: "100g",
			weight:      100,
			want1:       0,
			want2:       0,
		}, {
			description: "200g",
			weight:      200,
			want1:       4000,
			want2:       4000,
		}, {
			description: "320g",
			weight:      320,
			want1:       10000,
			want2:       10000,
		}, {
			description: "500g",
			weight:      500,
			want1:       16000,
			want2:       16000,
		}, {
			description: "700g",
			weight:      700,
			want1:       16000,
			want2:       24000,
		},
	}
	for _, tt := range Cases {
		t.Run("Get Overweight Price 1", func(t *testing.T) {
			price := GetOverweightPrice(tt.weight, pricingDetailOV1)
			if !assert.EqualValues(t, tt.want1, price) {
				t.Fatalf("FAIL %v", tt.description)
			}
			t.Log("PASS", tt.description)
		})
		t.Run("Get Overweight Price 2", func(t *testing.T) {
			price := GetOverweightPrice(tt.weight, pricingDetailOV2)
			if !assert.EqualValues(t, tt.want2, price) {
				t.Fatalf("FAIL %v", tt.description)
			}
			t.Log("PASS", tt.description)
		})
	}
}

func TestGetPriceRuleDetail(t *testing.T) {
	var priceRuleDetails0 = []*shipmentprice.PricingDetail{}
	var priceRuleDetails1 = []*shipmentprice.PricingDetail{
		{
			Weight: 100,
			Price:  20000,
		}, {
			Weight: 300,
			Price:  25000,
		}, {
			Weight: 500,
			Price:  30000,
		}, {
			Weight: 3000,
			Price:  35000,
			Overweight: []*shipmentprice.PricingDetailOverweight{
				{
					MinWeight:  3000,
					MaxWeight:  -1,
					WeightStep: 500,
					PriceStep:  4500,
				},
			},
		},
	}
	var priceRuleDetails2 = []*shipmentprice.PricingDetail{
		{
			Weight: 100,
			Price:  27500,
		}, {
			Weight: 300,
			Price:  30000,
		}, {
			Weight: 500,
			Price:  34000,
		},
	}
	var priceRuleDetails3 = []*shipmentprice.PricingDetail{
		{
			Weight: 500,
			Price:  34000,
		}, {
			Weight: 100,
			Price:  27500,
		}, {
			Weight: 300,
			Price:  30000,
		},
	}
	var cases = []struct {
		description string
		weight      int
		want1       int // priceRuleDetails1
		want2       int // priceRuleDetails2
		want3       int // priceRuleDetails3 => want: same as case 2
	}{
		{
			description: "0g",
			weight:      0,
			want1:       100,
			want2:       100,
			want3:       100,
		}, {
			description: "100g",
			weight:      100,
			want1:       100,
			want2:       100,
			want3:       100,
		}, {
			description: "200g",
			weight:      200,
			want1:       300,
			want2:       300,
			want3:       300,
		}, {
			description: "300g",
			weight:      300,
			want1:       300,
			want2:       300,
			want3:       300,
		}, {
			description: "400g",
			weight:      400,
			want1:       500,
			want2:       500,
			want3:       500,
		}, {
			description: "500g",
			weight:      500,
			want1:       500,
			want2:       500,
			want3:       500,
		}, {
			description: "1200g",
			weight:      1200,
			want1:       3000,
			want2:       500,
			want3:       500,
		}, {
			description: "3000g",
			weight:      3000,
			want1:       3000,
			want2:       500,
			want3:       500,
		}, {
			description: "4300g",
			weight:      4300,
			want1:       3000,
			want2:       500,
			want3:       500,
		},
	}
	t.Run("Get PriceRuleDetail 0", func(t *testing.T) {
		_, err := GetPriceRuleDetail(100, priceRuleDetails0)
		assert.EqualError(t, err, "Missing priceDetails")
	})

	for _, tt := range cases {
		t.Run("Get PriceRuleDetail 1", func(t *testing.T) {
			priceRuleDetail, _ := GetPriceRuleDetail(tt.weight, priceRuleDetails1)
			if !assert.EqualValues(t, tt.want1, priceRuleDetail.Weight) {
				t.Fatalf("FAIL %v", tt.description)
			}
			t.Log("PASS", tt.description)
		})

		t.Run("Get PriceRuleDetail 2", func(t *testing.T) {
			priceRuleDetail, _ := GetPriceRuleDetail(tt.weight, priceRuleDetails2)
			if !assert.EqualValues(t, tt.want2, priceRuleDetail.Weight) {
				t.Fatalf("FAIL %v", tt.description)
			}
			t.Log("PASS", tt.description)
		})

		t.Run("Get PriceRuleDetail 3", func(t *testing.T) {
			priceRuleDetail, _ := GetPriceRuleDetail(tt.weight, priceRuleDetails3)
			if !assert.EqualValues(t, tt.want3, priceRuleDetail.Weight) {
				t.Fatalf("FAIL %v", tt.description)
			}
			t.Log("PASS", tt.description)
		})
	}
}

func TestGetPriceByPricingDetail(t *testing.T) {
	ruleDetail1 := &shipmentprice.PricingDetail{
		Weight: 300,
		Price:  26000,
	}
	ruleDetail2 := &shipmentprice.PricingDetail{
		Weight: 500,
		Price:  20000,
		Overweight: []*shipmentprice.PricingDetailOverweight{
			{
				MinWeight:  500,
				MaxWeight:  4000,
				WeightStep: 500,
				PriceStep:  5000,
			},
			{
				MinWeight:  4000,
				MaxWeight:  -1,
				WeightStep: 500,
				PriceStep:  7000,
			},
		},
	}
	ruleDetail3 := &shipmentprice.PricingDetail{
		Weight: 1000,
		Price:  35000,
		Overweight: []*shipmentprice.PricingDetailOverweight{
			{
				MinWeight:  4000,
				MaxWeight:  -1,
				WeightStep: 500,
				PriceStep:  7000,
			},
			{
				MinWeight:  1000,
				MaxWeight:  4000,
				WeightStep: 500,
				PriceStep:  5000,
			},
		},
	}

	// case not cover max weight (in overweight)
	// expect does not apply this price rule when it happens
	ruleDetail4 := &shipmentprice.PricingDetail{
		Weight: 1000,
		Price:  35000,
		Overweight: []*shipmentprice.PricingDetailOverweight{
			{
				MinWeight:  1000,
				MaxWeight:  4000,
				WeightStep: 500,
				PriceStep:  5000,
			},
		},
	}

	var cases = []struct {
		description string
		weight      int
		want1       int // ruleDetail1
		want2       int // ruleDetail2
		want3       int // ruleDetail3
		want4       int // ruleDetail4
	}{
		{
			description: "0g",
			weight:      0,
			want1:       26000,
			want2:       20000,
			want3:       35000,
			want4:       35000,
		}, {
			description: "100g",
			weight:      100,
			want1:       26000,
			want2:       20000,
			want3:       35000,
			want4:       35000,
		}, {
			description: "200g",
			weight:      200,
			want1:       26000,
			want2:       20000,
			want3:       35000,
			want4:       35000,
		}, {
			description: "600g",
			weight:      600,
			want1:       0,
			want2:       25000,
			want3:       35000,
			want4:       35000,
		}, {
			description: "1200g",
			weight:      1200,
			want1:       0,
			want2:       30000,
			want3:       40000,
			want4:       40000,
		}, {
			description: "1800g",
			weight:      1800,
			want1:       0,
			want2:       35000,
			want3:       45000,
			want4:       45000,
		}, {
			description: "3000g",
			weight:      3000,
			want1:       0,
			want2:       45000,
			want3:       55000,
			want4:       55000,
		},
		{
			description: "5400g",
			weight:      5400,
			want1:       0,
			want2:       76000,
			want3:       86000,
			want4:       0,
		},
	}
	for _, tt := range cases {
		t.Run("Get PriceByPricingDetail 1", func(t *testing.T) {
			price, _ := GetPriceByPricingDetail(tt.weight, ruleDetail1)
			if !assert.EqualValues(t, tt.want1, price) {
				t.Fatalf("FAIL %v", tt.description)
			}
			t.Log("PASS", tt.description)
		})
		t.Run("Get PriceByPricingDetail 2", func(t *testing.T) {
			price, _ := GetPriceByPricingDetail(tt.weight, ruleDetail2)
			if !assert.EqualValues(t, tt.want2, price) {
				t.Fatalf("FAIL %v", tt.description)
			}
			t.Log("PASS", tt.description)
		})
		t.Run("Get PriceByPricingDetail 3", func(t *testing.T) {
			price, _ := GetPriceByPricingDetail(tt.weight, ruleDetail3)
			if !assert.EqualValues(t, tt.want3, price) {
				t.Fatalf("FAIL %v", tt.description)
			}
			t.Log("PASS", tt.description)
		})

		t.Run("Get PriceByPricingDetail 4", func(t *testing.T) {
			price, _ := GetPriceByPricingDetail(tt.weight, ruleDetail4)
			if !assert.EqualValues(t, tt.want4, price) {
				t.Fatalf("FAIL %v", tt.description)
			}
			t.Log("PASS", tt.description)
		})
	}
}

func TestCalcInsuranceFees(t *testing.T) {
	args := CalcAdditionalFeeArgs{
		AdditionalFeeTypes: []shipping_fee_type.ShippingFeeType{
			shipping_fee_type.Insurance,
		},
	}
	addFee := &shipmentprice.AdditionalFee{
		FeeType: shipping_fee_type.Insurance,
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
			},
			{
				MinValue:          50001,
				MaxValue:          100000,
				PriceModifierType: price_modifier_type.Percentage,
				Amount:            15,
				MinPrice:          0,
			},
			{
				MinValue:          50001,
				MaxValue:          100000,
				PriceModifierType: price_modifier_type.Percentage,
				Amount:            15,
				MinPrice:          0,
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
			insuranceFeeResult: 15000,
		}, {
			description:        "Giá trị đơn 150k",
			basketValue:        150000,
			insuranceFeeResult: 30000,
		},
	}
	for _, tt := range Cases {
		t.Run("Calc Insurance Fee", func(t *testing.T) {
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
