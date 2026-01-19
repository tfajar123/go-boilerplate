package database

import (
	"context"
	"go-boilerplate/apps/internal/config"
	"log"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func NewMinioClient(cfg config.MinioConfig) *minio.Client {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func EnsureBucket(
	ctx context.Context,
	client *minio.Client,
	bucket string,
) {
	exists, err := client.BucketExists(ctx, bucket)
	if err != nil {
		log.Fatal(err)
	}

	if !exists {
		err = client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
		if err != nil && !strings.Contains(err.Error(), "BucketAlreadyOwnedByYou") {
			log.Fatal(err)
		}
	}
}
