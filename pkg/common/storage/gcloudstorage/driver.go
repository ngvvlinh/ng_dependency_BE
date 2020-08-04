package gcloudstorage

import (
	"context"
	"crypto/sha256"
	"io"
	"sync"

	gstorage "cloud.google.com/go/storage"
	"google.golang.org/api/option"

	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/gcloud"
	"o.o/backend/pkg/common/storage"
)

const (
	ReadOnly    = `https://www.googleapis.com/auth/devstorage.read_only`
	ReadWrite   = `https://www.googleapis.com/auth/devstorage.read_write`
	FullControl = `https://www.googleapis.com/auth/devstorage.full_control`
)

type Config struct {
	Bucket string `yaml:"bucket"`

	GoogleAuthFile string `yaml:"google_auth_file"`
}

var _ storage.Bucket = &Driver{}

type Driver struct {
	bucket     *gstorage.BucketHandle
	client     *gstorage.Client
	projectID  string
	bucketName string
}

var registry = map[string]registryItem{}
var mu sync.Mutex

type registryItem struct {
	client *gstorage.Client
}

func Connect(ctx context.Context, cfg Config, scopes []string) (*Driver, error) {
	if cfg.GoogleAuthFile == "" || cfg.Bucket == "" {
		return nil, cm.Errorf(cm.Internal, nil, "invalid config")
	}

	mu.Lock()
	defer mu.Unlock()

	// load credentials
	cred, err := gcloud.LoadCredentialsFromFile(ctx, cfg.GoogleAuthFile, scopes)
	if err != nil {
		return nil, cm.Wrap(err)
	}

	// reuse the client if they share the same credentials
	key_ := sha256.Sum256(cred.JSON)
	key := string(key_[:])

	regClient := registry[key]
	client := regClient.client
	if client == nil {
		client, err = gstorage.NewClient(ctx, option.WithCredentials(cred))
		if err != nil {
			return nil, cm.Wrap(err)
		}
		registry[key] = registryItem{client: client}
	}

	bucket := client.Bucket(cfg.Bucket)
	_, err = bucket.Attrs(ctx)
	if err != nil {
		return nil, cm.Wrap(err)
	}
	return &Driver{
		bucket:     bucket,
		client:     client,
		projectID:  cred.ProjectID,
		bucketName: cfg.Bucket,
	}, nil
}

func (b *Driver) Name() string {
	return b.bucketName
}

func (b *Driver) OpenForWrite(ctx context.Context, path string) (io.WriteCloser, error) {
	obj := b.bucket.Object(path)
	return obj.NewWriter(ctx), nil
}

func (b *Driver) OpenForRead(ctx context.Context, path string) (io.ReadCloser, error) {
	obj := b.bucket.Object(path)
	return obj.NewReader(ctx)
}
