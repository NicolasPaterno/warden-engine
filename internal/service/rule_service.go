package service

import (
	"context"
	"time"

	engine "github.com/NicolasPaterno/warden-engine"
	"github.com/google/uuid"
)

type RuleService struct {
	repo engine.RuleRepository
}

func NewRuleService(repo engine.RuleRepository) *RuleService {
	return &RuleService{repo: repo}
}

func (s *RuleService) Create(ctx context.Context, rule engine.Rule) (engine.Rule, error) {
	rule.ID = uuid.NewString()
	rule.CreatedAt = time.Now()
	err := s.repo.Save(ctx, rule)
	if err != nil {
		return engine.Rule{}, err
	}
	return rule, nil
}

func (s *RuleService) GetByID(ctx context.Context, id string) (engine.Rule, error) {
	return s.repo.GetByID(ctx, id)
}
func (s *RuleService) GetAll(ctx context.Context) ([]engine.Rule, error) {
	return s.repo.GetAll(ctx)
}
func (s *RuleService) Update(ctx context.Context, rule engine.Rule) (engine.Rule, error) {
	return s.repo.Update(ctx, rule)
}
func (s *RuleService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
