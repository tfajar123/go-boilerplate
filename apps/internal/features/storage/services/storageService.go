package storageService

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"

	"go-boilerplate/apps/internal/database"
	"go-boilerplate/apps/internal/utils"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
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

// UploadBase64 uploads a file from a base64 encoded string to storage.
// Supports base64 string with or without data URI prefix (e.g. data:image/png;base64,...).
func (s *StorageService) UploadBase64(
	ctx context.Context,
	base64String string,
	folder string,
) (string, error) {
	var contentType string

	// Handle data URI format
	if strings.Contains(base64String, ";base64,") {
		parts := strings.SplitN(base64String, ";base64,", 2)
		contentType = strings.TrimPrefix(parts[0], "data:")
		base64String = parts[1]
	}

	// Decode base64
	data, err := base64.StdEncoding.DecodeString(base64String)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}

	// Detect content type if not provided via data URI
	if contentType == "" {
		contentType = http.DetectContentType(data)
	}

	// Get file extension from mime type
	ext, err := extensionFromMimeType(contentType)
	if err != nil {
		return "", err
	}

	// Generate unique object name
	objectName := fmt.Sprintf("%s/%s%s", folder, uuid.NewString(), ext)

	// Upload to storage
	_, err = s.storage.Client.PutObject(
		ctx,
		s.storage.Bucket,
		objectName,
		bytes.NewReader(data),
		int64(len(data)),
		minio.PutObjectOptions{
			ContentType: contentType,
		},
	)
	if err != nil {
		return "", fmt.Errorf("failed to upload base64 file: %w", err)
	}

	return objectName, nil
}

// extensionFromMimeType maps common mime types to file extensions.
func extensionFromMimeType(mimeType string) (string, error) {
	mimeToExt := map[string]string{
		"image/jpeg": ".jpg",
		"image/jpg":  ".jpg",
		"image/png":  ".png",
		"image/gif":  ".gif",
		"image/webp": ".webp",
		"application/pdf": ".pdf",
		"application/msword": ".doc",
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document": ".docx",
		"text/plain": ".txt",
	}

	if ext, ok := mimeToExt[mimeType]; ok {
		return ext, nil
	}

	return "", fmt.Errorf("unsupported file type: %s", mimeType)
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
