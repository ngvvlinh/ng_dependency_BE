package main

import (
	"encoding/base64"
	"io"
	"mime/multipart"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/apifw/httpx"
)

type reader struct {
	r   io.Reader
	buf []byte
	n   int
}

func newBufferReader(r io.Reader) (*reader, error) {
	// 4096 is default bufio size, do not change this value
	buf := make([]byte, 4096)
	for n := 0; n < 4096; {
		readN, err := r.Read(buf[n:])
		if err != nil {
			return nil, err
		}
		n += readN
	}
	return &reader{r: r, buf: buf, n: 0}, nil
}

func (r *reader) Read(p []byte) (n int, err error) {
	if r.n >= len(r.buf) {
		return r.r.Read(p)
	}
	copied := copy(p, r.buf[r.n:])
	if len(p) < len(r.buf)-r.n {
		r.n += len(p)
		return len(p), nil
	}
	r.n = len(r.buf)
	readN, err := r.r.Read(p[copied:])
	return copied + readN, err
}

func (r *reader) Peek(n int) ([]byte, error) {
	if n > len(r.buf) {
		n = len(r.buf)
	}
	return r.buf[0:n], nil
}

func (r *reader) Reset() {
	r.n = 0
}

type UploadedImages struct {
	Purpose Purpose
	Images  []UploadedImage
}

type UploadedImage struct {
	Source    io.Reader
	Filename  string
	Extension string
}

func getFormValue(form *multipart.Form, key string) string {
	values := form.Value[key]
	if values != nil {
		return values[0]
	}
	return ""
}

func readMultipartForm(c *httpx.Context) (result UploadedImages, err error) {
	form, err := c.MultipartForm()
	if err != nil {
		return
	}
	purpose := getFormValue(form, "purpose")
	if purpose == "" {
		purpose = getFormValue(form, "type") // backward-compatible
	}
	files := form.File["files"]
	if len(files) == 0 {
		err = cm.Error(cm.InvalidArgument, "No file", nil)
		return
	}
	images := make([]UploadedImage, len(files))
	for i, file := range files {
		multipartFile, err := file.Open()
		if err != nil {
			return result, cm.Error(cm.InvalidArgument, "Invalid ", err)
		}
		src, err := newBufferReader(multipartFile)
		if err != nil {
			return result, cm.Error(cm.InvalidArgument, "Invalid ", err)
		}
		format, err := verifyImage(file.Filename, int(file.Size), src)
		if err != nil {
			return result, cm.Error(cm.InvalidArgument, "Invalid ", err)
		}

		images[i] = UploadedImage{
			Source:    src,
			Filename:  file.Filename,
			Extension: format,
		}
	}
	result.Images = images
	return
}

func readBase64Images(c *httpx.Context) (result UploadedImages, err error) {
	filename := "uploaded file"
	src, err := newBufferReader(base64.NewDecoder(base64.StdEncoding, c.Req.Body))
	if err != nil {
		return result, cm.Error(cm.InvalidArgument, "Invalid ", err)
	}
	length := c.Req.ContentLength * 3 / 4
	format, err := verifyImage(filename, int(length), src)
	if err != nil {
		err = cm.Error(cm.InvalidArgument, "Invalid ", err)
		return
	}
	result.Purpose = Purpose(c.Req.Header.Get("Purpose"))
	result.Images = []UploadedImage{
		{
			Source:    src,
			Filename:  "",
			Extension: format,
		},
	}
	return
}
