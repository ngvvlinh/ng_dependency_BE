package upload

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	cm "o.o/backend/pkg/common"
	"o.o/common/l"
)

var ll = l.New()

type Config struct {
	DirImportShopOrder   string `yaml:"dir_import_shop_order"`
	DirImportShopProduct string `yaml:"dir_import_shop_product"`
}

type StoreFileCommand struct {
	UploadType string
	FileName   string
	Data       []byte `json:"-"`
}

type Uploader struct {
	dirs map[string]string
}

func NewUploader(dirs map[string]string) (*Uploader, error) {
	uploadDirs := make(map[string]string)
	for typ, dirPath := range dirs {
		fi, err := os.Stat(dirPath)
		if err != nil {
			return nil, fmt.Errorf("upload dir %v: %v", typ, err)
		}
		if !fi.IsDir() {
			return nil, fmt.Errorf("upload dir %v: `%v` is not a directory", typ, err)
		}

		absPath, err := filepath.Abs(dirPath)
		if err != nil {
			panic(err)
		}
		uploadDirs[typ] = absPath
	}
	return &Uploader{dirs: uploadDirs}, nil
}

func (u *Uploader) ExpectDir(uploadType string) {
	if u.dirs[uploadType] == "" {
		ll.Panic("Expected directory type", l.String("uploadType", uploadType))
	}
}

func (u *Uploader) StoreFile(cmd *StoreFileCommand) error {
	if cmd.UploadType == "" || cmd.FileName == "" || len(cmd.Data) == 0 {
		return cm.Error(cm.Internal, "Invalid arg", nil)
	}

	dir := u.dirs[cmd.UploadType]
	if dir == "" {
		return cm.Error(cm.Internal, "Can not upload file", nil).
			WithMetaf("reason", "no directory of type %v", cmd.UploadType)
	}

	filePath := filepath.Join(dir, cmd.FileName)
	err := ioutil.WriteFile(filePath, cmd.Data, os.ModePerm)
	if err != nil {
		return cm.Error(cm.Internal, "", err)
	}
	ll.Info("Stored file", l.String("path", filePath))
	return nil
}
