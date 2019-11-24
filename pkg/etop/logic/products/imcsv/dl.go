package imcsv

import (
	"bytes"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"

	"etop.vn/backend/pkg/common/gencode"
	"etop.vn/backend/res/dl/imports"
	"etop.vn/common/l"
)

var dlShopProductXlsx []byte

const (
	AssetShopProductPath     = "imports/shop_products.v1.xlsx"
	assetShopProductFilename = "shop_products.v1"

	AssetShopProductSimplifiedPath = "imports/shop_products.v1.simplified.xlsx"
)

func init() {
	err := loadImportFile()
	if err != nil {
		ll.Fatal("Init import product", l.Error(err))
	}
}

func loadImportFile() error {
	data, err := imports.Asset(AssetShopProductPath)
	if err != nil {
		return err
	}

	_, err = excelize.OpenReader(bytes.NewReader(data))
	if err != nil {
		return err
	}

	dlShopProductXlsx = data
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

	codePrefixP := codePrefix + "SP-"
	file.SetCellStr(sheetName, "C2", codePrefixP+"01")
	file.SetCellStr(sheetName, "C3", codePrefixP+"01")
	file.SetCellStr(sheetName, "C4", codePrefixP+"01")
	file.SetCellStr(sheetName, "C5", codePrefixP+"01")
	file.SetCellStr(sheetName, "C6", codePrefixP+"01")
	file.SetCellStr(sheetName, "C7", codePrefixP+"01")
	file.SetCellStr(sheetName, "C8", codePrefixP+"07")
	file.SetCellStr(sheetName, "C9", codePrefixP+"07")
	file.SetCellStr(sheetName, "C10", codePrefixP+"09")
	file.SetCellStr(sheetName, "C11", codePrefixP+"09")
	file.SetCellStr(sheetName, "C12", codePrefixP+"11")
	file.SetCellStr(sheetName, "C13", codePrefixP+"12")
	file.SetCellStr(sheetName, "C14", codePrefixP+"13")
	file.SetCellStr(sheetName, "C15", codePrefixP+"14")

	for i := 2; i <= 15; i++ {
		code := codePrefix + strconv.Itoa(100 + i - 1)[1:]
		file.SetCellStr(sheetName, "D"+strconv.Itoa(i), code)
	}

	_ = file.Write(w)
	return nil
}

func GenerateImportFile(w http.ResponseWriter) (filename string, data []byte, err error) {
	randomCode := gencode.GenerateCode(gencode.Alphabet22, 5)
	codePrefix := "TEST-" + randomCode + "-"
	filename = assetShopProductFilename + "." + strings.ToLower(randomCode) + ".xlsx"

	var b bytes.Buffer
	err = randomizeFileWithCode(dlShopProductXlsx, codePrefix, &b)
	return filename, b.Bytes(), err
}
