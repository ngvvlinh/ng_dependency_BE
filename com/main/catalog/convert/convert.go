package convert

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/microcosm-cc/bluemonday"

	"etop.vn/api/main/catalog"
	catalogmodel "etop.vn/backend/com/main/catalog/model"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/validate"
	"etop.vn/common/l"
)

// +gen:convert: etop.vn/backend/com/main/catalog/model->etop.vn/api/main/catalog,etop.vn/api/main/catalog/types
// +gen:convert: etop.vn/api/main/catalog

var ll = l.New()
var htmlPolicy = bluemonday.UGCPolicy()

const (
	MaxCodeNorm        = 999999
	MaxCodeNormVariant = 100
	codeRegex          = "^SP([0-9]{6})$"
	codePrefix         = "SP"
)

var reCode = regexp.MustCompile(codeRegex)

func ParseCodeNorm(code string) (_ int, ok bool) {
	parts := reCode.FindStringSubmatch(code)
	if len(parts) == 0 {
		return 0, false
	}
	number, err := strconv.Atoi(parts[1])
	if err != nil {
		ll.Panic("unexpected", l.Error(err))
	}
	return number, true
}

func GenerateCodeProduct(codeNorm int) string {
	return fmt.Sprintf("%v%06v", codePrefix, codeNorm)
}

func GenerateCodeVariant(productCode string, codeNorm int) string {
	return fmt.Sprintf("%v-%02v", productCode, codeNorm)
}

func shopProduct(in *catalogmodel.ShopProduct, out *catalog.ShopProduct) {
	metaFields := []*catalog.MetaField{}
	for _, metaField := range in.MetaFields {
		metaFields = append(metaFields, &catalog.MetaField{
			Key:   metaField.Key,
			Value: metaField.Value,
		})
	}

	convert_catalogmodel_ShopProduct_catalog_ShopProduct(in, out)
	out.MetaFields = metaFields
}

func shopProductDB(in *catalog.ShopProduct, out *catalogmodel.ShopProduct) {

	metaFields := []*catalogmodel.MetaField{}
	for _, metaField := range in.MetaFields {
		metaFields = append(metaFields, &catalogmodel.MetaField{
			Key:   metaField.Key,
			Value: metaField.Value,
		})
	}
	convert_catalog_ShopProduct_catalogmodel_ShopProduct(in, out)
	out.MetaFields = metaFields
}

func ShopProductWithVariants(in *catalogmodel.ShopProductWithVariants) (out *catalog.ShopProductWithVariants) {
	if in == nil {
		return nil
	}
	shopVariants := Convert_catalogmodel_ShopVariants_catalog_ShopVariants(in.Variants)
	out = &catalog.ShopProductWithVariants{
		ShopProduct: Convert_catalogmodel_ShopProduct_catalog_ShopProduct(in.ShopProduct, nil),
		Variants:    shopVariants,
	}
	return out
}

func ShopProductsWithVariants(ins []*catalogmodel.ShopProductWithVariants) (outs []*catalog.ShopProductWithVariants) {
	outs = make([]*catalog.ShopProductWithVariants, len(ins))
	for i, in := range ins {
		outs[i] = ShopProductWithVariants(in)
	}
	return outs
}

func shopVariantDB(in *catalog.ShopVariant, out *catalogmodel.ShopVariant) {
	convert_catalog_ShopVariant_catalogmodel_ShopVariant(in, out)
	attributes, attrNormKv := catalogmodel.NormalizeAttributes(in.Attributes)
	out.Attributes = Convert_catalogtypes_Attributes_catalogmodel_ProductAttributes(attributes)
	out.AttrNormKv = attrNormKv
}

func ShopVariantWithProduct(in *catalogmodel.ShopVariantWithProduct) (out *catalog.ShopVariantWithProduct) {
	if in == nil {
		return nil
	}
	out = &catalog.ShopVariantWithProduct{
		ShopProduct: Convert_catalogmodel_ShopProduct_catalog_ShopProduct(in.ShopProduct, nil),
		ShopVariant: Convert_catalogmodel_ShopVariant_catalog_ShopVariant(in.ShopVariant, nil),
	}
	return out
}

func ShopVariantsWithProduct(ins []*catalogmodel.ShopVariantWithProduct) (outs []*catalog.ShopVariantWithProduct) {
	outs = make([]*catalog.ShopVariantWithProduct, len(ins))
	for i, in := range ins {
		outs[i] = ShopVariantWithProduct(in)
	}
	return outs
}

func createShopBrand(args *catalog.CreateBrandArgs, out *catalog.ShopBrand) {
	apply_catalog_CreateBrandArgs_catalog_ShopBrand(args, out)
	out.ID = cm.NewID()
}

