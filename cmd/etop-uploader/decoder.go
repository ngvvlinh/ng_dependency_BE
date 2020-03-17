package main

import (
	"bytes"
	"encoding/base64"
	"io"
	"mime/multipart"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/apifw/httpx"
)

type UploadedImages struct {
	Purpose Purpose
	Images  []UploadedImage
}

type UploadedImage struct {
	Source    []byte
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
		format, data, err := verifyImage(file.Filename, int(file.Size), multipartFile)
		if err != nil {
			return result, cm.Error(cm.InvalidArgument, "Invalid ", err)
		}

		images[i] = UploadedImage{
			Source:    data,
			Filename:  file.Filename,
			Extension: format,
		}
	}
	result.Images = images
	return result, err
}

func readBase64Images(c *httpx.Context) (result UploadedImages, err error) {
	filename := "uploaded file"
	src := base64.NewDecoder(base64.StdEncoding, c.Req.Body)
	length := int(c.Req.ContentLength) * 3 / 4
	format, data, err := verifyImage(filename, length, src)
	if err != nil {
		err = cm.Error(cm.InvalidArgument, "Invalid ", err)
		return
	}
	result.Purpose = Purpose(c.Req.Header.Get("Purpose"))
	result.Images = []UploadedImage{
		{
			Source:    data,
			Filename:  "",
			Extension: format,
		},
	}
	return
}

func readAll(r io.Reader, size int) ([]byte, error) {
	var b bytes.Buffer
	b.Grow(size)
	_, err := b.ReadFrom(r)
	return b.Bytes(), err
}
