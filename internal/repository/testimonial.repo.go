package repository

import (
	"BE/internal/domain"
	"gorm.io/gorm"
)

type TestimonialRepository struct {
	db *gorm.DB
}

func NewTestimonialRepository(db *gorm.DB) *TestimonialRepository {
	return &TestimonialRepository{db: db}
}

func (r *TestimonialRepository) FindVisible() ([]domain.Testimonial, error) {
	var testimonials []domain.Testimonial
	err := r.db.Where("is_visible = ?", true).Order("created_at DESC").Find(&testimonials).Error
	return testimonials, err
}

func (r *TestimonialRepository) FindAll() ([]domain.Testimonial, error) {
	var testimonials []domain.Testimonial
	err := r.db.Order("created_at DESC").Find(&testimonials).Error
	return testimonials, err
}

func (r *TestimonialRepository) Create(t *domain.Testimonial) error {
	return r.db.Create(t).Error
}

func (r *TestimonialRepository) Update(t *domain.Testimonial) error {
	return r.db.Save(t).Error
}

func (r *TestimonialRepository) FindByID(id string) (*domain.Testimonial, error) {
	var t domain.Testimonial
	err := r.db.Where("id = ?", id).First(&t).Error
	return &t, err
}

func (r *TestimonialRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&domain.Testimonial{}).Error
}