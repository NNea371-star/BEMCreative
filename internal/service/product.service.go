package service

import (
	"BE/internal/domain"
	"BE/internal/hub"
	"BE/internal/repository"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ProductService struct {
	repo *repository.ProductRepository
	uploadService *UploadService
}

func NewProductService(repo *repository.ProductRepository, uploadService *UploadService) *ProductService {
	return &ProductService{repo: repo, uploadService: uploadService}
}

func (s *ProductService) GetAll() ([]domain.Product, error) {
	return s.repo.FindAll()
}

func (s *ProductService) GetByID(id string) (*domain.Product, error) {
	return s.repo.FindByID(id)
}

func (s *ProductService) Create(input map[string]any) (*domain.Product, error) {
	catID, err := uuid.Parse(input["category_id"].(string))
	if err != nil {
		return nil, errors.New("category_id tidak valid")
	}

	product := &domain.Product{
		ID:          uuid.New(),
		CategoryID:  catID,
		Name:        input["name"].(string),
		Description: getString(input, "description"),
		Price:       getFloat(input, "price"),
		Stock:       getInt(input, "stock"),
		ImageURL:    getString(input, "image_url"),
		IsAvailable: getBool(input, "is_available", true),
	}

	return product, s.repo.Create(product)
}

func (s *ProductService) Update(id string, input map[string]any) (*domain.Product, error) {
	product, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("produk tidak ditemukan")
	}

	if v, ok := input["name"].(string); ok {
		product.Name = v
	}
	if v, ok := input["description"].(string); ok {
		product.Description = v
	}
	if v, ok := input["price"].(float64); ok {
		product.Price = v
	}
	if v, ok := input["stock"].(float64); ok {
		product.Stock = int(v)
	}
	if v, ok := input["image_url"].(string); ok {
		product.ImageURL = v
	}
	if v, ok := input["is_available"].(bool); ok {
		product.IsAvailable = v
	}
	if v, ok := input["category_id"].(string); ok {
		catID, err := uuid.Parse(v)
		if err == nil {
			product.CategoryID = catID
		}
	}

	return product, s.repo.Update(product)
}

func (s *ProductService) UpdateStock(id string, stock int) error {
	err := s.repo.UpdateStock(id, stock)
	if err != nil {
		return err
	}
	// Broadcast ke semua pengunjung
	hub.H.BroadcastToAll(fiber.Map{
		"type":         "product:stock_updated",
		"product_id":   id,
		"stock":        stock,
		"is_available": stock > 0,
	})
	return nil
}

func (s *ProductService) GetAllCategories() ([]domain.ProductCategory, error) {
	return s.repo.FindAllCategories()
}

func (s *ProductService) CreateCategory(name string) (*domain.ProductCategory, error) {
	return s.repo.CreateCategory(name)
}

func (s *ProductService) DeleteCategory(id string) error {
	if s.repo.IsCategoryUsed(id) {
		return errors.New("kategori masih digunakan oleh produk, hapus produknya terlebih dahulu")
	}
	return s.repo.DeleteCategory(id)
}

// helpers
func getString(m map[string]any, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}

func getFloat(m map[string]any, key string) float64 {
	if v, ok := m[key].(float64); ok {
		return v
	}
	return 0
}

func getInt(m map[string]any, key string) int {
	if v, ok := m[key].(float64); ok {
		return int(v)
	}
	return 0
}

func getBool(m map[string]any, key string, def bool) bool {
	if v, ok := m[key].(bool); ok {
		return v
	}
	return def
}

func (s *ProductService) Delete(id string) error {
	product, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("produk tidak ditemukan")
	}

	// Hapus gambar dari Cloudinary jika ada
	if product.ImageURL != "" {
		publicID := ExtractPublicID(product.ImageURL)
		if publicID != "" {
			_ = s.uploadService.DeleteImage(publicID)
		}
	}

	return s.repo.Delete(id)
}