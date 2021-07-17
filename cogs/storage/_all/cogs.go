package storage_all

import (
	"context"

	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/storage"
	"o.o/backend/pkg/common/storage/awsstorage"
	"o.o/backend/pkg/common/storage/filestorage"
	"o.o/backend/pkg/common/storage/gcloudstorage"
)

type DriverConfig struct {
	File   *filestorage.Config   `yaml:"file"`
	Gcloud *gcloudstorage.Config `yaml:"gcloud"`
	Aws    *awsstorage.Config    `yaml:"aws"`
}

func Build(ctx context.Context, c DriverConfig) (storage.Bucket, error) {
	return c.Build(ctx)
}

func (c DriverConfig) Build(ctx context.Context) (storage.Bucket, error) {
	switch {
	case c.File != nil:
		return filestorage.Connect(*c.File)
	case c.Gcloud != nil:
		return gcloudstorage.Connect(ctx, *c.Gcloud, []string{gcloudstorage.ReadWrite})
	case c.Aws != nil:
		return awsstorage.Connect(ctx, *c.Aws)
	default:
		return nil, cm.Errorf(cm.Internal, nil, "invalid config for storage driver")
	}
}

func DefaultDriver() DriverConfig {
	fileDriver := &filestorage.Config{RootPath: "/tmp/upload"}
	return DriverConfig{
		File: fileDriver,
	}
}
