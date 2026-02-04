package utils

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"

	"go-boilerplate/apps/internal/database"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

func UploadToStorage(
	ctx context.Context,
	storage *database.Storage,
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

	_, err = storage.Client.PutObject(
		ctx,
		storage.Bucket,
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

func DeleteFromStorage(
	ctx context.Context,
	storage *database.Storage,
	object string,
) error {

	return storage.Client.RemoveObject(
		ctx,
		storage.Bucket,
		object,
		minio.RemoveObjectOptions{},
	)
}
