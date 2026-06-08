package postgres

import (
	"context"

	engine "github.com/NicolasPaterno/warden-engine"
	db "github.com/NicolasPaterno/warden-engine/db/generated"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AlertRepo struct {
	queries *db.Queries
}

func NewAlertRepo(pool *pgxpool.Pool) *AlertRepo {
	return &AlertRepo{
		queries: db.New(pool),
	}
}

func (r *AlertRepo) Save(ctx context.Context, alert engine.Alert) error
func (r *AlertRepo) GetByRoom(ctx context.Context, room string) ([]engine.Alert, error)
func (r *AlertRepo) Save(ctx context.Context, alert engine.Alert) (engine.Alert, error) {
	created, err := r.queries.CreateAlert(ctx, db.CreateAlertParams{
		ID:        alert.ID,
		RuleID:    alert.RuleID,
		Room:      alert.Room,
		Message:   alert.Message,
		Severity:  toDBAlertSeverity(alert.Severity),
		Value:     alert.Value,
		CreatedAt: pgtype.Timestamptz{Time: alert.CreatedAt, Valid: true},
	})
	if err != nil {
		return engine.Alert{}, err
	}
	return toEngineAlert(created), nil
}
func (r *AlertRepo) GetAll(ctx context.Context) ([]engine.Alert, error)

func toEngineAlert(a db.Alert) engine.Alert {
	return engine.Alert{
		ID:        a.ID,
		RuleID:    a.RuleID,
		Room:      a.Room,
		Message:   a.Message,
		Severity:  engine.AlertSeverity(a.Severity),
		Value:     a.Value,
		CreatedAt: a.CreatedAt.Time,
	}
}

func toDBAlertSeverity(s engine.AlertSeverity) db.AlertSeverity {
	return db.AlertSeverity(s)
}
