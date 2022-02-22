package ns3

import (
	"context"
	"github.com/minio/minio-go"
	"github.com/nbs-go/nlogger"
	"io"
)

var log = nlogger.Get()

type MinioOpt struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool
	BucketName      string
	ctx             *context.Context
}

func NewMinio(opt MinioOpt) (*Minio, error) {
	// Init minio client
	client, err := minio.New(opt.Endpoint, opt.AccessKeyID, opt.SecretAccessKey, opt.UseSSL)
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
		log.Error("error when upload object", nlogger.Error(err))
		return err
	}

	return nil
}

func (m *Minio) Remove(objectName string) error {
	err := m.Client.RemoveObject(m.BucketName, objectName)
	if err != nil {
		log.Error("error when removing object", nlogger.Error(err))
		return err
	}

	return nil
}
