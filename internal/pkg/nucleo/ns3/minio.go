package ns3

import (
	"context"
	"github.com/minio/minio-go"
	"io"
)

type MinioOpt struct {
	Endpoint        string
	AccessKeyId     string
	SecretAccessKey string
	UseSSL          bool
	BucketName      string
	ctx             *context.Context
}

func NewMinio(opt MinioOpt) (*Minio, error) {
	// Init minio client
	client, err := minio.New(opt.Endpoint, opt.AccessKeyId, opt.SecretAccessKey, opt.UseSSL)
	if err != nil {
		return nil, err
	}

	// Init s3 minio
	return &Minio{
		Client:     client,
		BucketName: opt.BucketName,
	}, nil
}

type Minio struct {
	Client     *minio.Client
	BucketName string
	Region     string
}

func (m *Minio) Upload(file io.Reader, contentType, dest string) error {
	_, err := m.Client.PutObject(m.BucketName, dest, file, -1, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return err
	}

	return nil
}
