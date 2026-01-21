package storageService

import (
	"context"
	"errors"
	"mime/multipart"
	"path/filepath"
	"strings"

	"go-boilerplate/apps/internal/utils"

	"github.com/minio/minio-go/v7"
)

type FileType string

const (
	FileImage FileType = "image"
	FilePDF   FileType = "pdf"
)

type StorageService struct {
	client *minio.Client
	bucket string
}

func NewStorageService(
	client *minio.Client,
	bucket string,
) *StorageService {
	return &StorageService{
		client: client,
		bucket: bucket,
	}
}

// ========================
// Upload (DINAMIS)
// ========================

func (s *StorageService) Upload(
	ctx context.Context,
	file *multipart.FileHeader,
	folder string,
	fileType FileType,
) (string, error) {

	if err := validateFile(file, fileType); err != nil {
		return "", err
	}

	return utils.UploadToMinio(
		ctx,
		s.client,
		s.bucket,
		file,
		folder,
	)
}

// ========================
// Delete
// ========================

func (s *StorageService) Delete(
	ctx context.Context,
	objectName string,
) error {
	return utils.DeleteFromMinio(
		ctx,
		s.client,
		s.bucket,
		objectName,
	)
}

// ========================
// Validator
// ========================

func validateFile(
	file *multipart.FileHeader,
	fileType FileType,
) error {

	ext := strings.ToLower(filepath.Ext(file.Filename))

	switch fileType {

	case FileImage:
		allowed := map[string]bool{
			".png":  true,
			".jpg":  true,
			".jpeg": true,
			".webp": true,
		}

		if !allowed[ext] {
			return errors.New("hanya file png, jpg, jpeg, webp yang diizinkan")
		}

		if file.Size > 3<<20 {
			return errors.New("ukuran gambar maksimal 5MB")
		}

	case FilePDF:
		if ext != ".pdf" {
			return errors.New("hanya file PDF yang diizinkan")
		}

		if file.Size > 5<<20 {
			return errors.New("ukuran PDF maksimal 10MB")
		}

	default:
		return errors.New("tipe file tidak dikenali")
	}

	return nil
}
