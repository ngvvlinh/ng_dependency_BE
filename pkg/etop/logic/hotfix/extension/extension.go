package extension

import (
	"bytes"
	"context"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize/v2"

	"o.o/api/etelecom"
	"o.o/api/main/identity"
	"o.o/api/top/types/etc/status3"
	com "o.o/backend/com/main"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/httpx"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/capi/dot"
)

type ExtensionService struct {
	etelecomAggr etelecom.CommandBus
	etelecomQS   etelecom.QueryBus
	dbEtelecom   *cmsql.Database
	identityQS   identity.QueryBus
}

func New(etelecomA etelecom.CommandBus, etelecomQ etelecom.QueryBus, dbEtelecom com.EtelecomDB, identityQ identity.QueryBus) *ExtensionService {
	return &ExtensionService{
		etelecomAggr: etelecomA,
		etelecomQS:   etelecomQ,
		dbEtelecom:   dbEtelecom,
		identityQS:   identityQ,
	}
}

func (s *ExtensionService) HandleImportExtension(c *httpx.Context) error {
	form, err := c.MultipartForm()
	if err != nil {
		return cm.Errorf(cm.InvalidArgument, nil, "Invalid request")
	}

	files := form.File["files"]
	switch len(files) {
	case 0:
		return cm.Errorf(cm.InvalidArgument, nil, "No file")
	case 1:
	// continue
	default:
		return cm.Errorf(cm.InvalidArgument, nil, "Too many files")
	}

	file, err := files[0].Open()
	if err != nil {
		return cm.Errorf(cm.InvalidArgument, err, "Can not read file")
	}
	defer file.Close()

	rawData, err := ioutil.ReadAll(file)
	if err != nil {
		return cm.Errorf(cm.InvalidArgument, err, "Không thể đọc được file.")
	}
	excelFile, err := excelize.OpenReader(bytes.NewReader(rawData))
	if err != nil {
		return cm.Errorf(cm.InvalidArgument, err, "Không thể đọc được file")
	}

	sheetName := excelFile.GetSheetName(0)
	rows, err := excelFile.GetRows(sheetName)
	if err != nil {
		return cm.Errorf(cm.InvalidArgument, err, "Không thể đọc được file")
	}
	if len(rows) <= 1 {
		return cm.Errorf(cm.InvalidArgument, nil, "File không có nội dung")
	}

	var lines []*Line
	ctx := bus.Ctx()
	for i, row := range rows {
		if i == 0 {
			continue
		}
		line, _err := s.parseRow(ctx, row)
		if _err != nil {
			return _err
		}
		lines = append(lines, line)
	}
	if err = s.createExtensions(ctx, lines); err != nil {
		return err
	}
	totalExt := len(lines)
	c.SetResult(map[string]interface{}{"code": "ok", "total": totalExt})
	return nil
}

type Line struct {
	HotlineID       dot.ID
	TenantID        dot.ID
	OwnerID         dot.ID
	AccountID       dot.ID
	ExtensionNumber string
	ExpiresAt       time.Time
}

var mapHotlines = make(map[string]*etelecom.Hotline)
var mapAccountUsers = make(map[dot.ID]*identity.AccountUser)
var dateLayouts = []string{
	// Date first, month later
	"02/01/2006", "02/01/06", "02-01-06",
	// Month first, date later
	"01/02/2006", "01/02/06", "01-02-06",
}

func (s *ExtensionService) parseRow(ctx context.Context, row []string) (*Line, error) {
	floatExt, err := strconv.ParseFloat(row[1], 32)
	if err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, err, "Can not parse extension number")
	}
	extNum := int(floatExt)
	extNumber := strconv.Itoa(extNum)

	hotlineNumber := row[2]
	expiresAtStr := row[3]

	var expiresAt time.Time
	for _, layout := range dateLayouts {
		expiresAt, err = time.ParseInLocation(layout, expiresAtStr, time.Local)
		if err != nil {
			continue
		}
		break
	}
	if expiresAt.IsZero() {
		return nil, cm.Errorf(cm.InvalidArgument, err, "ExpiresAt does not valid").WithMetap("row", row)
	}

	if extNumber == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Extension Number is empty").WithMetap("row", row)
	}
	if hotlineNumber == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Hotline number is empty").WithMetap("row", row)
	}

	hotline, ok := mapHotlines[hotlineNumber]
	if !ok {
		hotline, err = s.getHotlineByHotlineNumber(ctx, hotlineNumber)
		if err != nil {
			return nil, err
		}
		if hotline.Status != status3.P {
			return nil, cm.Errorf(cm.FailedPrecondition, nil, "Hotline was not activated").WithMetap("row", row)
		}
		mapHotlines[hotlineNumber] = hotline
	}

	var shopCode string
	if row[4] != "" {
		shopCode = strings.TrimSpace(row[4])
	}
	if shopCode == "" {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Missing shop code")
	}

	getShopQuery := &identity.GetShopByCodeQuery{
		Code: shopCode,
	}
	if err = s.identityQS.Dispatch(ctx, getShopQuery); err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Can not get shop by shop code")
	}
	shop := getShopQuery.Result

	return &Line{
		HotlineID:       hotline.ID,
		TenantID:        hotline.TenantID,
		OwnerID:         hotline.OwnerID,
		AccountID:       shop.ID,
		ExtensionNumber: extNumber,
		ExpiresAt:       expiresAt,
	}, nil
}

func (s *ExtensionService) createExtensions(ctx context.Context, lines []*Line) error {
	var importExtensions []*etelecom.ImportExtension
	for _, line := range lines {
		impExtension := &etelecom.ImportExtension{
			TenantID:        line.TenantID,
			OwnerID:         line.OwnerID,
			AccountID:       line.AccountID,
			ExtensionNumber: line.ExtensionNumber,
			ExpiresAt:       line.ExpiresAt,
			HotlineID:       line.HotlineID,
		}
		importExtensions = append(importExtensions, impExtension)
	}
	var cmd = &etelecom.ImportExtensionsCommand{
		ImportExtensions: importExtensions,
	}

	return s.etelecomAggr.Dispatch(ctx, cmd)
}

func (s *ExtensionService) getHotlineByHotlineNumber(ctx context.Context, hotlineNumber string) (*etelecom.Hotline, error) {
	queryHotline := &etelecom.GetHotlineByHotlineNumberQuery{
		Hotline: hotlineNumber,
	}
	if err := s.etelecomQS.Dispatch(ctx, queryHotline); err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, err, "Hotline không hợp lệ").WithMetap("hotline", hotlineNumber)
	}
	return queryHotline.Result, nil
}
