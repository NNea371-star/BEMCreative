package repository

import (
	"BE/internal/domain"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *UserRepository) FindByID(id string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("id = ?", id).First(&user).Error
	return &user, err
}

func (r *UserRepository) Update(user *domain.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) Create(user *domain.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) FindAll() ([]domain.User, error) {
	var users []domain.User
	err := r.db.Select("id, name, email, role, created_at, updated_at").
		Order("created_at DESC").Find(&users).Error
	return users, err
}

func (r *UserRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&domain.User{}).Error
}