package service

import (
	"BE/internal/domain"
	"BE/internal/repository"
	"errors"
	"time"

	"github.com/google/uuid"
)

type ArticleService struct {
	repo          *repository.ArticleRepository
	uploadService *UploadService
}

func NewArticleService(repo *repository.ArticleRepository, uploadService *UploadService) *ArticleService {
	return &ArticleService{repo: repo, uploadService: uploadService}
}

func (s *ArticleService) GetAll(status string) ([]domain.Article, error) {
	return s.repo.FindAll(status)
}

func (s *ArticleService) GetPublished() ([]domain.Article, error) {
	return s.repo.FindPublished()
}

func (s *ArticleService) GetBySlug(slug string) (*domain.Article, error) {
	a, err := s.repo.FindBySlug(slug)
	if err != nil {
		return nil, errors.New("artikel tidak ditemukan")
	}
	return a, nil
}

func (s *ArticleService) Create(input map[string]any, authorID string) (*domain.Article, error) {
	catID, err := uuid.Parse(getString(input, "category_id"))
	if err != nil {
		return nil, errors.New("category_id tidak valid")
	}
	aID, _ := uuid.Parse(authorID)

	status := getString(input, "status")
	if status == "" {
		status = "draft"
	}

	var publishedAt *time.Time
	if status == "published" {
		now := time.Now()
		publishedAt = &now
	}

	article := &domain.Article{
		ID:           uuid.New(),
		AuthorID:     aID,
		CategoryID:   catID,
		Title:        getString(input, "title"),
		Slug:         s.repo.GenerateSlug(getString(input, "title")),
		Content:      getString(input, "content"),
		ThumbnailURL: getString(input, "thumbnail_url"),
		Status:       status,
		PublishedAt:  publishedAt,
	}

	if article.Title == "" {
		return nil, errors.New("judul artikel wajib diisi")
	}

	return article, s.repo.Create(article)
}

func (s *ArticleService) Update(id string, input map[string]any) (*domain.Article, error) {
	article, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("artikel tidak ditemukan")
	}

	if v, ok := input["title"].(string); ok { article.Title = v }
	if v, ok := input["content"].(string); ok { article.Content = v }
	if v, ok := input["thumbnail_url"].(string); ok { article.ThumbnailURL = v }
	if v, ok := input["category_id"].(string); ok {
		if catID, err := uuid.Parse(v); err == nil {
			article.CategoryID = catID
		}
	}
	if v, ok := input["status"].(string); ok {
		if article.Status != "published" && v == "published" {
			now := time.Now()
			article.PublishedAt = &now
		}
		article.Status = v
	}

	return article, s.repo.Update(article)
}

func (s *ArticleService) Delete(id string) error {
	article, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("artikel tidak ditemukan")
	}

	if article.ThumbnailURL != "" {
		publicID := ExtractPublicID(article.ThumbnailURL)
		if publicID != "" {
			_ = s.uploadService.DeleteImage(publicID)
		}
	}

	return s.repo.Delete(id)
}

func (s *ArticleService) GetAllCategories() ([]domain.BlogCategory, error) {
	return s.repo.FindAllCategories()
}

func (s *ArticleService) CreateCategory(name string) (*domain.BlogCategory, error) {
	return s.repo.CreateCategory(name)
}

func (s *ArticleService) DeleteCategory(id string) error {
	if s.repo.IsCategoryUsed(id) {
		return errors.New("kategori masih digunakan oleh artikel, hapus artikelnya terlebih dahulu")
	}
	return s.repo.DeleteCategory(id)
}