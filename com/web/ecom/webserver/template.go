package webserver

import (
	"html/template"
	"io"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"o.o/api/webserver"
	"o.o/backend/pkg/common/cmenv"
	"o.o/capi/dot"
)

type Templates struct {
	templates   *template.Template
	filePattern string
	autoReload  bool
	lastReload  time.Time
}

var funcMap = template.FuncMap{
	"mul":            Mul,
	"add":            Add,
	"arr":            ArrNumber,
	"bigger":         BiggerNumber,
	"address":        AddressString,
	"comp":           GetFirstComparePrice,
	"checkContainID": checkContainID,
	"formatNumber":   formatPrice,
}

func parseTemplates(pattern string) (*Templates, error) {
	templates, err := template.New("template").Funcs(funcMap).ParseGlob(pattern)
	if err != nil {
		return nil, err
	}
	t := &Templates{
		templates:   templates,
		filePattern: pattern,
		autoReload:  cmenv.IsDev(),
	}
	return t, nil
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	if t.autoReload && time.Now().Sub(t.lastReload) > 5*time.Second {
		templates, err := template.New("template").Funcs(funcMap).ParseGlob(t.filePattern)
		if err != nil {
			return err
		}
		t.lastReload, t.templates = time.Now(), templates
		return templates.ExecuteTemplate(w, name, data)
	}
	return t.templates.ExecuteTemplate(w, name, data)
}

func (t *Templates) Reload() error {
	templates, err := template.ParseGlob(t.filePattern)
	if err != nil {
		return err
	}
	t.templates = templates
	return nil
}

func BiggerNumber(param1 int, param2 int) bool {
	if param1 > param2 {
		return true
	}
	return false
}

func ArrNumber(param int) []int {
	var result []int
	for i := 1; i <= param; i++ {
		result = append(result, i)
	}
	return result
}

func Mul(param1 int, param2 int) int {
	return param1 * param2
}

func Add(param1 int, param2 int) int {
	return param1 + param2
}

func AddressString(arg *webserver.AddressShopInfo) string {
	if arg == nil {
		return ""
	}
	result := ""
	if arg.Address != "" {
		result += arg.Address
	}
	if arg.Ward != "" {
		if result != "" {
			result += ", "
		}
		result += arg.Ward
	}
	if arg.District != "" {
		if result != "" {
			result += ", "
		}
		result += arg.District
	}
	if arg.Province != "" {
		if result != "" {
			result += ", "
		}
		result += arg.Province
	}
	if result != "" {
		result += "."
	}
	return result
}

func GetFirstComparePrice(arg *webserver.WsProduct) int {
	comparePrice := 0
	if arg == nil {
		return comparePrice
	}
	if len(arg.Product.Variants) == 0 {
		return comparePrice
	}
	comparePrice = arg.Product.Variants[0].RetailPrice
	for _, v := range arg.ComparePrice {
		if v.VariantID == arg.Product.Variants[0].VariantID {
			comparePrice = v.ComparePrice
			break
		}
	}
	return comparePrice
}

func checkContainID(ids []dot.ID, id dot.ID) bool {
	for _, v := range ids {
		if v == id {
			return true
		}
	}
	return false
}

func formatPrice(n int) string {
	p := message.NewPrinter(language.Vietnamese)
	return p.Sprint(n)
}
