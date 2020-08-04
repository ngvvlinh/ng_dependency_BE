package upload

import (
	"bytes"
	"context"
	"io"
	"path/filepath"

	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/storage"
	"o.o/common/l"
)

var ll = l.New()

type StoreFileCommand struct {
	UploadType string
	FileName   string
	Data       []byte `json:"-"`
}

type Uploader struct {
	driver storage.Bucket
	dirs   map[string]string
}

func NewUploader(driver storage.Bucket, dirs map[string]string) (*Uploader, error) {
	return &Uploader{driver: driver, dirs: dirs}, nil
}

func (u *Uploader) ExpectDir(uploadType string) {
	if u.dirs[uploadType] == "" {
		ll.Panic("Expected directory type", l.String("uploadType", uploadType))
	}
}

func (u *Uploader) StoreFile(cmd *StoreFileCommand) (_err error) {
	if cmd.UploadType == "" || cmd.FileName == "" || len(cmd.Data) == 0 {
		return cm.Error(cm.Internal, "Invalid arg", nil)
	}

	dir := u.dirs[cmd.UploadType]
	if dir == "" {
		return cm.Errorf(cm.Internal, nil, "Can not upload file").
			WithMetaf("reason", "no directory of type %v", cmd.UploadType)
	}

	filePath := filepath.Join(dir, cmd.FileName)
	dst, err := u.driver.OpenForWrite(context.Background(), filePath)
	if err != nil {
		return cm.Errorf(cm.Internal, err, "Can not upload file")
	}
	defer func() {
		err2 := dst.Close()
		if _err == nil {
			_err = err2
		}
	}()

	_, err = io.Copy(dst, bytes.NewReader(cmd.Data))
	if err != nil {
		return err
	}
	ll.Info("Stored file", l.String("path", filePath))
	return nil
}
