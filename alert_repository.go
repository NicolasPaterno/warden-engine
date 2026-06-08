package engine

import "context"

type AlertRepository interface {
	Save(ctx context.Context, alert Alert) (Alert, error)
	GetByRoom(ctx context.Context, room string) ([]Alert, error)
	GetAll(ctx context.Context) ([]Alert, error)
}
