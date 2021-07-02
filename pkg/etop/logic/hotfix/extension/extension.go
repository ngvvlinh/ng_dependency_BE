package extension

import (
	"bytes"
	"context"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"o.o/api/etelecom"
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
}

func New(etelecomA etelecom.CommandBus, etelecomQ etelecom.QueryBus, dbEtelecom com.EtelecomDB) *ExtensionService {
	return &ExtensionService{
		etelecomAggr: etelecomA,
		etelecomQS:   etelecomQ,
		dbEtelecom:   dbEtelecom,
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
	ownerID, err := dot.ParseID(GetFormValue(form.Value["owner_id"]))
	if err != nil {
		return cm.Errorf(cm.InvalidArgument, nil, "owner ID does not valid")
	}
	accountID, err := dot.ParseID(GetFormValue(form.Value["account_id"]))
	if err != nil {
		return cm.Errorf(cm.InvalidArgument, nil, "account ID does not valid")
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
		line, _err := s.parseRow(ctx, row, ownerID, accountID)
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
var dateLayouts = []string{
	// Date first, month later
	"02/01/2006", "02/01/06", "02-01-06",
	// Month first, date later
	"01/02/2006", "01/02/06", "01-02-06",
}

func (s *ExtensionService) parseRow(ctx context.Context, row []string, ownerID, accountID dot.ID) (*Line, error) {
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
		queryHotline := &etelecom.GetHotlineByHotlineNumberQuery{
			Hotline: hotlineNumber,
			OwnerID: ownerID,
		}
		if err = s.etelecomQS.Dispatch(ctx, queryHotline); err != nil {
			return nil, cm.Errorf(cm.InvalidArgument, err, "Hotline không hợp lệ").WithMetap("hotline", hotlineNumber).WithMetap("owner_id", ownerID)
		}
		hotline = queryHotline.Result
		mapHotlines[hotlineNumber] = hotline
	}

	return &Line{
		HotlineID:       hotline.ID,
		TenantID:        hotline.TenantID,
		OwnerID:         ownerID,
		AccountID:       accountID,
		ExtensionNumber: extNumber,
		ExpiresAt:       expiresAt,
	}, nil
}

func GetFormValue(ss []string) string {
	if ss == nil {
		return ""
	}
	return ss[0]
}

func (s *ExtensionService) createExtensions(ctx context.Context, lines []*Line) error {
	cmd := &etelecom.ImportExtensionsCommand{
		TenantID:   lines[0].TenantID,
		OwnerID:    lines[0].OwnerID,
		AccountID:  lines[0].AccountID,
		Extensions: []*etelecom.ImportExtensionInfo{},
	}
	for _, line := range lines {
		cmd.Extensions = append(cmd.Extensions, &etelecom.ImportExtensionInfo{
			ExtensionNumber: line.ExtensionNumber,
			ExpiresAt:       line.ExpiresAt,
			HotlineID:       line.HotlineID,
		})
	}

	return s.etelecomAggr.Dispatch(ctx, cmd)
}
