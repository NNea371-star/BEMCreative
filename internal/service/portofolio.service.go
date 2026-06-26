package service

import (
	"BE/internal/domain"
	"BE/internal/repository"
	"errors"

	"github.com/google/uuid"
)

type PortfolioService struct {
	repo          *repository.PortfolioRepository
	uploadService *UploadService
}

func NewPortfolioService(repo *repository.PortfolioRepository, uploadService *UploadService) *PortfolioService {
	return &PortfolioService{repo: repo, uploadService: uploadService}
}

func (s *PortfolioService) GetAll() ([]domain.Portfolio, error) {
	return s.repo.FindAll()
}

func (s *PortfolioService) GetByID(id string) (*domain.Portfolio, error) {
	return s.repo.FindByID(id)
}

func (s *PortfolioService) Create(input map[string]any) (*domain.Portfolio, error) {
	p := &domain.Portfolio{
		ID:          uuid.New(),
		Title:       getString(input, "title"),
		Description: getString(input, "description"),
		ImageURL:    getString(input, "image_url"),
		ClientName:  getString(input, "client_name"),
		Year:        getInt(input, "year"),
	}
	if p.Title == "" || p.ClientName == "" {
		return nil, errors.New("judul dan nama klien wajib diisi")
	}
	return p, s.repo.Create(p)
}

func (s *PortfolioService) Update(id string, input map[string]any) (*domain.Portfolio, error) {
	p, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("karya tidak ditemukan")
	}
	if v, ok := input["title"].(string); ok { p.Title = v }
	if v, ok := input["description"].(string); ok { p.Description = v }
	if v, ok := input["image_url"].(string); ok { p.ImageURL = v }
	if v, ok := input["client_name"].(string); ok { p.ClientName = v }
	if v, ok := input["year"].(float64); ok { p.Year = int(v) }
	return p, s.repo.Update(p)
}

func (s *PortfolioService) Delete(id string) error {
	p, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("karya tidak ditemukan")
	}

	if p.ImageURL != "" {
		publicID := ExtractPublicID(p.ImageURL)
		if publicID != "" {
			_ = s.uploadService.DeleteImage(publicID)
		}
	}

	return s.repo.Delete(id)
}