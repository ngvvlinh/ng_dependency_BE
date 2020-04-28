package imcsv

import (
	"context"

	"o.o/backend/com/main/catalog/convert"
	catalogmodel "o.o/backend/com/main/catalog/model"
	"o.o/backend/pkg/common/imcsv"
	"o.o/backend/pkg/common/validate"
	"o.o/capi/dot"
)

var schemaV0 = imcsv.Schema{
	{
		Name:    "_",
		Display: "Loại hàng",
		Norm:    "loai hang",
	},
	{
		Name:    "category",
		Display: "Nhóm hàng",
		Norm:    "nhom hang 3 cap",
	},
	{
		Name:    "collections",
		Display: "Bộ sưu tập",
		Norm:    "bo suu tap ngang hang",
	},
	{
		Name:    "product_code",
		Display: "Mã hàng",
		Norm:    "ma hang",
	},
	{
		Name:    "variant_code",
		Display: "Mã thuộc tính",
		Norm:    "ma thuoc tinh",
	},
	{
		Name:    "product_name",
		Display: "Tên hàng hóa",
		Norm:    "ten hang hoa",
	},
	{
		Name:    "attributes",
		Display: "Thuộc tính",
		Norm:    "thuoc tinh",
	},
	{
		Name:    "list_price",
		Display: "Giá bán",
		Norm:    "gia ban",
	},
	{
		Name:    "cost_price",
		Display: "Giá vốn",
		Norm:    "gia von",
	},
	{
		Name:    "quantity_available",
		Display: "Tồn kho",
		Norm:    "ton kho",
	},
	{
		Name:    "unit",
		Display: "Đơn vị tính",
		Norm:    "dvt",
	},
	{
		Name:    "images",
		Display: "Hình ảnh",
		Norm:    "hinh anh url1 url2",
	},
	{
		Name:    "weight",
		Display: "Trọng lượng",
		Norm:    "trong luong",
	},
	{
		Name:    "description",
		Display: "Mô tả",
		Norm:    "mo ta",
	},
}

var schemaV1 = imcsv.Schema{
	{
		Name:    "category",
		Display: "Danh mục",
		Norm:    "danh muc 3 cap",
	},
	{
		Name:    "collections",
		Display: "Bộ sưu tập",
		Norm:    "bo suu tap ngang hang",
	},
	{
		Name:    "product_code",
		Display: "Mã sản phẩm",
		Norm:    "ma san pham",
	},
	{
		Name:    "variant_code",
		Display: "Mã phiên bản sản phẩm",
		Norm:    "ma phien ban san pham",
	},
	{
		Name:    "product_name",
		Display: "Tên sản phẩm",
		Norm:    "ten san pham",
	},
	{
		Name:    "attributes",
		Display: "Thuộc tính",
		Norm:    "thuoc tinh",
	},
	{
		Name:    "list_price",
		Display: "Giá bán",
		Norm:    "gia ban",
	},
	{
		Name:    "cost_price",
		Display: "Giá vốn",
		Norm:    "gia von",
	},
	{
		Name:    "quantity_available",
		Display: "Tồn kho",
		Norm:    "ton kho",
	},
	{
		Name:    "unit",
		Display: "Đơn vị tính",
		Norm:    "don vi tinh",
	},
	{
		Name:    "images",
		Display: "Hình ảnh",
		Norm:    "hinh anh link1 link2",
	},
	{
		Name:    "weight",
		Display: "Khối lượng (kg)",
		Norm:    "khoi luong kg",
	},
	{
		Name:    "description",
		Display: "Mô tả",
		Norm:    "mo ta",
	},
}

type indexes struct {
	indexer imcsv.Indexer

	category      int
	collections   int
	productCode   int
	variantCode   int
	productName   int
	attributes    int
	listPrice     int
	costPrice     int
	quantityAvail int
	unit          int
	images        int
	weight        int
	description   int
}

var schemas = []imcsv.Schema{schemaV0, schemaV1}
var idxes = []indexes{
	initIndexes(schemaV0),
	initIndexes(schemaV1),
}

func initIndexes(schema imcsv.Schema) indexes {
	indexer := schema.Indexer()
	return indexes{
		category:      indexer("category"),
		collections:   indexer("collections"),
		productCode:   indexer("product_code"),
		variantCode:   indexer("variant_code"),
		productName:   indexer("product_name"),
		attributes:    indexer("attributes"),
		listPrice:     indexer("list_price"),
		costPrice:     indexer("cost_price"),
		quantityAvail: indexer("quantity_available"),
		unit:          indexer("unit"),
		images:        indexer("images"),
		weight:        indexer("weight"),
		description:   indexer("description"),
	}
}

func validateSchema(ctx context.Context, headerRow *[]string) (schema imcsv.Schema, idx indexes, errs []error, err error) {
	i, indexer, errs, err := imcsv.ValidateAgainstSchemas(ctx, headerRow, schemas)
	if err != nil || errs != nil {
		return
	}

	idx = idxes[i] // clone the struct
	idx.indexer = indexer
	return schemas[i], idx, nil, nil
}

type RowProduct struct {
	RowIndex      int
	Category      [3]string
	Collections   []string
	ProductCode   string
	VariantCode   string
	ProductName   string
	Attributes    []*catalogmodel.ProductAttribute
	ListPrice     int
	CostPrice     int
	QuantityAvail int
	Unit          string
	ImageURLs     []string
	Weight        int // gram
	Description   string

	categoryID    dot.ID
	collectionIDs []dot.ID
	nameNormUa    string
	attrNormKv    string
}

func (m *RowProduct) Validate(schema imcsv.Schema, idx indexes, mode Mode) (errs []error) {
	var col int
	r := m.RowIndex

	col = idx.productName
	if m.ProductName == "" {
		errs = append(errs, imcsv.CellError(idx.indexer, r, col, "%v không được để trống.", schema[col].Display))
	}

	col = idx.listPrice
	if m.ListPrice == 0 {
		errs = append(errs, imcsv.CellError(idx.indexer, r, col, "%v không được để trống.", schema[col].Display))
	}

	return errs
}

// GetProductKey returns ProductCode or fallback to normalize(ProductName). It
// returns a key for grouping variants into a product. All rows must provide all
// codes or leave them all empty so this function can rely.
func (m *RowProduct) GetProductKey() string {
	if m.ProductCode != "" {
		return m.ProductCode
	}
	if m.nameNormUa != "" {
		return m.nameNormUa
	}
	m.nameNormUa = validate.NormalizeUnaccent(m.ProductName)
	return m.nameNormUa
}

func (m *RowProduct) GetProductCodeOrName() string {
	if m.ProductCode != "" {
		return m.ProductCode
	}
	return m.ProductName
}

func (m *RowProduct) GetProductNameOrCode() string {
	if m.ProductName != "" {
		return m.ProductName
	}
	return m.ProductCode
}

func (m *RowProduct) GetVariantAttrNorm() string {
	if m.attrNormKv != "" {
		return m.attrNormKv
	}
	_, m.attrNormKv = catalogmodel.NormalizeAttributes(convert.Convert_catalogmodel_ProductAttributes_catalogtypes_Attributes(m.Attributes))
	if m.attrNormKv != "_" {
		m.attrNormKv = validate.NormalizedSearchToTsVector(m.attrNormKv)
	}
	return m.attrNormKv
}
