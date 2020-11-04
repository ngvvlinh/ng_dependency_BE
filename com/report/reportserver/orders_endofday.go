package reportserver

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"o.o/api/main/identity"
	"o.o/api/main/reporting"
	"o.o/backend/pkg/etop/authorize/session"
	"o.o/capi/dot"
)

var funcMap = template.FuncMap{
	"formatPrice": formatPrice,
}

type ReportService struct {
	ReportQuery   reporting.QueryBus
	IdentityQuery identity.QueryBus
}

type OrdersEndOfDayData struct {
	Now       string
	ShopName  string
	StaffName string
	CreatedAt string
	Summary   OrderEndOfDaySummary
	Lines     []OrderEndOfDayLine
}

type OrderEndOfDaySummary struct {
	TotalOrders   int
	TotalItems    int
	TotalAmount   int
	TotalFee      int
	TotalDiscount int
	TotalRevenue  int
}

type OrderEndOfDayLine struct {
	OrderCode     string
	CreatedAt     string
	TotalItems    int
	TotalAmount   int
	TotalFee      int
	TotalDiscount int
	Revenue       int
}

var tmpl *template.Template

func init() {
	tmpl = parseTemplate("orders_endofday.html")
}

func (s *ReportService) ReportOrdersEndOfDay(w http.ResponseWriter, r *http.Request) {
	data := s.getData(w, r)

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *ReportService) ExportReportOrdersEndOfDay(w http.ResponseWriter, r *http.Request) {
	fileTyp := getFileType(r)
	if fileTyp == "" {
		http.Error(w, "file_type is required", http.StatusBadRequest)
		return
	}

	data := s.getData(w, r)

	switch fileTyp {
	case fileTypePDF:
		var html bytes.Buffer
		if err := tmpl.Execute(&html, data); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		exportPDF(w, html)
	case fileTypeExcel:
		exportExcel(w, data.ToExcelData())
	default:
		http.Error(w, "định dạng file không hợp lệ", http.StatusBadRequest)
		return
	}
}

func (s *ReportService) getParams(ctx context.Context, r *http.Request) (staffName string, createdAt time.Time, createdBy dot.ID, err error) {
	createdAtArg := r.URL.Query().Get("created_at")
	createdByArg := r.URL.Query().Get("created_by")

	if createdAtArg != "" {
		createdAt, err = parseDateArg(createdAtArg)
		if err != nil {
			return "", time.Time{}, 0, err
		}
	} else {
		// current day
		now := time.Now()
		d := 24 * time.Hour
		createdAt = now.Truncate(d)
	}

	if createdByArg != "" {
		createdBy, err = dot.ParseID(createdByArg)
		if err != nil {
			return "", time.Time{}, 0, err
		}
		getUserQuery := &identity.GetUserByIDQuery{
			UserID: createdBy,
		}
		if err := s.IdentityQuery.Dispatch(ctx, getUserQuery); err != nil {
			return "", time.Time{}, 0, err
		}
		staffName = getUserQuery.Result.FullName
	}
	return
}

func (s *ReportService) getData(w http.ResponseWriter, r *http.Request) (data OrdersEndOfDayData) {
	ctx := r.Context()
	ss := session.GetSessionFromCtx(ctx)

	var (
		createdAt time.Time
		createdBy dot.ID
		staffName string
		err       error
	)

	staffName, createdAt, createdBy, err = s.getParams(ctx, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := &reporting.ReportOrdersQuery{
		ShopID:        ss.SS.Shop().ID,
		CreatedAtFrom: createdAt,
		CreatedAtTo:   createdAt.Add(24 * time.Hour),
		CreatedBy:     createdBy,
	}
	if err := s.ReportQuery.Dispatch(ctx, query); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data = OrdersEndOfDayData{
		Now:       time.Now().Format("02/01/2006 15:04"),
		CreatedAt: createdAt.Format("02/01/2006"),
		ShopName:  ss.SS.Shop().Name,
		StaffName: staffName,
		Summary: OrderEndOfDaySummary{
			TotalOrders: len(query.Result),
		},
	}

	for _, reportLine := range query.Result {
		data.Lines = append(data.Lines, OrderEndOfDayLine{
			OrderCode:     reportLine.OrderCode,
			CreatedAt:     reportLine.CreatedAt.Format("02/01/2006"),
			TotalItems:    reportLine.TotalItems,
			TotalAmount:   reportLine.TotalAmount,
			TotalFee:      reportLine.TotalFee,
			TotalDiscount: reportLine.TotalDiscount,
			Revenue:       reportLine.Revenue,
		})
		data.Summary.TotalItems += reportLine.TotalItems
		data.Summary.TotalAmount += reportLine.TotalAmount
		data.Summary.TotalFee += reportLine.TotalFee
		data.Summary.TotalDiscount += reportLine.TotalDiscount
		data.Summary.TotalRevenue += reportLine.Revenue
	}
	return
}

func (d OrdersEndOfDayData) ToExcelData() (res ReportData) {
	res.ReportDate = fmt.Sprintf("Ngày lập: %s", d.Now)
	res.ReportName = "Báo cáo cuối ngày về bán hàng"

	var infos []string
	infos = append(infos, fmt.Sprintf("Ngày bán: %s", d.CreatedAt))
	infos = append(infos, fmt.Sprintf("Cửa hàng: %s", d.ShopName))
	if d.StaffName != "" {
		infos = append(infos, fmt.Sprintf("Nhân viên: %s", d.StaffName))
	}
	res.ReportInfos = infos

	var table [][]Cell
	table = append(table, []Cell{
		{
			IsHeader: true,
			Value:    "Mã hoá đơn",
		}, {
			IsHeader: true,
			Value:    "Thời gian",
		}, {
			IsHeader: true,
			Value:    "SL sản phẩm",
		}, {
			IsHeader: true,
			Value:    "Doanh thu",
		}, {
			IsHeader: true,
			Value:    "Thu khác",
		}, {
			IsHeader: true,
			Value:    "Giảm giá",
		}, {
			IsHeader: true,
			Value:    "Thực thu",
		}})

	table = append(table, []Cell{
		{
			IsHeader: true,
			Value:    fmt.Sprintf("Hoá đơn: %d", len(d.Lines)),
		}, {
			Value: "",
		}, {
			IsHeader: true,
			Value:    d.Summary.TotalOrders,
		}, {
			IsHeader: true,
			Value:    d.Summary.TotalAmount,
		}, {
			IsHeader: true,
			Value:    d.Summary.TotalFee,
		}, {
			IsHeader: true,
			Value:    d.Summary.TotalDiscount,
		}, {
			IsHeader: true,
			Value:    d.Summary.TotalRevenue,
		}})

	for _, line := range d.Lines {
		table = append(table, []Cell{
			{
				Value: line.OrderCode,
			},
			{
				Value: line.CreatedAt,
			},
			{
				Value: line.TotalItems,
			},
			{
				Value: line.TotalAmount,
			},
			{
				Value: line.TotalFee,
			},
			{
				Value: line.TotalDiscount,
			},
			{
				Value: line.Revenue,
			},
		})
	}
	res.ReportTable = table

	return res
}
