package service

import (
	"context"

	engine "github.com/NicolasPaterno/warden-engine"
	sensorv1 "github.com/NicolasPaterno/warden-proto/gen/go/warden/sensor/v1"
)

type EngineService struct {
	rules  engine.RuleRepository
	alerts engine.AlertRepository
}

func NewEngineService(rules engine.RuleRepository, alerts engine.AlertRepository) *EngineService {
	return &EngineService{
		rules:  rules,
		alerts: alerts,
	}
}

func (s *EngineService) Evaluate(ctx context.Context, reading *sensorv1.SensorReading) error {
	// 1. buscar só regras habilitadas: s.rules.GetEnabled(ctx)
	// 2. se erro, retornar
	// 3. para cada rule no slice:
	//      a. pular se rule.Room != reading.Room
	//      b. pular se string(rule.Condition.SensorType) != reading.Type.String()
	//         dica: reading.Type é um enum proto — use .String() e compare lowercase
	//      c. chamar evaluate(rule, reading.Value)
	//      d. se true:
	//           - criar engine.Alert{ID: uuid.New().String(), RuleID: rule.ID, ...}
	//           - salvar com s.alerts.Save(ctx, alert)
	//           - logar com slog.Info
	// 4. retornar nil
}

func evaluate(rule engine.Rule, value float64) bool {
	// switch em rule.Condition.Operator:
	//   OperatorGT  → value > rule.Condition.Threshold
	//   OperatorLT  → value < rule.Condition.Threshold
	//   OperatorGTE → value >= rule.Condition.Threshold
	//   OperatorLTE → value <= rule.Condition.Threshold
	//   OperatorEQ  → value == rule.Condition.Threshold
	//   default     → false
}
