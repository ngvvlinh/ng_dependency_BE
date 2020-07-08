package money_transaction

import (
	"bytes"
	"context"
	"io/ioutil"
	"mime/multipart"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/valyala/tsvreader"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/whitelabel/wl"
)

func ReadXLSXFile(ctx context.Context, file multipart.File) (rows [][]string, _ error) {
	rawData, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, err, "Không thể đọc được file. Vui lòng kiểm tra lại hoặc liên hệ %v.", wl.X(ctx).CSEmail).WithMeta("reason", "can not open file")
	}

	excelFile, err := excelize.OpenReader(bytes.NewReader(rawData))
	if err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, err, "Không thể đọc được file. Vui lòng kiểm tra lại hoặc liên hệ %v.", wl.X(ctx).CSEmail).WithMeta("reason", "invalid file format")
	}

	sheetName := excelFile.GetSheetName(1)
	rows = excelFile.GetRows(sheetName)
	if len(rows) <= 1 {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "File không có nội dung. Vui lòng tải lại file import hoặc liên hệ %v.", wl.X(ctx).CSEmail).WithMeta("reason", "no rows")
	}
	return
}

func ReadCSVFile(ctx context.Context, file multipart.File, nColumn int) (rows [][]string, _ error) {
	r := tsvreader.New(file)
	for r.Next() {
		var row []string
		for i := 0; i < nColumn; i++ {
			row = append(row, r.String())
		}
		r.SkipCol()
		rows = append(rows, row)
	}
	if err := r.Error(); err != nil {
		return nil, cm.Errorf(cm.InvalidArgument, err, "Không thể đọc được file. Vui lòng kiểm tra lại hoặc liên hệ %v.", wl.X(ctx).CSEmail).WithMeta("reason", "invalid file format")
	}
	return
}
