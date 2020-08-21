package imcsv

import (
	"bytes"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
	"o.o/backend/pkg/common/code/gencode"
	"o.o/backend/res/dl/imports"
	"o.o/common/l"
)

var (
	dlShopFulfillmentXlsx           []byte
	dlShopFulfillmentSimplifiedXlsx []byte
)

const (
	AssetShopFulfillmentPath               = "shop_fulfillments.v1.xlsx"
	assetShopFulfillmentFileName           = "shop_fulfillments.v1"
	AssetShopFulfillmentSimplifiedPath     = "shop_fulfillments.v1.simplified.xlsx"
	assetShopFulfillmentSimplifiedFileName = "shop_fulfillments.v1.simplified"
)

func init() {
	err := loadImportFile()
	if err != nil {
		ll.Fatal("Init import fulfillment", l.Error(err))
	}
}

func loadImportFile() error {
	data, err := imports.Asset(AssetShopFulfillmentPath)
	if err != nil {
		return err
	}

	_, err = excelize.OpenReader(bytes.NewReader(data))
	if err != nil {
		return err
	}

	dlShopFulfillmentXlsx = data

	_data, err := imports.Asset(AssetShopFulfillmentSimplifiedPath)
	if err != nil {
		return err
	}

	if _, err := excelize.OpenReader(bytes.NewReader(_data)); err != nil {
		return err
	}

	dlShopFulfillmentSimplifiedXlsx = _data

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

	for i := 2; i <= 7; i++ {
		code := codePrefix + strconv.Itoa(100 + i - 1)[1:]
		file.SetCellStr(sheetName, "A"+strconv.Itoa(i), code)
	}

	_ = file.Write(w)
	return nil
}

func GenerateImportFile(w http.ResponseWriter) (filename string, data []byte, err error) {
	randomCode := gencode.GenerateCode(gencode.Alphabet22, 5)
	codePrefix := "TEST-" + randomCode + "-"
	filename = assetShopFulfillmentFileName + "." + strings.ToLower(randomCode) + ".xlsx"

	var b bytes.Buffer
	err = randomizeFileWithCode(dlShopFulfillmentXlsx, codePrefix, &b)
	return filename, b.Bytes(), err
}

func GenerateImportSimplifiedFile(w http.ResponseWriter) (filename string, data []byte, err error) {
	randomCode := gencode.GenerateCode(gencode.Alphabet22, 5)
	codePrefix := "TEST-" + randomCode + "-"
	filename = assetShopFulfillmentSimplifiedFileName + "." + strings.ToLower(randomCode) + ".xlsx"

	var b bytes.Buffer
	err = randomizeFileWithCode(dlShopFulfillmentSimplifiedXlsx, codePrefix, &b)
	return filename, b.Bytes(), err
}
