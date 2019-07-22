package cmService

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"etop.vn/backend/doc"
	"etop.vn/backend/pkg/common/idemp"
	"etop.vn/backend/pkg/etop/dl"
	"etop.vn/common/l"
)

const (
	MIMEExcel = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	MIMEOctet = "application/octet-stream"
)

var ll = l.New()
var idempgroup = idemp.NewGroup()

var mimeMap = map[string]string{
	".xlsx": MIMEExcel,
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

func SwaggerHandler(docFile string) http.Handler {
	data, err := doc.Asset(docFile)
	if err != nil {
		panic(err)
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(data)
	})
}

func ServeAssets(path string, contentType string) http.Handler {
	if contentType == "" {
		contentType = GetMIME(path)
	}

	data, err := dl.Asset(path)
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
			w.Write([]byte(`Lỗi không xác định. Vui lòng liên hệ hotro@etop.vn.`))
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

func RedocHandler() http.HandlerFunc {
	const tpl = `<!DOCTYPE html>
<html>
	<head>
	<title>eTop API Documentation</title>
	<meta charset="utf-8"/>
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<style>body { margin: 0; padding: 0 }</style>
	</head>
	<body>
	<redoc spec-url='%v'></redoc>
	<script src="https://rebilly.github.io/ReDoc/releases/latest/redoc.min.js"></script>
	</body>
</html>`

	return func(w http.ResponseWriter, r *http.Request) {
		path := filepath.Join(r.URL.Path, "swagger.json")
		_, _ = fmt.Fprintf(w, tpl, path)
	}
}

type Shutdowner interface {
	Register(fn func())
	Done() <-chan struct{}
}

type ShutdownImpl struct {
	funcs  []func()
	ctx    context.Context
	cancel func()
}

func NewShutdowner() *ShutdownImpl {
	newCtx, cancel := context.WithCancel(context.Background())
	return &ShutdownImpl{ctx: newCtx, cancel: cancel}
}

func (s *ShutdownImpl) Register(fn func()) {
	s.funcs = append(s.funcs, fn)
}

func (s *ShutdownImpl) ShutdownAll() {
	s.cancel()
	for _, fn := range s.funcs {
		fn()
	}
}

func (s *ShutdownImpl) Done() <-chan struct{} {
	return s.ctx.Done()
}
