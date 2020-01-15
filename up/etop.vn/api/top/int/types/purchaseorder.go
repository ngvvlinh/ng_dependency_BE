package types

import "etop.vn/common/jsonx"

type AdjustmentLine struct {
	Note   string `json:"note"`
	Amount int    `json:"amount"`
}

func (m *AdjustmentLine) Reset()         { *m = AdjustmentLine{} }
func (m *AdjustmentLine) String() string { return jsonx.MustMarshalToString(m) }

type DiscountLine struct {
	Note   string `json:"note"`
	Amount int    `json:"amount"`
}

func (m *DiscountLine) Reset()         { *m = DiscountLine{} }
func (m *DiscountLine) String() string { return jsonx.MustMarshalToString(m) }

type FeeLine struct {
	Note   string `json:"note"`
	Amount int    `json:"amount"`
}

func (m *FeeLine) Reset()         { *m = FeeLine{} }
func (m *FeeLine) String() string { return jsonx.MustMarshalToString(m) }
