package reportserver

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"sort"
	"strconv"
	"time"

	"o.o/api/main/reporting"
	"o.o/api/top/types/etc/report_time_filter"
	"o.o/backend/pkg/etop/authorize/session"
)

type IncomeStatementData struct {
	Now      string
	ShopName string
	Cols     []IncomeStatementCol
}

type IncomeStatementCol struct {
	Header          string
	Revenue         int
	Discounts       int
	NetRevenue      int
	CostPrice       int
	GrossProfit     int
	Expenses        int
	DeliveryFee     int
	Discards        int
	Others          int
	ProfitStatement int
	OtherIncomes    int
	OtherExpenses   int
	NetProfit       int
}

var incomeStatementTmpl *template.Template

func init() {
	incomeStatementTmpl = parseTemplate("income_statement.html")
}

func (s *ReportService) ReportIncomeStatement(w http.ResponseWriter, r *http.Request) {
	data := s.getIncomesStatement(w, r)

	if err := incomeStatementTmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *ReportService) getIncomesStatement(w http.ResponseWriter, r *http.Request) (data IncomeStatementData) {
	ctx := r.Context()
	ss := session.GetSessionFromCtx(ctx)

	var (
		timeFilter report_time_filter.TimeFilter
		year       int64
		err        error
	)

	timeFilterArg := r.URL.Query().Get("time_filter")
	yearArg := r.URL.Query().Get("year")

	switch timeFilterArg {
	case report_time_filter.Month.String():
		timeFilter = report_time_filter.Month
	case report_time_filter.Quater.String():
		timeFilter = report_time_filter.Quater
	case report_time_filter.Year.String():
		timeFilter = report_time_filter.Year
	default:
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if yearArg != "" {
		year, err = strconv.ParseInt(yearArg, 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if year > int64(time.Now().Year()) {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
	}

	query := &reporting.ReportIncomeStatementQuery{
		ShopID:     ss.SS.Shop().ID,
		Year:       int(year),
		TimeFilter: timeFilter,
	}
	if err := s.ReportQuery.Dispatch(ctx, query); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result := query.Result
	var timeSeries []int
	for key := range result {
		timeSeries = append(timeSeries, key)
	}
	sort.Ints(timeSeries)
	{
		data.Cols = make([]IncomeStatementCol, len(timeSeries)+1)
		for i, key := range timeSeries {
			incomeStatementCol := result[key]
			data.Cols[i] = IncomeStatementCol{
				Header:          fmt.Sprintf("%s%d", timeFilter.ShortName(), key),
				Revenue:         incomeStatementCol.Revenue,
				Discounts:       incomeStatementCol.Discounts,
				NetRevenue:      incomeStatementCol.NetRevenue,
				CostPrice:       incomeStatementCol.CostPrice,
				GrossProfit:     incomeStatementCol.GrossProfit,
				Expenses:        incomeStatementCol.Expenses,
				DeliveryFee:     incomeStatementCol.ShippingFee,
				Discards:        incomeStatementCol.Discards,
				Others:          incomeStatementCol.Others,
				ProfitStatement: incomeStatementCol.ProfitStatement,
				OtherIncomes:    incomeStatementCol.OtherIncomes,
				OtherExpenses:   incomeStatementCol.OtherExpenses,
				NetProfit:       incomeStatementCol.NetProfit,
			}
		}

		// Add col "Tổng"
		idx := len(timeSeries)
		data.Cols[idx].Header = "Tổng"
		for _, key := range timeSeries {
			incomeStatementCol := result[key]
			data.Cols[idx].Revenue += incomeStatementCol.Revenue
			data.Cols[idx].Discounts += incomeStatementCol.Discounts
			data.Cols[idx].NetRevenue += incomeStatementCol.NetRevenue
			data.Cols[idx].CostPrice += incomeStatementCol.CostPrice
			data.Cols[idx].GrossProfit += incomeStatementCol.GrossProfit
			data.Cols[idx].Expenses += incomeStatementCol.Expenses
			data.Cols[idx].DeliveryFee += incomeStatementCol.ShippingFee
			data.Cols[idx].Discards += incomeStatementCol.Discards
			data.Cols[idx].Others += incomeStatementCol.Others
			data.Cols[idx].ProfitStatement += incomeStatementCol.ProfitStatement
			data.Cols[idx].OtherIncomes += incomeStatementCol.OtherIncomes
			data.Cols[idx].OtherExpenses += incomeStatementCol.OtherExpenses
			data.Cols[idx].NetProfit += incomeStatementCol.NetProfit
		}
	}

	data.Now = time.Now().Format("02/01/2006 15:04")
	data.ShopName = ss.SS.Shop().Name
	return
}

func (s *ReportService) ExportReportIncomeStatement(w http.ResponseWriter, r *http.Request) {
	fileTyp := getFileType(r)
	if fileTyp == "" {
		http.Error(w, "file_type is required", http.StatusBadRequest)
		return
	}

	data := s.getIncomesStatement(w, r)

	switch fileTyp {
	case fileTypePDF:
		var html bytes.Buffer
		if err := incomeStatementTmpl.Execute(&html, data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		exportPDF(w, html, &PDFOptions{
			width:  uint(100 + 30*len(data.Cols)),
			height: 297,
		})
	case fileTypeExcel:
		exportExcel(w, data.ToExcelData())
	default:
		http.Error(w, "định dạng file không hợp lệ", http.StatusBadRequest)
		return
	}

}

func (d IncomeStatementData) ToExcelData() (res ReportData) {
	res.ReportDate = fmt.Sprintf("Ngày lập: %s", d.Now)
	res.ReportName = "Báo cáo hoạt động kinh doanh"

	var infos []string
	infos = append(infos, fmt.Sprintf("Cửa hàng: %s", d.ShopName))
	res.ReportInfos = infos

	var table [][]Cell

	// header
	table = append(table,
		[]Cell{{}}, // empty cell
		[]Cell{{IsHeader: true, Value: "Doanh thu bán hàng (1)", StyleOption: styleLeftHeader}},
		[]Cell{{IsHeader: true, Value: "Giảm trừ Doanh thu (2)", StyleOption: styleLeftHeader}},
		[]Cell{{IsHeader: true, Value: "Doanh thu thuần (3=1-2)", StyleOption: styleLeftHeader}},
		[]Cell{{IsHeader: true, Value: "Giá vốn hàng bán (4)", StyleOption: styleLeftHeader}},
		[]Cell{{IsHeader: true, Value: "Lợi nhuận gộp về bán hàng (5=3-4)", StyleOption: styleLeftHeader}},
		[]Cell{{IsHeader: true, Value: "Chi phí (6 = 6.1 + 6.2 + 6.3)", StyleOption: styleLeftHeader}},
		[]Cell{{IsHeader: true, Value: "Phí giao hàng (6.1)", StyleOption: styleLeftHeader}},
		[]Cell{{IsHeader: true, Value: "Xuất hủy hàng hóa (6.2)", StyleOption: styleLeftHeader}},
		[]Cell{{IsHeader: true, Value: "Khác (6.3)", StyleOption: styleLeftHeader}},
		[]Cell{{IsHeader: true, Value: "Lợi nhuận từ hoạt động kinh doanh (7=5-6)", StyleOption: styleLeftHeader}},
		[]Cell{{IsHeader: true, Value: "Thu nhập khác (8)", StyleOption: styleLeftHeader}},
		[]Cell{{IsHeader: true, Value: "Chi phí khác (9)", StyleOption: styleLeftHeader}},
		[]Cell{{IsHeader: true, Value: "Lợi nhuận thuần (10=(7+8)-9)", StyleOption: styleLeftHeader}},
	)
	for _, val := range d.Cols {
		table[0] = append(table[0], Cell{IsHeader: true, Value: val.Header})
		table[1] = append(table[1], Cell{Value: val.Revenue})
		table[2] = append(table[2], Cell{Value: val.Discounts})
		table[3] = append(table[3], Cell{Value: val.NetRevenue})
		table[4] = append(table[4], Cell{Value: val.CostPrice})
		table[5] = append(table[5], Cell{Value: val.GrossProfit})
		table[6] = append(table[6], Cell{Value: val.Expenses})
		table[7] = append(table[7], Cell{Value: val.DeliveryFee})
		table[8] = append(table[8], Cell{Value: val.Discards})
		table[9] = append(table[9], Cell{Value: val.Others})
		table[10] = append(table[10], Cell{Value: val.ProfitStatement})
		table[11] = append(table[11], Cell{Value: val.OtherIncomes})
		table[12] = append(table[12], Cell{Value: val.OtherExpenses})
		table[13] = append(table[13], Cell{Value: val.NetProfit})
	}
	res.ReportTable = table
	return
}
