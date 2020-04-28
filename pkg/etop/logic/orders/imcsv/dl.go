package imcsv

import (
	"bytes"
	"io"
	"net/http"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"

	"o.o/backend/pkg/common/code/gencode"
	"o.o/backend/res/dl/imports"
	"o.o/common/l"
)

var dlShopOrderXlsx []byte

const (
	AssetShopOrderPath     = "imports/shop_orders.v1.xlsx"
	assetShopOrderFilename = "shop_orders.v1"
)

func init() {
	err := loadImportFile()
	if err != nil {
		ll.Fatal("Init import order", l.Error(err))
	}
}

func loadImportFile() error {
	data, err := imports.Asset(AssetShopOrderPath)
	if err != nil {
		return err
	}

	_, err = excelize.OpenReader(bytes.NewReader(data))
	if err != nil {
		return err
	}

	dlShopOrderXlsx = data
	return nil
}

func randomizeFileWithCode(data []byte, codePrefix string, w io.Writer) error {
	file, err := excelize.OpenReader(bytes.NewReader(data))
	if err != nil {
		return err
	}

	sheetName, err := validateSheets(file)
	if err != nil {
		return err
	}

	file.SetCellStr(sheetName, "A2", codePrefix+"1")
	file.SetCellStr(sheetName, "A3", codePrefix+"1")
	file.SetCellStr(sheetName, "A4", codePrefix+"1")

	file.SetCellStr(sheetName, "A5", codePrefix+"2")

	file.SetCellStr(sheetName, "A6", codePrefix+"3")
	file.SetCellStr(sheetName, "A7", codePrefix+"3")
	// skip A8
	file.SetCellStr(sheetName, "A9", codePrefix+"4")

	_ = file.Write(w)
	return nil
}

func GenerateImportFile(w http.ResponseWriter) (filename string, data []byte, err error) {
	randomCode := gencode.GenerateCode(gencode.Alphabet22, 5)
	codePrefix := "TEST-" + randomCode + "-"
	filename = assetShopOrderFilename + "." + strings.ToLower(randomCode) + ".xlsx"

	var b bytes.Buffer
	err = randomizeFileWithCode(dlShopOrderXlsx, codePrefix, &b)
	return filename, b.Bytes(), err
}
