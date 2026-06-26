package service

import (
	"BE/internal/domain"
	"BE/internal/repository"
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetByID(id string) (*domain.User, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("user tidak ditemukan")
	}
	return user, nil
}

func (s *UserService) UpdateProfile(id, name, email string) (*domain.User, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("user tidak ditemukan")
	}

	if name == "" || email == "" {
		return nil, errors.New("nama dan email wajib diisi")
	}

	// Cek email duplikat
	existing, err := s.repo.FindByEmail(email)
	if err == nil && existing.ID != user.ID {
		return nil, errors.New("email sudah digunakan")
	}

	user.Name = name
	user.Email = email

	if err := s.repo.Update(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) ChangePassword(id, currentPassword, newPassword string) error {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("user tidak ditemukan")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(currentPassword)); err != nil {
		return errors.New("password saat ini tidak benar")
	}

	if len(newPassword) < 6 {
		return errors.New("password baru minimal 6 karakter")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.PasswordHash = string(hash)
	return s.repo.Update(user)
}

func (s *UserService) Create(name, email, password, role string) (*domain.User, error) {
	if name == "" || email == "" || password == "" {
		return nil, errors.New("nama, email, dan password wajib diisi")
	}

	// Cek email duplikat
	if _, err := s.repo.FindByEmail(email); err == nil {
		return nil, errors.New("email sudah terdaftar")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	if role == "" {
		role = "admin"
	}

	user := &domain.User{
		ID:           uuid.New(),
		Name:         name,
		Email:        email,
		PasswordHash: string(hash),
		Role:         role,
	}

	return user, s.repo.Create(user)
}