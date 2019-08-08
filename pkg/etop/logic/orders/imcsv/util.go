package imcsv

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"

	"etop.vn/api/main/location"
	catalogsqlstore "etop.vn/backend/com/main/catalog/sqlstore"
	pbcm "etop.vn/backend/pb/common"
	pbsheet "etop.vn/backend/pb/common/spreadsheet"
	pborder "etop.vn/backend/pb/etop/order"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/cmsql"
	"etop.vn/backend/pkg/common/httpx"
	"etop.vn/backend/pkg/common/idemp"
	"etop.vn/backend/pkg/common/imcsv"
	"etop.vn/backend/pkg/common/redis"
	cmservice "etop.vn/backend/pkg/common/service"
	"etop.vn/backend/pkg/etop/model"
	"etop.vn/backend/pkg/etop/upload"
	"etop.vn/common/l"
)

var ll = l.New()
var idempgroup *idemp.RedisGroup
var uploader *upload.Uploader
var locationBus location.QueryBus
var shopVariantStore catalogsqlstore.ShopVariantStoreFactory

const PrefixIdemp = "IdempImportOrder"

func Init(_locationBus location.QueryBus, sd cmservice.Shutdowner, rd redis.Store, ul *upload.Uploader, db cmsql.Database) {
	locationBus = _locationBus
	idempgroup = idemp.NewRedisGroup(rd, PrefixIdemp, 5*60) // 5 minutes
	sd.Register(idempgroup.Shutdown)
	shopVariantStore = catalogsqlstore.NewShopVariantStore(db)

	if ul != nil {
		uploader = ul
		ul.ExpectDir(string(model.ImportTypeShopOrder))
	}
}

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
	GHNNoteCode string // postcheck_note

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

func (imp *Importer) generateErrorResponse(idx imcsv.Indexer, errs []error) (*pborder.ImportOrdersResponse, error) {
	resp := &pborder.ImportOrdersResponse{
		Data:       imp.toSpreadsheetData(idx),
		CellErrors: pbcm.PbErrors(errs),
	}
	return resp, nil
}

func (imp *Importer) toSpreadsheetData(idx imcsv.Indexer) *pbsheet.SpreadsheetData {
	return pbsheet.ToSpreadsheetData(imp.Schema, idx, imp.Rows, imp.LastRow)
}

func uploadFile(id int64, data []byte) (*upload.StoreFileCommand, error) {
	fileName := fmt.Sprintf("%v.xlsx", id)
	uploadCmd := &upload.StoreFileCommand{
		UploadType: string(model.ImportTypeShopOrder),
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

	if err := json.Unmarshal([]byte(debugHeader), &debug); err != nil {
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
