package servedoc

import (
	"fmt"
	"net/http"
	"path/filepath"

	"o.o/backend/doc"
	"o.o/backend/pkg/common/cmenv"
)

func SwaggerHandler(docFile string) http.Handler {
	data, err := doc.Asset(docFile)
	if err != nil {
		panic(err)
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if cmenv.IsDev() { // always reload file on dev
			data, err = doc.Asset(docFile)
			if err != nil {
				panic(err)
			}
		}
		_, _ = w.Write(data)
	})
}

func RedocHandler() http.HandlerFunc {
	const tpl = `<!DOCTYPE html>
<html>
	<head>
	<title>API Documentation</title>
	<meta charset="utf-8"/>
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<style>body { margin: 0; padding: 0 }</style>
	</head>
	<body>
	<redoc spec-url='%v'></redoc>
	<script src="https://cdn.jsdelivr.net/npm/redoc@next/bundles/redoc.standalone.js"></script>
	</body>
</html>`

	return func(w http.ResponseWriter, r *http.Request) {
		path := filepath.Join(r.URL.Path, "swagger.json")
		_, _ = fmt.Fprintf(w, tpl, path)
	}
}
