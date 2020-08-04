package storage

import (
	"context"
	"io"

	cc "o.o/backend/pkg/common/config"
)

type Config struct {
	cc.GenericConfig
}

func Build(ctx context.Context, cfg Config) (result map[string]Bucket, err error) {
	err = cfg.Build(ctx, &result)
	return
}

type Drivers []DriverConfig

func (ds Drivers) Get(name string) Bucket {
	for _, d := range ds {
		if d.Name == name {
			return d.Driver
		}
	}
	return nil
}

type DriverConfig struct {
	Name   string
	Driver Bucket
}

type Bucket interface {
	Name() string

	OpenForRead(ctx context.Context, path string) (io.ReadCloser, error)

	OpenForWrite(ctx context.Context, path string) (io.WriteCloser, error)
}
