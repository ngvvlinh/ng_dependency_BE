package main

import (
	"image"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/httpx"
	"etop.vn/backend/pkg/common/l"
)

const (
	minSize = 100
	maxSize = 1024 * 1024 // 1MB
	minWH   = 200
	maxWH   = 2000
)

func NewUploadError(code cm.Code, msg, filename string) error {
	return cm.Error(code, msg, nil).
		WithMeta("filename", filename)
}

func UploadHandler(c *httpx.Context) error {

	// Multipart form
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}

	files := form.File["files"]
	if len(files) == 0 {
		return cm.Error(cm.InvalidArgument, "No file", nil)
	}

	exts := make([]string, len(files))
	for i, file := range files {
		format, err := verifyImage(file)
		if err != nil {
			return cm.Error(cm.InvalidArgument, "Invalid ", err)
		}

		// Haravan does not accept .jpeg, so we have to change the extension
		if format == "jpeg" {
			format = "jpg"
		}
		exts[i] = format
	}

	countOK := 0
	errors := make([]error, len(files))
	result := make([]interface{}, len(files))
	for i, file := range files {
		id := cm.NewBase54ID()
		genName := id + "." + exts[i]
		src, err := file.Open()
		if err != nil {
			ll.Info("Unexpected", l.Error(err))
			return cm.Error(cm.InvalidArgument, "", err)
		}

		dst, err := os.Create(filepath.Join(cfg.UploadDirImg, genName))
		if err != nil {
			errors[i] = NewUploadError(cm.Internal, cm.Internal.String(), file.Filename)
			continue
		}
		defer dst.Close()

		if _, err = io.Copy(dst, src); err != nil {
			ll.Info("Error writing file", l.Error(err))
			errors[i] = NewUploadError(cm.Internal, cm.Internal.String(), file.Filename)
			continue
		}

		ll.Info("Uploaded", l.String("filename", genName))
		resp := map[string]interface{}{
			"id":       id,
			"filename": file.Filename,
		}
		if cfg.URLPrefix != "" {
			resp["url"] = cfg.URLPrefix + "/" + genName
		}
		result[i] = resp
		errors[i] = NewUploadError(cm.NoError, "", file.Filename)
		countOK++
	}

	if countOK == 0 {
		return errors[0]
	}
	c.SetResult(map[string]interface{}{
		"result": result,
		"errors": errors,
	})
	return nil
}

func verifyImage(file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	if file.Size < minSize {
		return "", NewUploadError(cm.InvalidArgument, "Invalid filesize", file.Filename)
	}
	if file.Size > maxSize {
		return "", NewUploadError(cm.InvalidArgument, "File is too big (maximum 1MB)", file.Filename)
	}

	img, format, err := image.DecodeConfig(src)
	if err != nil {
		ll.Error("Unrecognized image file", l.String("filename", file.Filename), l.Error(err))
		return "", NewUploadError(cm.InvalidArgument, "Unrecognized image file", file.Filename)
	}

	if img.Width < minWH || img.Width > maxWH || img.Height < minWH || img.Height > maxWH {
		return "", NewUploadError(cm.InvalidArgument, "Image must be at least 200px * 200px and at most 2500px * 2500px", file.Filename)
	}

	return format, nil
}
