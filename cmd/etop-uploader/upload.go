package main

import (
	"bytes"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
	"path/filepath"

	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/apifw/httpx"
	"o.o/common/l"
	"o.o/common/xerrors"
)

const (
	minSize = 100
	maxSize = 1024 * 1024 // 1MB
	minWH   = 200
	maxWH   = 2000
)

func getImageConfig(purpose Purpose) (*ImageConfig, error) {
	if purpose == "" {
		purpose = PurposeDefault
	}
	config := imageConfigs[purpose]
	if config == nil {
		return nil, cm.Errorf(cm.InvalidArgument, nil, "Invalid image purpose")
	}
	return config, nil
}

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

	imgConfig, err := getImageConfig(uploadedImages.Purpose)
	if err != nil {
		return err
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
		if err := ensureDir(dirPath); err != nil {
			return err
		}
		if !func() bool {
			dst, err := os.Create(filepath.Join(dirPath, genName))
			if err != nil {
				ll.Error("error creating file", l.Error(err))
				errors[i] = NewUploadError(cm.Internal, cm.Internal.String(), uploadedImage.Filename).
					Log("can not upload", l.Error(err))
				return false
			}
			defer func() { _ = dst.Close() }()

			if _, err = io.Copy(dst, bytes.NewReader(uploadedImage.Source)); err != nil {
				ll.Error("error writing file", l.Error(err))
				errors[i] = NewUploadError(cm.Internal, cm.Internal.String(), uploadedImage.Filename).
					Log("can not upload", l.Error(err))
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

func ensureDir(dir string) error {
	return os.MkdirAll(dir, 0755)
}

func convertAPIErrorsToTwErrors(errs []*xerrors.APIError) []*xerrors.ErrorJSON {
	var result []*xerrors.ErrorJSON
	for _, err := range errs {
		result = append(result, xerrors.ToErrorJSON(xerrors.TwirpError(err)))
	}
	return result
}
