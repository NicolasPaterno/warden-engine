package engine

import "context"

type RuleRepository interface {
	Save(ctx context.Context, rule Rule) error
	GetByID(ctx context.Context, id string) (Rule, error)
	GetAll(ctx context.Context) ([]Rule, error)
	Update(ctx context.Context, rule Rule) error
	Delete(ctx context.Context, id string) error
}
