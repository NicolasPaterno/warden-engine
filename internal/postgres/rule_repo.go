package postgres

import (
	"context"

	engine "github.com/NicolasPaterno/warden-engine"
	db "github.com/NicolasPaterno/warden-engine/db/generated"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RuleRepo struct {
	queries *db.Queries
}

func NewRuleRepo(pool *pgxpool.Pool) *RuleRepo {
	return &RuleRepo{
		queries: db.New(pool),
	}
}

func (r *RuleRepo) Save(ctx context.Context, rule engine.Rule) (engine.Rule, error) {
	created, err := r.queries.CreateRule(ctx, db.CreateRuleParams{
		ID:         rule.ID,
		Name:       rule.Name,
		Room:       rule.Room,
		SensorType: toDBSensorType(rule.Condition.SensorType),
		Operator:   toDBOperator(rule.Condition.Operator),
		Threshold:  rule.Condition.Threshold,
		ActionType: toDBActionType(rule.Action.Type),
		Payload:    rule.Action.Payload,
		Enabled:    rule.Enabled,
		CreatedAt:  pgtype.Timestamptz{Time: rule.CreatedAt, Valid: true},
	})
	if err != nil {
		return engine.Rule{}, err
	}
	return toEngineRule(created), nil
}
func (r *RuleRepo) GetByID(ctx context.Context, id string) (engine.Rule, error)
func (r *RuleRepo) GetAll(ctx context.Context) ([]engine.Rule, error)
func (r *RuleRepo) GetEnabled(ctx context.Context) ([]engine.Rule, error)
func (r *RuleRepo) Update(ctx context.Context, rule engine.Rule) (engine.Rule, error)
func (r *RuleRepo) Delete(ctx context.Context, id string) error

func toEngineRule(r db.Rule) engine.Rule {
	return engine.Rule{
		ID:   r.ID,
		Name: r.Name,
		Room: r.Room,
		Condition: engine.Condition{
			SensorType: engine.SensorType(r.SensorType),
			Operator:   engine.Operator(r.Operator),
			Threshold:  r.Threshold,
		},
		Action: engine.Action{
			Type:    engine.ActionType(r.ActionType),
			Payload: r.Payload,
		},
		Enabled:   r.Enabled,
		CreatedAt: r.CreatedAt.Time,
	}
}
func toDBOperator(o engine.Operator) db.Operator {
	return db.Operator(o)
}
func toDBSensorType(s engine.SensorType) db.SensorType {
	return db.SensorType(s)
}
func toDBActionType(a engine.ActionType) db.ActionType {
	return db.ActionType(a)
}
