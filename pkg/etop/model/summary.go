package model

import "etop.vn/capi/dot"

type SummarizeFulfillmentsRequest struct {
	ShopID   dot.ID
	DateFrom string
	DateTo   string

	Result struct {
		Tables []*SummaryTable
	}
}

type SummarizePOSResquest struct {
	ShopID   dot.ID
	DateFrom string
	DateTo   string

	Result struct {
		Tables []*SummaryTable
	}
}

type SummaryTable struct {
	Label string
	Tags  []string
	Cols  []SummaryColRow
	Rows  []SummaryColRow
	Data  []SummaryItem
}

type SummaryColRow struct {
	Label  string
	Spec   string
	Unit   string
	Indent int
}

type SummaryItem struct {
	Spec  string
	Value int
	Unit  string
}
