package cmService

import (
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"o.o/backend/pkg/common/apifw/idemp"
	"o.o/backend/res/dl/fabo"
	"o.o/backend/res/dl/imports"
	"o.o/common/l"
)

const (
	MIMEExcel = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	MIMEOctet = "application/octet-stream"
	MINEPNG   = "image/png"
)

var ll = l.New()
var idempgroup = idemp.NewGroup()

var mimeMap = map[string]string{
	".xlsx": MIMEExcel,
	".png":  MINEPNG,
}

func GetMIME(filename string) string {
	ext := filepath.Ext(filename)
	typ := mimeMap[ext]
	if typ == "" {
		typ = MIMEOctet
	}
	return typ
}

func GetMIMEByExt(ext string) string {
	typ := mimeMap[ext]
	if typ == "" {
		typ = MIMEOctet
	}
	return typ
}

func ServeAssets(path string, contentType string) http.Handler {
	if contentType == "" {
		contentType = GetMIME(path)
	}

	data, err := imports.Asset(path)
	if err != nil {
		panic(err)
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", contentType)
		_, _ = w.Write(data)
	})
}

func ServeFaboAssets(path, contentType string) http.Handler {
	if contentType == "" {
		contentType = GetMIME(path)
	}

	data, err := fabo.Asset(path)
	if err != nil {
		panic(err)
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", contentType)
		_, _ = w.Write(data)
	})
}

type AssetsContent struct {
	FileName string
	Data     []byte
}

func ServeAssetsByContentGenerator(
	contentType, name string, timeout time.Duration,
	fn func(w http.ResponseWriter) (filename string, data []byte, err error),
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assetContent, err, _ := idempgroup.Do(name, timeout, func() (interface{}, error) {
			filename, data, err := fn(w)
			return &AssetsContent{
				FileName: filename,
				Data:     data,
			}, err
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`Lỗi không xác định.`))
			return
		}

		asset := assetContent.(*AssetsContent)

		h := w.Header()
		h.Add("Content-Type", contentType)
		h.Add("Content-Disposition", fmt.Sprintf(`attachment; filename=%q`, asset.FileName))
		h.Add("Cache-Control", "no-cache, no-store, must-revalidate")
		w.WriteHeader(http.StatusOK)
		w.Write(asset.Data)
	})
}
