package repository

import (
	"BE/database"
	"BE/internal/domain"
)

type StatsRepository interface {
	GetSiteStats() (*domain.SiteStats, error)
}

type statsRepository struct{}

func NewStatsRepository() StatsRepository {
	return &statsRepository{}
}

func (r *statsRepository) GetSiteStats() (*domain.SiteStats, error) {
	var stats domain.SiteStats

	var machineTypes int64
	if err := database.DB.Model(&domain.ProductCategory{}).Count(&machineTypes).Error; err != nil {
		return nil, err
	}
	stats.MachineTypes = int(machineTypes)

	var machinesBuilt int64
	if err := database.DB.Model(&domain.Product{}).Count(&machinesBuilt).Error; err != nil {
		return nil, err
	}
	stats.MachinesBuilt = int(machinesBuilt)

	var worksProduced int64
	if err := database.DB.Model(&domain.Portfolio{}).Count(&worksProduced).Error; err != nil {
		return nil, err
	}
	stats.WorksProduced = int(worksProduced)

	var firstPortfolio domain.Portfolio
	if err := database.DB.Order("year asc").First(&firstPortfolio).Error; err == nil {
		currentYear := 2026
		stats.ExperienceYears = currentYear - firstPortfolio.Year
	} else {
		stats.ExperienceYears = 3
	}

	return &stats, nil
}