package awsstorage

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/storage"
	"o.o/common/l"
)

var ll = l.New()

type Config struct {
	Bucket          string `yaml:"bucket"`
	Region          string `yaml:"region"`
	AccessKeyID     string `yaml:"access_key_id"`
	SecretAccessKey string `yaml:"secret_access_key"`
	EndpointURL     string `yaml:"endpoint_url"`
}

var _ storage.Bucket = &Driver{}

type Driver struct {
	bucket   string
	session  *session.Session
	uploader *s3manager.Uploader
}

func Connect(ctx context.Context, cfg Config) (*Driver, error) {
	if cfg.AccessKeyID == "" || cfg.SecretAccessKey == "" || cfg.Bucket == "" {
		return nil, cm.Errorf(cm.Internal, nil, "invalid config aws")
	}
	region := cfg.Region
	if region == "" {
		region = "ap-southeast-1"
	}
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
		Credentials: credentials.NewStaticCredentials(
			cfg.AccessKeyID,
			cfg.SecretAccessKey,
			"", // a token will be created when the session it's used.
		),
		Endpoint: aws.String(cfg.EndpointURL),
	})
	if err != nil {
		return nil, cm.Wrap(err)
	}

	uploader := s3manager.NewUploader(sess)
	return &Driver{
		bucket:   cfg.Bucket,
		session:  sess,
		uploader: uploader,
	}, nil

}

func (d *Driver) Name() string {
	return d.bucket
}

func (d *Driver) OpenForRead(ctx context.Context, path string) (io.ReadCloser, error) {
	panic("implement me")
}

func (d *Driver) OpenForWrite(ctx context.Context, path string) (io.WriteCloser, error) {
	pr, pw := io.Pipe()
	go func() {
		_, err := d.uploader.Upload(&s3manager.UploadInput{
			Bucket: aws.String(d.bucket),
			ACL:    aws.String("public-read"),
			Key:    aws.String(path),
			Body:   pr,
		})
		if err != nil {
			ll.Error("error creating file", l.Error(err))
			return
		}
	}()

	return pw, nil
}
