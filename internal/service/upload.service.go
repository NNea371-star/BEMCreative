package service

import (
	"context"
	"errors"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type UploadService struct {
	cld *cloudinary.Cloudinary
}

func NewUploadService(cloudName, apiKey, apiSecret string) (*UploadService, error) {
	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		return nil, err
	}
	return &UploadService{cld: cld}, nil
}

func (s *UploadService) UploadImage(file multipart.File, header *multipart.FileHeader, folder string) (string, error) {
	// Validasi ekstensi
	ext := strings.ToLower(filepath.Ext(header.Filename))
	allowed := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true}
	if !allowed[ext] {
		return "", errors.New("format file tidak didukung, gunakan JPG, PNG, atau WebP")
	}

	// Validasi ukuran — max 5MB
	if header.Size > 5*1024*1024 {
		return "", errors.New("ukuran file maksimal 5MB")
	}

	ctx := context.Background()

	// Upload ke Cloudinary
	resp, err := s.cld.Upload.Upload(ctx, file, uploader.UploadParams{
		Folder:         "mcreative/" + folder,
		Transformation: "f_auto,q_auto,w_1200",
	})
	if err != nil {
		return "", errors.New("gagal mengupload gambar: " + err.Error())
	}

	return resp.SecureURL, nil
}

func (s *UploadService) DeleteImage(publicID string) error {
	ctx := context.Background()
	_, err := s.cld.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: publicID,
	})
	return err
}

// Helper — ambil public ID dari URL Cloudinary
func ExtractPublicID(url string) string {
	// https://res.cloudinary.com/demo/image/upload/v123/mcreative/products/abc.jpg
	// → mcreative/products/abc
	parts := strings.Split(url, "/upload/")
	if len(parts) < 2 {
		return ""
	}
	path := parts[1]
	// Hapus versi jika ada (v123456/)
	if strings.HasPrefix(path, "v") {
		idx := strings.Index(path, "/")
		if idx != -1 {
			path = path[idx+1:]
		}
	}
	// Hapus ekstensi
	ext := filepath.Ext(path)
	return strings.TrimSuffix(path, ext)
}