package repository

import (
	"BE/internal/domain"
	"gorm.io/gorm"
)

type PortfolioRepository struct {
	db *gorm.DB
}

func NewPortfolioRepository(db *gorm.DB) *PortfolioRepository {
	return &PortfolioRepository{db: db}
}

func (r *PortfolioRepository) FindAll() ([]domain.Portfolio, error) {
	var portfolios []domain.Portfolio
	err := r.db.Order("year DESC, created_at DESC").Find(&portfolios).Error
	return portfolios, err
}

func (r *PortfolioRepository) FindByID(id string) (*domain.Portfolio, error) {
	var portfolio domain.Portfolio
	err := r.db.Where("id = ?", id).First(&portfolio).Error
	return &portfolio, err
}

func (r *PortfolioRepository) Create(p *domain.Portfolio) error {
	return r.db.Create(p).Error
}

func (r *PortfolioRepository) Update(p *domain.Portfolio) error {
	return r.db.Save(p).Error
}

func (r *PortfolioRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&domain.Portfolio{}).Error
}