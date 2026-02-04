package storageService

import (
	"context"
	"errors"
	"mime/multipart"
	"path/filepath"
	"strings"

	"go-boilerplate/apps/internal/database"
	"go-boilerplate/apps/internal/utils"
)

type FileType string

const (
	FileImage FileType = "image"
	FilePDF   FileType = "pdf"
)

type StorageService struct {
	storage *database.Storage
}

func NewStorageService(storage *database.Storage) *StorageService {
	return &StorageService{
		storage: storage,
	}
}

func (s *StorageService) Upload(
	ctx context.Context,
	file *multipart.FileHeader,
	folder string,
	fileType FileType,
) (string, error) {

	if err := validateFile(file, fileType); err != nil {
		return "", err
	}

	return utils.UploadToStorage(
		ctx,
		s.storage,
		file,
		folder,
	)
}

func (s *StorageService) Delete(
	ctx context.Context,
	objectName string,
) error {

	return utils.DeleteFromStorage(
		ctx,
		s.storage,
		objectName,
	)
}

func validateFile(
	file *multipart.FileHeader,
	fileType FileType,
) error {

	ext := strings.ToLower(filepath.Ext(file.Filename))
	contentType := file.Header.Get("Content-Type")

	switch fileType {

	case FileImage:

		allowedExt := map[string]bool{
			".png":  true,
			".jpg":  true,
			".jpeg": true,
			".webp": true,
		}

		allowedMime := map[string]bool{
			"image/png":  true,
			"image/jpeg": true,
			"image/webp": true,
		}

		if !allowedExt[ext] {
			return errors.New("hanya file png, jpg, jpeg, webp yang diizinkan")
		}

		if !allowedMime[contentType] {
			return errors.New("content-type gambar tidak valid")
		}

		if file.Size > 5<<20 { // 5MB
			return errors.New("ukuran gambar maksimal 5MB")
		}

	case FilePDF:

		if ext != ".pdf" {
			return errors.New("hanya file PDF yang diizinkan")
		}

		if contentType != "application/pdf" {
			return errors.New("content-type harus application/pdf")
		}

		if file.Size > 10<<20 { // 10MB
			return errors.New("ukuran PDF maksimal 10MB")
		}

	default:
		return errors.New("tipe file tidak dikenali")
	}

	return nil
}
