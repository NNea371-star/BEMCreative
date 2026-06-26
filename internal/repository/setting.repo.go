package repository

import (
	"BE/internal/domain"
	"gorm.io/gorm"
)

type SettingRepository struct {
	db *gorm.DB
}

func NewSettingRepository(db *gorm.DB) *SettingRepository {
	return &SettingRepository{db: db}
}

func (r *SettingRepository) FindAll() ([]domain.SiteSetting, error) {
	var settings []domain.SiteSetting
	err := r.db.Find(&settings).Error
	return settings, err
}

func (r *SettingRepository) Upsert(key, value string) error {
	// Cari record berdasarkan key
	var setting domain.SiteSetting
	err := r.db.Where("key = ?", key).First(&setting).Error
	
	if err == gorm.ErrRecordNotFound {
		// Jika tidak ditemukan, buat baru
		return r.db.Create(&domain.SiteSetting{
			Key:   key,
			Value: value,
		}).Error
	}

	if err != nil {
		return err
	}

	// Update value (termasuk jika value kosong)
	return r.db.Model(&setting).Update("value", value).Error
}

func (r *SettingRepository) GetWhatsappConfig() (*domain.WhatsappConfig, error) {
	var config domain.WhatsappConfig
	err := r.db.First(&config).Error
	return &config, err
}

func (r *SettingRepository) UpdateWhatsappConfig(config *domain.WhatsappConfig) error {
	return r.db.Save(config).Error
}