package imcsv

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
	"unicode"

	"o.o/api/top/int/types"
	pbsheet "o.o/api/top/int/types/spreadsheet"
	"o.o/api/top/types/etc/ghn_note_code"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/cmapi"
	"o.o/backend/pkg/common/apifw/httpx"
	"o.o/backend/pkg/common/apifw/idemp"
	"o.o/backend/pkg/common/imcsv"
	"o.o/backend/pkg/etop/model"
	"o.o/backend/pkg/etop/upload"
	"o.o/capi/dot"
	"o.o/common/jsonx"
	"o.o/common/l"
)

var ll = l.New()
var idempgroup *idemp.RedisGroup

const PrefixIdemp = "IdempImportOrder"

type Mode int

const (
	ModeOnlyOrder   Mode = 1
	ModeFulfillment Mode = 2

	ModeEtopCode Mode = 10
	ModeEdCode   Mode = 11
)

const (
	MaxFilesize   = 2 * 1024 * 1024 // 2 MB
	MaxRows       = 1200
	MaxCellErrors = 20
	MaxEmptyRows  = 20
)

func ParseMode(v string) (Mode, error) {
	switch v {
	case "only_order":
		return ModeOnlyOrder, nil
	case "fulfillment":
		return ModeFulfillment, nil
	default:
		return 0, cm.Errorf(cm.InvalidArgument, nil, "Invalid mode")
	}
}

func ParseCodeMode(v string) (Mode, error) {
	switch v {
	case "custom":
		return ModeEdCode, nil
	case "default":
		return ModeEtopCode, nil
	default:
		return 0, cm.Errorf(cm.InvalidArgument, nil, "Invalid code mode")
	}
}

type Importer struct {
	Schema      imcsv.Schema
	Mode        Mode
	CodeMode    Mode
	GHNNoteCode ghn_note_code.GHNNoteCode // postcheck_note

	Rows    [][]string
	LastRow int
	File    *multipart.FileHeader
}

func parseRequest(c *httpx.Context) (*Importer, error) {
	// Limit the max file size
	c.Req.Body = http.MaxBytesReader(c.Resp, c.Req.Body, MaxFilesize)

	mode, err := ParseMode(c.Req.FormValue("mode"))
	if err != nil {
		return nil, err
	}

	codeMode, err := ParseCodeMode(c.Req.FormValue("code_mode"))
	if err != nil {
		return nil, err
	}

	ghnNoteCode, err := parseAsGHNNoteCode(c.Req.FormValue("postcheck_note"))
	if err != nil {
		return nil, cm.Error(cm.InvalidArgument, err.Error(), err)
	}

	form, err := c.MultipartForm()
	if err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Invalid request")
	}

	files := form.File["files"]
	switch len(files) {
	case 0:
		return nil, cm.Errorf(cm.InvalidArgument, nil, "No file")
	case 1:
		// continue
	default:
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Too many files")
	}

	req := &Importer{
		Schema:      schema,
		Mode:        mode,
		CodeMode:    codeMode,
		GHNNoteCode: ghnNoteCode,
		File:        files[0],
	}
	return req, nil
}

func (imp *Importer) generateErrorResponse(idx imcsv.Indexer, errs []error) (*types.ImportOrdersResponse, error) {
	resp := &types.ImportOrdersResponse{
		Data:       imp.toSpreadsheetData(idx),
		CellErrors: cmapi.PbErrors(errs),
	}
	return resp, nil
}

func (imp *Importer) toSpreadsheetData(idx imcsv.Indexer) *pbsheet.SpreadsheetData {
	return imcsv.ToSpreadsheetData(imp.Schema, idx, imp.Rows, imp.LastRow)
}

func uploadFile(uploader *upload.Uploader, id dot.ID, data []byte) (*upload.StoreFileCommand, error) {
	fileName := fmt.Sprintf("%v.xlsx", id)
	uploadCmd := &upload.StoreFileCommand{
		UploadType: model.ImportTypeShopOrder.String(),
		FileName:   fileName,
		Data:       data,
	}

	if uploader == nil {
		ll.Warn("Disabled uploading file")
		return uploadCmd, nil
	}
	return uploadCmd, uploader.StoreFile(uploadCmd)
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

func isRandomFail(percent int) bool {
	var b [2]byte
	rand.Read(b[:])
	u16 := binary.LittleEndian.Uint16(b[:])
	return u16%100 < uint16(percent)
}

// cleanRows removes unicode characters which are not printable
func cleanRows(rows [][]string) {
	for i := range rows {
		for j := range rows[i] {
			rows[i][j] = cleanString(rows[i][j])
		}
	}
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

// For phone number, user may store it as Number instead of Text. For example:
// "976123456" instead of "0976358769", and Excel may stores it as
// "9.76123456E8". We have to parse it to number and convert again.
func convertExcelNumberToText(s string) string {
	s = strings.TrimSpace(s)
	if s == "" || s[0] == '0' { // fast path: starts with 0
		return s
	}
	for i := range s {
		c := s[i] // fast path: contains a character which can not be a number
		if !(c >= '0' && c <= '9' || c == '.' || c == 'e' || c == 'E') {
			return s
		}
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return s
	}
	return strconv.FormatFloat(f, 'f', 0, 64)
}
