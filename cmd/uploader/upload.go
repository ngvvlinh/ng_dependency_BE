package main

import (
	"bytes"
	"context"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"path/filepath"

	"o.o/backend/cmd/uploader/config"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/httpx"
	"o.o/common/l"
	"o.o/common/xerrors"
)

type Purpose = config.Purpose

const (
	minSize = 100
	maxSize = 3 * 1024 * 1024 // 3MB
	minWH   = 200
	maxWH   = 2500
)

func NewUploadError(code xerrors.Code, msg, filename string) *xerrors.APIError {
	return cm.Error(code, msg, nil).
		WithMeta("filename", filename)
}

// UploaderHandler handle uploading with 2 different cases
//     1. Multipart form
//       - type: ""|"ahamove_verification"
//       - files: multiple files
//
//     2. Base64
//       - type: Header
//       - body: base64 encoded single file
func UploadHandler(c *httpx.Context) error {

	contentType := c.Req.Header.Get("Content-Type")
	reader := readMultipartForm
	if contentType == "application/base64" {
		reader = readBase64Images
	}
	uploadedImages, err := reader(c)
	if err != nil {
		return err
	}

	imgConfig, ok := cfg.Dirs.Get(uploadedImages.Purpose)
	if !ok {
		return cm.Errorf(cm.Internal, nil, "invalid purpose")
	}
	path := imgConfig.Path
	urlPrefix := imgConfig.URLPrefix

	countOK := 0
	errors := make([]*xerrors.APIError, len(uploadedImages.Images))
	result := make([]interface{}, len(uploadedImages.Images))
	for i, uploadedImage := range uploadedImages.Images {
		id := cm.NewBase54ID()
		genName := id + "." + uploadedImage.Extension
		subFolder := genName[:3]

		dirPath := filepath.Join(path, subFolder)
		filePath := filepath.Join(dirPath, genName)

		if !func() bool {

			// NOTE(vu): we use context.Background() here, instead of
			// c.Context(), to make the uploading task independent to the
			// incoming http request
			dst, err2 := bucket.OpenForWrite(context.Background(), filePath)
			if err2 != nil {
				ll.Error("error creating file", l.Error(err2))
				errors[i] = NewUploadError(cm.Internal, cm.Internal.String(), uploadedImage.Filename).
					Log("can not upload", l.Error(err2))
				return false
			}
			defer func() { _ = dst.Close() }()

			if _, err2 = io.Copy(dst, bytes.NewReader(uploadedImage.Source)); err2 != nil {
				ll.Error("error writing file", l.Error(err2))
				errors[i] = NewUploadError(cm.Internal, cm.Internal.String(), uploadedImage.Filename).
					Log("can not upload", l.Error(err2))
				return false
			}
			return true
		}() {
			continue
		}

		ll.Debug("Uploaded", l.String("filename", genName))
		resp := map[string]interface{}{
			"id":       id,
			"filename": uploadedImage.Filename,
		}
		if urlPrefix != "" {
			resp["url"] = fmt.Sprintf("%v/%v/%v", urlPrefix, subFolder, genName)
		}
		result[i] = resp
		errors[i] = NewUploadError(cm.NoError, "", uploadedImage.Filename)
		countOK++
	}

	if countOK == 0 {
		return errors[0]
	}
	c.SetResult(map[string]interface{}{
		"result": result,
		"errors": convertAPIErrorsToTwErrors(errors),
	})
	return nil
}

func verifyImage(filename string, size int, src io.Reader) (format string, data []byte, err error) {
	if size < minSize {
		return "", nil, NewUploadError(cm.InvalidArgument, "Invalid filesize", filename)
	}
	if size > maxSize {
		return "", nil, NewUploadError(cm.InvalidArgument, "File is too big (maximum 1MB)", filename)
	}
	data, err = readAll(src, size)
	if err != nil {
		return "", nil, NewUploadError(cm.InvalidArgument, fmt.Sprintf("Can not read file: %v", err), filename)
	}
	img, format, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		ll.Error("Unrecognized image file", l.String("filename", filename), l.Error(err))
		return "", nil, NewUploadError(cm.InvalidArgument, "Unrecognized image file", filename)
	}
	if img.Width < minWH || img.Width > maxWH || img.Height < minWH || img.Height > maxWH {
		return "", nil, NewUploadError(cm.InvalidArgument, "Image must be at least 200px * 200px and at most 2500px * 2500px", filename)
	}

	// Haravan does not accept .jpeg, so we have to change the extension
	if format == "jpeg" {
		format = "jpg"
	}
	return format, data, nil
}

func convertAPIErrorsToTwErrors(errs []*xerrors.APIError) []*xerrors.ErrorJSON {
	var result []*xerrors.ErrorJSON
	for _, err := range errs {
		result = append(result, xerrors.ToErrorJSON(xerrors.TwirpError(err)))
	}
	return result
}
