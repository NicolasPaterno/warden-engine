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

func (r *AlertRepo) Save(ctx context.Context, alert engine.Alert) error {
	_, err := r.queries.CreateAlert(ctx, db.CreateAlertParams{
		ID:        alert.ID,
		RuleID:    alert.RuleID,
		Room:      alert.Room,
		Message:   alert.Message,
		Severity:  toDBAlertSeverity(alert.Severity),
		Value:     alert.Value,
		CreatedAt: pgtype.Timestamptz{Time: alert.CreatedAt, Valid: true},
	})
	return err
}

func (r *AlertRepo) GetByRoom(ctx context.Context, room string) ([]engine.Alert, error) {
	rows, err := r.queries.GetAlertsByRoom(ctx, room)
	if err != nil {
		return nil, err
	}
	var alerts []engine.Alert
	for _, row := range rows {
		alerts = append(alerts, toEngineAlert(row))
	}
	return alerts, nil
}

func (r *AlertRepo) GetAll(ctx context.Context) ([]engine.Alert, error) {
	rows, err := r.queries.GetAllAlerts(ctx)
	if err != nil {
		return nil, err
	}
	var alerts []engine.Alert
	for _, row := range rows {
		alerts = append(alerts, toEngineAlert(row))
	}
	return alerts, nil
}

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
