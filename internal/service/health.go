package service

import "github.com/soundmarket/backend/internal/config"

type HealthService struct {
	cfg *config.Config
}

func NewHealthService(cfg *config.Config) *HealthService {
	return &HealthService{cfg: cfg}
}

func (s *HealthService) Status() map[string]string {
	return map[string]string{
		"status": "ok",
		"env":    s.cfg.AppEnv,
	}
}
