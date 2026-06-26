package service

import (
	"BE/internal/domain"
	"BE/internal/repository"
)

type SettingService struct {
	settingRepo *repository.SettingRepository
	productRepo *repository.ProductRepository
}

func NewSettingService(
	settingRepo *repository.SettingRepository,
	productRepo *repository.ProductRepository,
) *SettingService {
	return &SettingService{
		settingRepo: settingRepo,
		productRepo: productRepo,
	}
}

func (s *SettingService) GetPublicSettings() ([]domain.SiteSetting, error) {
	return s.settingRepo.FindAll()
}

func (s *SettingService) UpdateSettings(settings []struct {
	Key   string
	Value string
}) error {
	for _, item := range settings {
		if err := s.settingRepo.Upsert(item.Key, item.Value); err != nil {
			return err
		}
	}
	return nil
}

func (s *SettingService) GetWhatsappConfig() (*domain.WhatsappConfig, error) {
	return s.settingRepo.GetWhatsappConfig()
}

func (s *SettingService) UpdateWhatsappConfig(phone, orderTemplate, hireTemplate string) error {
	config, err := s.settingRepo.GetWhatsappConfig()
	if err != nil {
		return err
	}
	config.PhoneNumber = phone
	config.OrderTemplate = orderTemplate
	config.HireTemplate = hireTemplate
	return s.settingRepo.UpdateWhatsappConfig(config)
}