package utils

import (
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"

	"go-boilerplate/apps/internal/config"
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

func BuildStorageURL(cfg *config.Config, object string) string {
	return fmt.Sprintf("%s/%s/%s", cfg.Storage.PublicURL, cfg.Storage.Bucket, object)
}

// // Contoh penggunaan
// obj, stat, err := utils.GetFileFromStorage(ctx, storage, objectPath)
// if err != nil {
//     // handle error
// }
// defer obj.Close() // JANGAN LUPA close object setelah selesai pakai

// // Dapatkan metadata:
// fmt.Println("Ukuran file:", stat.Size)
// fmt.Println("Content Type:", stat.ContentType)

// // Baca isi file:
// data, err := io.ReadAll(obj)

func GetFileFromStorage(
	ctx context.Context,
	storage *database.Storage,
	object string,
) (*minio.Object, minio.ObjectInfo, error) {
	obj, err := storage.Client.GetObject(
		ctx,
		storage.Bucket,
		object,
		minio.GetObjectOptions{},
	)
	if err != nil {
		return nil, minio.ObjectInfo{}, err
	}

	stat, err := obj.Stat()
	if err != nil {
		obj.Close()
		return nil, minio.ObjectInfo{}, err
	}

	return obj, stat, nil
}
