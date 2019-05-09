package model

type SummarizeFulfillmentsRequest struct {
	ShopID   int64
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
