package service

import (
	"BE/internal/domain"
	"BE/internal/repository"
	"errors"

	"github.com/google/uuid"
)

type TestimonialService struct {
	repo *repository.TestimonialRepository
}

func NewTestimonialService(repo *repository.TestimonialRepository) *TestimonialService {
	return &TestimonialService{repo: repo}
}

func (s *TestimonialService) GetVisible() ([]domain.Testimonial, error) {
	return s.repo.FindVisible()
}

func (s *TestimonialService) GetAll() ([]domain.Testimonial, error) {
	return s.repo.FindAll()
}

func (s *TestimonialService) Create(input map[string]any) (*domain.Testimonial, error) {
	t := &domain.Testimonial{
		ID:            uuid.New(),
		ClientName:    getString(input, "client_name"),
		ClientCompany: getString(input, "client_company"),
		Message:       getString(input, "message"),
		Rating:        getInt(input, "rating"),
		IsVisible:     getBool(input, "is_visible", false),
	}
	if t.ClientName == "" || t.Message == "" {
		return nil, errors.New("nama dan pesan wajib diisi")
	}
	if t.Rating < 1 || t.Rating > 5 {
		t.Rating = 5
	}
	return t, s.repo.Create(t)
}

func (s *TestimonialService) Update(id string, input map[string]any) (*domain.Testimonial, error) {
	t, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("testimoni tidak ditemukan")
	}
	if v, ok := input["client_name"].(string); ok { t.ClientName = v }
	if v, ok := input["client_company"].(string); ok { t.ClientCompany = v }
	if v, ok := input["message"].(string); ok { t.Message = v }
	if v, ok := input["rating"].(float64); ok { t.Rating = int(v) }
	if v, ok := input["is_visible"].(bool); ok { t.IsVisible = v }
	return t, s.repo.Update(t)
}

func (s *TestimonialService) Delete(id string) error {
	return s.repo.Delete(id)
}