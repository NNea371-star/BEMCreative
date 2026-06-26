package repository

import (
	"time"
	"BE/internal/domain"
	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(o *domain.OrderLog) error {
	return r.db.Create(o).Error
}

func (r *OrderRepository) FindAll() ([]domain.OrderLog, error) {
	var orders []domain.OrderLog
	err := r.db.Order("created_at DESC").Find(&orders).Error
	return orders, err
}

func (r *OrderRepository) FindAllWithPagination(limit, offset int) ([]domain.OrderLog, error) {
	var orders []domain.OrderLog
	err := r.db.Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&orders).Error
	return orders, err
}

func (r *OrderRepository) FindWithFilters(limit, offset int, status, search, from, to string) ([]domain.OrderLog, error) {
	query := r.db.Model(&domain.OrderLog{})

	// Filter status
	if status != "" && status != "all" {
		query = query.Where("status = ?", status)
	}

	// Filter search (name, phone, product_name)
	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where(
			"visitor_name ILIKE ? OR visitor_wa ILIKE ? OR product_name ILIKE ? OR project_type ILIKE ?",
			searchPattern, searchPattern, searchPattern, searchPattern,
		)
	}

	// Filter date range
	if from != "" {
		query = query.Where("created_at >= ?", from)
	}
	if to != "" {
		query = query.Where("created_at <= ?", to)
	}

	var orders []domain.OrderLog
	err := query.Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&orders).Error
	return orders, err
}

func (r *OrderRepository) FindByID(id string) (*domain.OrderLog, error) {
	var order domain.OrderLog
	err := r.db.Where("id = ?", id).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepository) FindRecent(limit int) ([]domain.OrderLog, error) {
	var orders []domain.OrderLog
	err := r.db.Order("created_at DESC").
		Limit(limit).
		Find(&orders).Error
	return orders, err
}

func (r *OrderRepository) UpdateStatus(id string, status string) error {
	return r.db.Model(&domain.OrderLog{}).
		Where("id = ?", id).
		Update("status", status).Error
}

func (r *OrderRepository) Update(order *domain.OrderLog) error {
	return r.db.Save(order).Error
}

func (r *OrderRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&domain.OrderLog{}).Error
}

func (r *OrderRepository) BulkDelete(ids []string) (int64, error) {
	result := r.db.Where("id IN ?", ids).Delete(&domain.OrderLog{})
	return result.RowsAffected, result.Error
}

func (r *OrderRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&domain.OrderLog{}).Count(&count).Error
	return count, err
}

func (r *OrderRepository) CountByStatus(status string) (int64, error) {
	var count int64
	err := r.db.Model(&domain.OrderLog{}).
		Where("status = ?", status).
		Count(&count).Error
	return count, err
}

func (r *OrderRepository) CountByDateRange(start, end time.Time) (int64, error) {
	var count int64
	err := r.db.Model(&domain.OrderLog{}).
		Where("created_at BETWEEN ? AND ?", start, end).
		Count(&count).Error
	return count, err
}

func (r *OrderRepository) CountWithFilters(status, search, from, to string) (int64, error) {
	query := r.db.Model(&domain.OrderLog{})

	// Filter status
	if status != "" && status != "all" {
		query = query.Where("status = ?", status)
	}

	// Filter search
	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where(
			"visitor_name ILIKE ? OR visitor_wa ILIKE ? OR product_name ILIKE ? OR project_type ILIKE ?",
			searchPattern, searchPattern, searchPattern, searchPattern,
		)
	}

	// Filter date range
	if from != "" {
		query = query.Where("created_at >= ?", from)
	}
	if to != "" {
		query = query.Where("created_at <= ?", to)
	}

	var count int64
	err := query.Count(&count).Error
	return count, err
}