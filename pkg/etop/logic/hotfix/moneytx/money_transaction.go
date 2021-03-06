package moneytx

import (
	"bytes"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"

	com "o.o/backend/com/main"
	identitymodel "o.o/backend/com/main/identity/model"
	txmodel "o.o/backend/com/main/moneytx/model"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/httpx"
	"o.o/backend/pkg/common/apifw/whitelabel/wl"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/code/gencode"
	"o.o/backend/pkg/common/sql/cmsql"
)

type HotFixMoneyTxService struct {
	db *cmsql.Database
}

func New(database com.MainDB) *HotFixMoneyTxService {
	return &HotFixMoneyTxService{db: database}
}

type Line struct {
	ShopCode  string
	ShopPhone string
	ShopName  string
	Amount    int
	Note      string
}

func (s *HotFixMoneyTxService) HandleImportMoneyTransactionManual(c *httpx.Context) error {
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
		return cm.Errorf(cm.InvalidArgument, err, "Không thể đọc được file. Vui lòng kiểm tra lại hoặc liên hệ %v.", wl.X(c.Context()).CSEmail).WithMeta("reason", "can not open file")
	}
	excelFile, err := excelize.OpenReader(bytes.NewReader(rawData))
	if err != nil {
		return cm.Errorf(cm.InvalidArgument, err, "Không thể đọc được file. Vui lòng kiểm tra lại hoặc liên hệ %v.", wl.X(c.Context()).CSEmail).WithMeta("reason", "invalid file format")
	}
	sheetName := excelFile.GetSheetName(1)
	rows := excelFile.GetRows(sheetName)
	if len(rows) <= 1 {
		return cm.Errorf(cm.InvalidArgument, nil, "File không có nội dung. Vui lòng tải lại file import hoặc liên hệ %v.", wl.X(c.Context()).CSEmail).WithMeta("reason", "no rows")
	}

	var lines []*Line
	for i, row := range rows {
		if i == 0 {
			continue
		}
		line, err := parseRow(row)
		if err != nil {
			return err
		}
		lines = append(lines, line)
	}

	_err := s.db.InTransaction(bus.Ctx(), func(s cmsql.QueryInterface) error {
		for _, line := range lines {
			if err := createMoneyTransactionShipping(s, line); err != nil {
				return err
			}
		}
		return nil
	})

	total := strconv.Itoa(len(lines))
	c.SetResult(map[string]string{"code": "ok", "total": total})
	return _err
}

func createMoneyTransactionShipping(s cmsql.QueryInterface, line *Line) error {
	var shop identitymodel.Shop
	if err := s.Table("shop").Where("code = ? AND phone = ?", line.ShopCode, line.ShopPhone).ShouldGet(&shop); err != nil {
		return cm.Errorf(cm.InvalidArgument, err, "Shop không tồn tại").WithMetap("line", line)
	}

	code := gencode.GenerateCodeWithType("M", shop.Code, time.Now())
	moneyTransaction := &txmodel.MoneyTransactionShipping{
		ID:          cm.NewID(),
		ShopID:      shop.ID,
		TotalAmount: line.Amount,
		Code:        code,
		Type:        "manual",
		Note:        line.Note,
	}
	if err := s.Table("money_transaction_shipping").ShouldInsert(moneyTransaction); err != nil {
		return cm.Errorf(cm.InvalidArgument, err, "Lỗi tạo phiên shop").WithMetap("line", line)
	}
	return nil
}

func parseRow(row []string) (*Line, error) {
	code := row[1]
	name := row[2]
	phone := row[3]
	amount, err := strconv.ParseFloat(row[4], 64)
	if err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, err, "Wrong amount format").WithMetap("row", row)
	}
	note := row[5]

	return &Line{
		ShopCode:  code,
		ShopPhone: phone,
		ShopName:  name,
		Amount:    int(amount),
		Note:      note,
	}, nil
}
