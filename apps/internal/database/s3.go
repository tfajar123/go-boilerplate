package database

import (
	"context"
	"go-boilerplate/apps/internal/config"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Storage struct {
	Client *minio.Client
	Bucket string
}

func NewStorage(cfg config.StorageConfig) *Storage {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
		Region: cfg.Region,
	})
	if err != nil {
		log.Fatal(err)
	}

	s := &Storage{
		Client: client,
		Bucket: cfg.Bucket,
	}

	s.ensureBucket(context.Background())

	return s
}

func (s *Storage) ensureBucket(ctx context.Context) {
	exists, err := s.Client.BucketExists(ctx, s.Bucket)
	if err != nil {
		log.Fatal(err)
	}

	if !exists {
		err = s.Client.MakeBucket(ctx, s.Bucket, minio.MakeBucketOptions{
			Region: "",
		})
		if err != nil {
			log.Fatal(err)
		}
	}
}
