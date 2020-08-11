package filestorage

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"

	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/cmenv"
	"o.o/backend/pkg/common/storage"
	"o.o/common/l"
)

var ll = l.New()

type Config struct {
	RootPath string `yaml:"root_dir"`
}

var _ storage.Bucket = &Driver{}

type Driver struct {
	rootPath    string
	rootAbsPath string
}

func Connect(cfg Config) (*Driver, error) {
	absPath, err := filepath.Abs(cfg.RootPath)
	if err != nil {
		return nil, cm.Errorf(cm.Internal, err, "%v", err)
	}
	stat, err := os.Stat(absPath)
	if cmenv.IsDev() && os.IsNotExist(err) {
		ll.Info("directory does not exist, auto create directory", l.String("path", absPath))
		err = os.MkdirAll(absPath, 0755)
		if err != nil {
			return nil, err
		}
		stat, err = os.Stat(absPath)
	}
	if err != nil {
		return nil, cm.Errorf(cm.Internal, err, "%v", err)
	}
	if !stat.IsDir() {
		return nil, cm.Errorf(cm.Internal, nil, "%v is not a directory", absPath)
	}
	if !isRegularDir(stat.Mode()) {
		return nil, cm.Errorf(cm.Internal, nil, "%v is not a regular directory")
	}
	return &Driver{
		rootPath:    cfg.RootPath,
		rootAbsPath: absPath,
	}, nil
}

func (b *Driver) Name() string {
	return b.rootPath
}

func (b *Driver) ensureDir(path string) (string, error) {
	if path == "" || path == "." { // current directory, ok
		return b.rootAbsPath, nil
	}
	if filepath.IsAbs(path) {
		return "", cm.Errorf(cm.Internal, nil, "path is absolute")
	}
	if strings.Contains(path, "..") {
		return "", cm.Errorf(cm.Internal, nil, "invalid path `%v`", path)
	}
	absPath := filepath.Join(b.rootAbsPath, path)

	// verify that the directory is inside absPath
	rel, err := filepath.Rel(b.rootAbsPath, absPath)
	if err != nil {
		return "", cm.Wrap(err)
	}
	if strings.Contains(rel, "..") {
		return "", cm.Errorf(cm.Internal, nil, "invalid path `%v`", path)
	}

	return absPath, os.MkdirAll(absPath, 0755)
}

func (b *Driver) OpenForWrite(ctx context.Context, path string) (io.WriteCloser, error) {
	dirPath := filepath.Dir(path)
	basePath := filepath.Base(path)
	absDirPath, err := b.ensureDir(dirPath)
	if err != nil {
		return nil, cm.Wrap(err)
	}
	file, err := os.Create(filepath.Join(absDirPath, basePath))
	if err != nil {
		return nil, cm.Wrap(err)
	}
	return file, nil
}

func (b *Driver) OpenForRead(ctx context.Context, path string) (io.ReadCloser, error) {
	panic("implement me")
}

// isRegularDir checks whether the directory is a regular directory (no special
// bit is set)
func isRegularDir(mode os.FileMode) bool {
	regularDir := mode & (os.ModeType & ^os.ModeDir)
	return regularDir == 0
}
