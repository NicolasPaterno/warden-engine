package postgres

import (
	"context"

	engine "github.com/NicolasPaterno/warden-engine"
	db "github.com/NicolasPaterno/warden-engine/db/generated"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AlertRepo struct {
	queries *db.Queries
}

func NewAlertRepo(pool *pgxpool.Pool) *AlertRepo

func (r *AlertRepo) Save(ctx context.Context, alert engine.Alert) error
func (r *AlertRepo) GetByRoom(ctx context.Context, room string) ([]engine.Alert, error)
func (r *AlertRepo) GetAll(ctx context.Context) ([]engine.Alert, error)

func toEngineAlert(a db.Alert) engine.Alert
func toDBAlertSeverity(s engine.AlertSeverity) db.AlertSeverity
