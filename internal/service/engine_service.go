package service

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	engine "github.com/NicolasPaterno/warden-engine"
	sensorv1 "github.com/NicolasPaterno/warden-proto/gen/go/warden/sensor/v1"
	"github.com/google/uuid"
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
	rules, err := s.rules.GetEnabled(ctx)
	if err != nil {
		return err
	}
	for _, rule := range rules {
		if rule.Room != reading.Room {
			continue
		}
		if rule.Condition.SensorType != fromProtoSensorType(reading.Type) {
			continue
		}
		if evaluate(rule, reading.Value) {
			var payload struct {
				Message  string `json:"message"`
				Severity string `json:"severity"`
			}
			if err := json.Unmarshal([]byte(rule.Action.Payload), &payload); err != nil {
				slog.Warn("invalid rule payload", "rule", rule.Name, "error", err)
				continue
			}

			alert := engine.Alert{
				ID:        uuid.NewString(),
				RuleID:    rule.ID,
				Room:      rule.Room,
				Message:   payload.Message,
				Severity:  engine.AlertSeverity(payload.Severity),
				Value:     reading.Value,
				CreatedAt: time.Now(),
			}

			if err := s.alerts.Save(ctx, alert); err != nil {
				return err
			}

			slog.Info("rule fired", "rule", rule.Name, "room", rule.Room, "value", reading.Value)
		}
	}
	return nil
}

func evaluate(rule engine.Rule, value float64) bool {
	switch rule.Condition.Operator {
	case engine.OperatorGT:
		return value > rule.Condition.Threshold
	case engine.OperatorLT:
		return value < rule.Condition.Threshold
	case engine.OperatorGTE:
		return value >= rule.Condition.Threshold
	case engine.OperatorLTE:
		return value <= rule.Condition.Threshold
	case engine.OperatorEQ:
		return value == rule.Condition.Threshold
	default:
		return false
	}
}

func fromProtoSensorType(t sensorv1.SensorType) engine.SensorType {
	switch t {
	case sensorv1.SensorType_SENSOR_TYPE_TEMPERATURE:
		return engine.Temperature
	case sensorv1.SensorType_SENSOR_TYPE_HUMIDITY:
		return engine.Humidity
	case sensorv1.SensorType_SENSOR_TYPE_MOTION:
		return engine.Motion
	case sensorv1.SensorType_SENSOR_TYPE_CO2:
		return engine.CO2
	default:
		return ""
	}
}
