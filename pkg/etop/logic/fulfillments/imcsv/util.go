package imcsv

import (
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"strings"
	"unicode"

	"o.o/api/top/int/types"
	pbsheet "o.o/api/top/int/types/spreadsheet"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/common/apifw/idemp"
	"o.o/backend/pkg/common/imcsv"
	"o.o/backend/pkg/common/validate"
	"o.o/backend/pkg/etop/model"
	"o.o/backend/pkg/etop/upload"
	"o.o/capi/dot"
	"o.o/common/jsonx"
	"o.o/common/l"
)

var ll = l.New()
var idempgroup *idemp.RedisGroup

const PrefixIdemp = "IdempImportFulfillment"

const (
	MaxFilesize = 2 * 1024 * 1024 // 2 MB
	MaxRows     = 1200
)

type Importer struct {
	Schema imcsv.Schema

	Rows    [][]string
	LastRow int
	File    *multipart.FileHeader
}

type Debug struct {
	FailPercent int `json:"fail_percent"`
}

func parseDebugHeader(h http.Header) (debug Debug, err error) {
	debugHeader := h.Get("debug")
	if debugHeader == "" {
		return
	}

	if err := jsonx.Unmarshal([]byte(debugHeader), &debug); err != nil {
		return Debug{}, cm.Error(cm.InvalidArgument, "Can not read debug header", err)
	}
	if debug.FailPercent < 0 || debug.FailPercent > 100 {
		return Debug{}, cm.Error(cm.InvalidArgument, "Invalid fail_percent", nil)
	}
	return
}

func (imp *Importer) generateErrorResponse(idx imcsv.Indexer, errs []error) (*types.ImportFulfillmentsResponse, error) {
	resp := &types.ImportFulfillmentsResponse{
		Data:       imp.toSpreadsheetData(idx),
		CellErrors: cmapi.PbErrors(errs),
	}
	return resp, nil
}

func (imp *Importer) toSpreadsheetData(idx imcsv.Indexer) *pbsheet.SpreadsheetData {
	return imcsv.ToSpreadsheetData(imp.Schema, idx, imp.Rows, imp.LastRow)
}

func (imp *Importer) toSpecificColumns(idx imcsv.Indexer) []*types.SpecificColumn {
	specificColumns := make([]*types.SpecificColumn, idx.Length())
	columns := imp.Schema

	for i, col := range columns {
		realIdx := idx.MapIndex(i)
		if realIdx < 0 {
			continue
		}
		specificColumns[realIdx] = &types.SpecificColumn{
			Fields: []string{col.Name},
			Label:  col.Display,
			Name:   col.Name,
		}
	}

	return specificColumns
}

func uploadFile(uploader *upload.Uploader, id dot.ID, data []byte) (*upload.StoreFileCommand, error) {
	fileName := fmt.Sprintf("%v.xlsx", id)
	uploadCmd := &upload.StoreFileCommand{
		UploadType: model.ImportTypeShopFulfillment.String(),
		FileName:   fileName,
		Data:       data,
	}

	if uploader == nil {
		ll.Warn("Disabled uploading file")
		return uploadCmd, nil
	}
	return uploadCmd, uploader.StoreFile(uploadCmd)
}

// cleanRows removes unicode characters which are not printable
func cleanRows(rows [][]string) {
	for i := range rows {
		for j := range rows[i] {
			rows[i][j] = cleanString(rows[i][j])
		}
	}
}

// clean empty rows
func cleanEmptyRows(rows [][]string) [][]string {
	var result [][]string

	for i := range rows {
		countEmptyCells := 0
		for j := range rows[i] {
			if strings.TrimSpace(rows[i][j]) == "" {
				countEmptyCells += 1
			}
		}
		if countEmptyCells != len(rows[i]) {
			result = append(result, rows[i])
		}
	}

	return result
}

func acceptCharacter(c rune) bool {
	return unicode.IsPrint(c) || c == ' ' || c == '\n'
}

func cleanString(s string) string {
	for _, c := range s {
		if !acceptCharacter(c) {
			return cleanString0(s)
		}
	}
	return s
}

func cleanString0(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for _, c := range s {
		switch {
		case unicode.IsPrint(c):
			b.WriteRune(c)
		case c == ' ' || c == '\n':
			b.WriteRune(c)
		case unicode.IsSpace(c):
			b.WriteRune(' ')
		default:
			// do not include in the result string
		}
	}
	return b.String()
}

func parseBool(v string) (bool, error) {
	switch validate.NormalizeSearchSimple(v) {
	case "x":
		return true, nil
	case "":
		return false, nil
	}
	return false, errors.New("Giá trị không hợp lệ, cần một trong các giá trị 'x' hoặc 'X'.")
}

func parseMoney(v string) (int, error) {
	var chars []rune
	for _, ch := range v {
		if !(ch == '.' || ch == ',' || ch == ' ' || ch == 'd' || ch == 'đ') {
			chars = append(chars, ch)
		}
	}

	return imcsv.ParseUint(string(chars))
}
