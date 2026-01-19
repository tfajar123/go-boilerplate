package utils

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

// ========================
// Upload
// ========================

func UploadToMinio(
	ctx context.Context,
	client *minio.Client,
	bucket string,
	file *multipart.FileHeader,
	folder string,
) (string, error) {

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	ext := filepath.Ext(file.Filename)

	objectName := fmt.Sprintf(
		"%s/%s%s",
		folder,
		uuid.NewString(),
		ext,
	)

	_, err = client.PutObject(
		ctx,
		bucket,
		objectName,
		src,
		file.Size,
		minio.PutObjectOptions{
			ContentType: file.Header.Get("Content-Type"),
		},
	)

	if err != nil {
		return "", err
	}

	return objectName, nil
}

// ========================
// Delete
// ========================

func DeleteFromMinio(
	ctx context.Context,
	client *minio.Client,
	bucket string,
	object string,
) error {

	return client.RemoveObject(
		ctx,
		bucket,
		object,
		minio.RemoveObjectOptions{},
	)
}
