package reportserver

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	wkhtml "github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"golang.org/x/text/language"
	"golang.org/x/text/message"

	cm "o.o/backend/pkg/common"
)

var upperAlphabets = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}

const dateLayoutArg = `2006-01-02`

type fileType string

var (
	fileTypePDF     fileType = "pdf"
	fileTypeExcel   fileType = "excel"
	fileTypeUnknown fileType = "unknown"
)

func getFileType(r *http.Request) fileType {
	fileTypeArg := r.URL.Query().Get("file_type")

	switch fileTypeArg {
	case string(fileTypePDF):
		return fileTypePDF
	case string(fileTypeExcel):
		return fileTypeExcel
	default:
		return fileTypeUnknown
	}
}

func exportPDF(w http.ResponseWriter, html bytes.Buffer) {
	pdfg, err := wkhtml.NewPDFGenerator()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	pdfg.AddPage(wkhtml.NewPageReader(strings.NewReader(html.String())))
	pdfg.SetOutput(w)

	w.Header().Set("Content-Type", "application/pdf")

	if err := pdfg.Create(); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func exportExcel(w http.ResponseWriter, data ReportData) {
	excelFile := excelize.NewFile()

	defaultSheetName := "sheet1"

	var numCols = 1
	if len(data.ReportTable) != 0 && len(data.ReportTable[0]) != 0 {
		numCols = len(data.ReportTable[0])
	}

	// date
	excelFile.MergeCell(defaultSheetName, "A1", convertColName(numCols-1)+"1")
	excelFile.SetCellValue(defaultSheetName, "A1", data.ReportDate)
	dateStyleID, err := excelFile.NewStyle(`{"alignment": {"horizontal":"left"}, "font": {"family":"Calibri","size":9}}`)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	excelFile.SetCellStyle(defaultSheetName, "A1", "A1", dateStyleID)

	// name
	excelFile.MergeCell(defaultSheetName, "A2", convertColName(numCols-1)+"2")
	excelFile.SetCellValue(defaultSheetName, "A2", data.ReportName)
	nameStyleID, err := excelFile.NewStyle(`{"alignment": {"horizontal":"center"}, "font": {"family":"Calibri","size":17, "bold": true}}`)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	excelFile.SetCellStyle(defaultSheetName, "A2", "A2", nameStyleID)
	excelFile.SetRowHeight(defaultSheetName, 2, 30)

	// infos
	infoStyleID, err := excelFile.NewStyle(`{"alignment": {"horizontal":"center"}, "font": {"family":"Calibri","size":11}}`)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	for idx, reportInfo := range data.ReportInfos {
		excelFile.MergeCell(defaultSheetName, fmt.Sprintf("A%d", 3+idx), fmt.Sprintf("%s%d", convertColName(numCols-1), 3+idx))
		excelFile.SetCellValue(defaultSheetName, fmt.Sprintf("A%d", 3+idx), reportInfo)
		excelFile.SetCellStyle(defaultSheetName, fmt.Sprintf("A%d", 3+idx), fmt.Sprintf("A%d", 3+idx), infoStyleID)
	}

	// empty line
	currentRowIdx := 3 + len(data.ReportInfos)
	excelFile.MergeCell(defaultSheetName, fmt.Sprintf("A%d", currentRowIdx), fmt.Sprintf("%s%d", convertColName(numCols-1), currentRowIdx))

	headerStyleID, err := excelFile.NewStyle(`{"alignment": {"horizontal":"center"}, "font": {"family":"Calibri","size":11, "bold": true}}`)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// tables
	currentRowIdx += 1
	for rowIdx := range data.ReportTable {
		for colIdx := range data.ReportTable[rowIdx] {
			cell := data.ReportTable[rowIdx][colIdx]
			axis := fmt.Sprintf("%s%d", convertColName(colIdx), currentRowIdx+rowIdx)
			excelFile.SetCellValue(defaultSheetName, axis, cell.Value)
			if cell.StyleOption != "" {
				styleID, err := excelFile.NewStyle(cell.StyleOption)
				if err != nil {
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}
				excelFile.SetCellStyle(defaultSheetName, axis, axis, styleID)
			} else if cell.IsHeader {
				excelFile.SetCellStyle(defaultSheetName, axis, axis, headerStyleID)
			}

		}
	}

	w.Header().Add("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Add("Content-Disposition", "attachment; filename="+"report.xlsx")
	w.Header().Add("Cache-Control", "no-cache, no-store, must-revalidate")
	w.WriteHeader(http.StatusOK)
	if err := excelFile.Write(w); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func parseDateArg(date string) (time.Time, error) {
	res, err := time.ParseInLocation(dateLayoutArg, date, time.Local)
	if err != nil {
		return time.Time{}, cm.Errorf(cm.InvalidArgument, nil, "invalid date", err)
	}
	return res, nil
}

type Cell struct {
	IsHeader    bool
	StyleOption string
	Value       interface{}
}

type ReportData struct {
	ReportDate  string
	ReportName  string
	ReportInfos []string
	ReportTable [][]Cell
}

func convertColName(idx int) (result string) {
	for idx >= 0 {
		result = upperAlphabets[idx%26] + result
		idx = idx/26 - 1
	}
	return
}

func formatPrice(n int) string {
	p := message.NewPrinter(language.Vietnamese)
	return p.Sprint(n)
}