func createShopProduct(arg *catalog.CreateShopProductArgs, out *catalog.ShopProduct) {
	apply_catalog_CreateShopProductArgs_catalog_ShopProduct(arg, out)
	out.ProductID = cm.NewID()
	out.Code = NormalizeExternalCode(arg.Code)

	out.ShortDesc = arg.ShortDesc
	out.Description = arg.Description
	out.DescHTML = htmlPolicy.Sanitize(arg.DescHTML)
	out.CostPrice = arg.CostPrice
	out.ListPrice = arg.ListPrice
	out.RetailPrice = arg.RetailPrice
}

func updateShopProduct(args *catalog.UpdateShopProductInfoArgs, in *catalog.ShopProduct) *catalog.ShopProduct {
	if in == nil {
		return nil
	}
	apply_catalog_UpdateShopProductInfoArgs_catalog_ShopProduct(args, in)
	in.UpdatedAt = time.Now()
	if args.DescHTML.Valid == true {
		var descHTML = htmlPolicy.Sanitize(args.DescHTML.String)
		in.DescHTML = descHTML
	}
	if args.Code.Valid {
		in.Code = NormalizeExternalCode(args.Code.String)
	}
	return in
}

func createShopVariant(arg *catalog.CreateShopVariantArgs, out *catalog.ShopVariant) {
	apply_catalog_CreateShopVariantArgs_catalog_ShopVariant(arg, out)
	out.VariantID = cm.NewID()
	out.Status = 0
	out.Code = NormalizeExternalCode(arg.Code)
	out.DescHTML = htmlPolicy.Sanitize(arg.DescHTML)
	out.ShortDesc = arg.ShortDesc
	out.Description = arg.Description
	out.CostPrice = arg.CostPrice
	out.ListPrice = arg.ListPrice
	out.RetailPrice = arg.RetailPrice
}

func updateShopVariant(args *catalog.UpdateShopVariantInfoArgs, in *catalog.ShopVariant) *catalog.ShopVariant {
	if in == nil {
		return nil
	}
	apply_catalog_UpdateShopVariantInfoArgs_catalog_ShopVariant(args, in)
	in.UpdatedAt = time.Now()
	if args.DescHTML.Valid == true {
		var descHTML = htmlPolicy.Sanitize(args.DescHTML.String)
		in.DescHTML = descHTML
	}
	if args.Code.Valid {
		in.Code = NormalizeExternalCode(args.Code.String)
	}
	return in
}

func updateShopCollection(args *catalog.UpdateShopCollectionArgs, in *catalog.ShopCollection) (out *catalog.ShopCollection) {
	if in == nil {
		return nil
	}
	apply_catalog_UpdateShopCollectionArgs_catalog_ShopCollection(args, in)
	in.UpdatedAt = time.Now()
	return in
}
func updateShopCategory(args *catalog.UpdateShopCategoryArgs, in *catalog.ShopCategory) (out *catalog.ShopCategory) {
	if in == nil {
		return nil
	}
	apply_catalog_UpdateShopCategoryArgs_catalog_ShopCategory(args, in)
	in.UpdatedAt = time.Now()
	return in
}
func updateShopProductCategory(args *catalog.UpdateShopProductCategoryArgs, in *catalog.ShopProduct) (out *catalog.ShopProduct) {
	if in == nil {
		return nil
	}
	apply_catalog_UpdateShopProductCategoryArgs_catalog_ShopProduct(args, in)
	in.UpdatedAt = time.Now()
	return in
}

func NormalizeExternalCode(s string) string {
	s = strings.ReplaceAll(s, " ", "-")
	if len(s) > 0 && (string(s[len(s)-1]) == "-" || string(s[0]) == "-") {
		return normalizeExternalCode(s)
	}
	for i := 0; i < len(s); i++ {
		if !validate.ExternalCodeCharacter(s[i]) || ((i < len(s)-1) && string(s[i]) == "-" && string(s[i+1]) == "-") {
			return normalizeExternalCode(s)
		}
	}
	return s
}

func normalizeExternalCode(s string) string {
	res := make([]byte, 0, len(s))
	spaceMark := true
	for i := 0; i < len(s); i++ {
		if validate.ExternalCodeCharacter(s[i]) {

			if string(s[i]) == "-" && !spaceMark && (i < len(s)-1 && (s[i+1] != "-"[0])) {
				res = append(res, "-"[0])
			}
			if string(s[i]) != "-" {
				res = append(res, s[i])
				spaceMark = false
			}
		}
	}
	return string(res)
}
