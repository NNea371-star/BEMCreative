package repository

import (
	"BE/internal/domain"

	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) FindAll() ([]domain.Product, error) {
	var products []domain.Product
	err := r.db.Preload("Category").Order("created_at DESC").Find(&products).Error
	return products, err
}

func (r *ProductRepository) FindByID(id string) (*domain.Product, error) {
	var product domain.Product
	err := r.db.Preload("Category").Where("id = ?", id).First(&product).Error
	return &product, err
}

func (r *ProductRepository) Create(product *domain.Product) error {
	return r.db.Create(product).Error
}

func (r *ProductRepository) Update(product *domain.Product) error {
	return r.db.Save(product).Error
}

func (r *ProductRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&domain.Product{}).Error
}

func (r *ProductRepository) UpdateStock(id string, stock int) error {
	return r.db.Model(&domain.Product{}).Where("id = ?", id).
		Updates(map[string]any{"stock": stock, "is_available": stock > 0}).Error
}

// Category
func (r *ProductRepository) FindAllCategories() ([]domain.ProductCategory, error) {
	var categories []domain.ProductCategory
	err := r.db.Order("name ASC").Find(&categories).Error
	return categories, err
}

func (r *ProductRepository) CreateCategory(name string) (*domain.ProductCategory, error) {
	cat := &domain.ProductCategory{
		Name: name,
		Slug: slug.Make(name),
	}
	err := r.db.Create(cat).Error
	return cat, err
}

func (r *ProductRepository) DeleteCategory(id string) error {
	// Cek apakah masih dipakai produk
	var count int64
	r.db.Model(&domain.Product{}).Where("category_id = ?", id).Count(&count)
	if count > 0 {
		return gorm.ErrRecordNotFound
	}
	return r.db.Where("id = ?", id).Delete(&domain.ProductCategory{}).Error
}

func (r *ProductRepository) IsCategoryUsed(id string) bool {
	var count int64
	r.db.Model(&domain.Product{}).Where("category_id = ?", id).Count(&count)
	return count > 0
}