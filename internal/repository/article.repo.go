package repository

import (
	"BE/internal/domain"

	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

type ArticleRepository struct {
	db *gorm.DB
}

func NewArticleRepository(db *gorm.DB) *ArticleRepository {
	return &ArticleRepository{db: db}
}

func (r *ArticleRepository) FindAll(status string) ([]domain.Article, error) {
	var articles []domain.Article
	q := r.db.Preload("Author").Preload("Category").Order("created_at DESC")
	if status != "" {
		q = q.Where("status = ?", status)
	}
	err := q.Find(&articles).Error
	return articles, err
}

func (r *ArticleRepository) FindPublished() ([]domain.Article, error) {
	var articles []domain.Article
	err := r.db.Preload("Author").Preload("Category").
		Where("status = ?", "published").
		Order("published_at DESC").
		Find(&articles).Error
	return articles, err
}

func (r *ArticleRepository) FindBySlug(slug string) (*domain.Article, error) {
	var article domain.Article
	err := r.db.Preload("Author").Preload("Category").
		Where("slug = ? AND status = ?", slug, "published").
		First(&article).Error
	return &article, err
}

func (r *ArticleRepository) FindByID(id string) (*domain.Article, error) {
	var article domain.Article
	err := r.db.Preload("Author").Preload("Category").
		Where("id = ?", id).First(&article).Error
	return &article, err
}

func (r *ArticleRepository) Create(a *domain.Article) error {
	return r.db.Create(a).Error
}

func (r *ArticleRepository) Update(a *domain.Article) error {
	return r.db.Save(a).Error
}

func (r *ArticleRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&domain.Article{}).Error
}

func (r *ArticleRepository) GenerateSlug(title string) string {
	base := slug.Make(title)
	s := base
	var count int64
	for i := 1; ; i++ {
		r.db.Model(&domain.Article{}).Where("slug = ?", s).Count(&count)
		if count == 0 {
			break
		}
		s = base + "-" + string(rune('0'+i))
	}
	return s
}

// Category
func (r *ArticleRepository) FindAllCategories() ([]domain.BlogCategory, error) {
	var cats []domain.BlogCategory
	err := r.db.Order("name ASC").Find(&cats).Error
	return cats, err
}

func (r *ArticleRepository) CreateCategory(name string) (*domain.BlogCategory, error) {
	cat := &domain.BlogCategory{
		Name: name,
		Slug: slug.Make(name),
	}
	return cat, r.db.Create(cat).Error
}

func (r *ArticleRepository) DeleteCategory(id string) error {
	var count int64
	r.db.Model(&domain.Article{}).Where("category_id = ?", id).Count(&count)
	if count > 0 {
		return gorm.ErrRecordNotFound
	}
	return r.db.Where("id = ?", id).Delete(&domain.BlogCategory{}).Error
}

func (r *ArticleRepository) IsCategoryUsed(id string) bool {
	var count int64
	r.db.Model(&domain.Article{}).Where("category_id = ?", id).Count(&count)
	return count > 0
}