package engine

import "context"

type AlertRepository interface {
	Save(ctx context.Context, alert Alert) error
	GetByRoom(ctx context.Context, room string) ([]Alert, error)
	GetByRule(ctx context.Context, ruleID string) ([]Alert, error)
	GetAll(ctx context.Context) ([]Alert, error)
}
