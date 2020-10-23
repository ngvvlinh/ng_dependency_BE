package reportserver

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"time"

	wkhtml "github.com/SebastiaanKlippert/go-wkhtmltopdf"

	cm "o.o/backend/pkg/common"
)

const dateLayoutArg = `2006-01-02`

type fileType string

var (
	fileTypePDF   fileType = "pdf"
	fileTypeExcel fileType = "excel"
)

func getFileType(r *http.Request) fileType {
	fileTypeArg := r.URL.Query().Get("file_type")

	switch fileTypeArg {
	case string(fileTypePDF):
		return fileTypePDF
	case string(fileTypeExcel):
		return fileTypeExcel
	default:
		return fileType(fileTypeArg)
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

func exportFile(fileTyp fileType, w http.ResponseWriter, args ...interface{}) {
	switch fileTyp {
	case fileTypePDF:
		exportPDF(w, args[0].(bytes.Buffer))
	case fileTypeExcel:
		http.Error(w, "unsupported file type Excel", http.StatusBadRequest)
		return
	default:
		http.Error(w, fmt.Sprintf("unsupported file type %s", fileTyp), http.StatusBadRequest)
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
