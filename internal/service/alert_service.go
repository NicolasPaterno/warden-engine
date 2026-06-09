package service

import (
	"context"

	engine "github.com/NicolasPaterno/warden-engine"
)

type AlertService struct {
	repo engine.AlertRepository
}

func NewAlertService(repo engine.AlertRepository) *AlertService {
	return &AlertService{repo: repo}
}

func (s *AlertService) GetAll(ctx context.Context) ([]engine.Alert, error) {
	return s.repo.GetAll(ctx)
}
func (s *AlertService) GetByRoom(ctx context.Context, room string) ([]engine.Alert, error) {
	return s.repo.GetByRoom(ctx, room)
}
func (s *AlertService) GetByRule(ctx context.Context, ruleID string) ([]engine.Alert, error) {
	return s.repo.GetByRule(ctx, ruleID)
}
