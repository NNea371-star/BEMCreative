package service

import (
	"BE/internal/domain"
	"BE/internal/repository"
)

type StatsService interface {
	GetSiteStats() (*domain.SiteStats, error)
}

type statsService struct {
	statsRepo repository.StatsRepository
}

func NewStatsService(statsRepo repository.StatsRepository) StatsService {
	return &statsService{
		statsRepo: statsRepo,
	}
}

func (s *statsService) GetSiteStats() (*domain.SiteStats, error) {
	return s.statsRepo.GetSiteStats()
}