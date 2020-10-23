package reportserver

import (
	"fmt"
	"html/template"

	"o.o/backend/pkg/common/projectpath"
)

func parseTemplate(pattern string) *template.Template {
	dirPath := projectpath.GetPath() + "/com/report/templates"
	return template.Must(template.New("layout.html").Funcs(funcMap).
		ParseFiles(dirPath+"/layout.html",
			fmt.Sprintf("%s/%s", dirPath, pattern)))
}
