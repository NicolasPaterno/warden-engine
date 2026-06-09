package service

import (
	"context"
	"errors"
	"testing"

	engine "github.com/NicolasPaterno/warden-engine"
	sensorv1 "github.com/NicolasPaterno/warden-proto/gen/go/warden/sensor/v1"
)

// --- fakes ---------------------------------------------------------------

// fakeRuleRepo is a hand-rolled stub of engine.RuleRepository. Evaluate only
// touches GetEnabled, so the other methods exist solely to satisfy the
// interface and panic if a future change starts calling them unexpectedly.
type fakeRuleRepo struct {
	enabled    []engine.Rule
	enabledErr error
}

func (f *fakeRuleRepo) GetEnabled(context.Context) ([]engine.Rule, error) {
	return f.enabled, f.enabledErr
}
func (f *fakeRuleRepo) Save(context.Context, engine.Rule) error { panic("unexpected Save") }
func (f *fakeRuleRepo) GetByID(context.Context, string) (engine.Rule, error) {
	panic("unexpected GetByID")
}
func (f *fakeRuleRepo) GetAll(context.Context) ([]engine.Rule, error) { panic("unexpected GetAll") }
func (f *fakeRuleRepo) Update(context.Context, engine.Rule) (engine.Rule, error) {
	panic("unexpected Update")
}
func (f *fakeRuleRepo) Delete(context.Context, string) error { panic("unexpected Delete") }

// fakeAlertRepo records every alert handed to Save and can be told to fail.
type fakeAlertRepo struct {
	saved   []engine.Alert
	saveErr error
}

func (f *fakeAlertRepo) Save(_ context.Context, a engine.Alert) error {
	if f.saveErr != nil {
		return f.saveErr
	}
	f.saved = append(f.saved, a)
	return nil
}
func (f *fakeAlertRepo) GetByRoom(context.Context, string) ([]engine.Alert, error) {
	panic("unexpected GetByRoom")
}
func (f *fakeAlertRepo) GetByRule(context.Context, string) ([]engine.Alert, error) {
	panic("unexpected GetByRule")
}
func (f *fakeAlertRepo) GetAll(context.Context) ([]engine.Alert, error) { panic("unexpected GetAll") }

// tempRule is a rule that fires on bedroom temperature above 30°C.
func tempRule() engine.Rule {
	return engine.Rule{
		ID:   "rule-1",
		Name: "bedroom hot",
		Room: "bedroom",
		Condition: engine.Condition{
			SensorType: engine.Temperature,
			Operator:   engine.OperatorGT,
			Threshold:  30,
		},
		Action: engine.Action{
			Type:    engine.ActionAlert,
			Payload: `{"message":"temp alta","severity":"warning"}`,
		},
		Enabled: true,
	}
}

func reading(room string, t sensorv1.SensorType, value float64) *sensorv1.SensorReading {
	return &sensorv1.SensorReading{Room: room, Type: t, Value: value}
}

// --- evaluate (pure operator logic) -------------------------------------

func TestEvaluate(t *testing.T) {
	cond := func(op engine.Operator, threshold float64) engine.Rule {
		return engine.Rule{Condition: engine.Condition{Operator: op, Threshold: threshold}}
	}
	tests := []struct {
		name  string
		rule  engine.Rule
		value float64
		want  bool
	}{
		{"gt above", cond(engine.OperatorGT, 30), 31, true},
		{"gt equal", cond(engine.OperatorGT, 30), 30, false},
		{"lt below", cond(engine.OperatorLT, 30), 29, true},
		{"lt equal", cond(engine.OperatorLT, 30), 30, false},
		{"gte equal", cond(engine.OperatorGTE, 30), 30, true},
		{"gte below", cond(engine.OperatorGTE, 30), 29, false},
		{"lte equal", cond(engine.OperatorLTE, 30), 30, true},
		{"lte above", cond(engine.OperatorLTE, 30), 31, false},
		{"eq match", cond(engine.OperatorEQ, 30), 30, true},
		{"eq miss", cond(engine.OperatorEQ, 30), 30.1, false},
		{"unknown operator", cond("between", 30), 30, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := evaluate(tt.rule, tt.value); got != tt.want {
				t.Errorf("evaluate(%v, %v) = %v, want %v", tt.rule.Condition, tt.value, got, tt.want)
			}
		})
	}
}

func TestFromProtoSensorType(t *testing.T) {
	tests := []struct {
		in   sensorv1.SensorType
		want engine.SensorType
	}{
		{sensorv1.SensorType_SENSOR_TYPE_TEMPERATURE, engine.Temperature},
		{sensorv1.SensorType_SENSOR_TYPE_HUMIDITY, engine.Humidity},
		{sensorv1.SensorType_SENSOR_TYPE_MOTION, engine.Motion},
		{sensorv1.SensorType_SENSOR_TYPE_CO2, engine.CO2},
		{sensorv1.SensorType_SENSOR_TYPE_UNSPECIFIED, ""},
	}
	for _, tt := range tests {
		if got := fromProtoSensorType(tt.in); got != tt.want {
			t.Errorf("fromProtoSensorType(%v) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

// --- Evaluate (end-to-end over the fakes) -------------------------------

func TestEngineServiceEvaluate(t *testing.T) {
	tempReading := reading("bedroom", sensorv1.SensorType_SENSOR_TYPE_TEMPERATURE, 31)

	t.Run("fires and persists an alert when everything matches", func(t *testing.T) {
		alerts := &fakeAlertRepo{}
		svc := NewEngineService(&fakeRuleRepo{enabled: []engine.Rule{tempRule()}}, alerts)

		if err := svc.Evaluate(context.Background(), tempReading); err != nil {
			t.Fatalf("Evaluate returned error: %v", err)
		}
		if len(alerts.saved) != 1 {
			t.Fatalf("expected 1 alert saved, got %d", len(alerts.saved))
		}
		got := alerts.saved[0]
		if got.RuleID != "rule-1" {
			t.Errorf("alert.RuleID = %q, want rule-1", got.RuleID)
		}
		if got.Room != "bedroom" {
			t.Errorf("alert.Room = %q, want bedroom", got.Room)
		}
		if got.Message != "temp alta" {
			t.Errorf("alert.Message = %q, want 'temp alta'", got.Message)
		}
		if got.Severity != engine.SeverityWarning {
			t.Errorf("alert.Severity = %q, want warning", got.Severity)
		}
		if got.Value != 31 {
			t.Errorf("alert.Value = %v, want 31 (the reading, not the threshold)", got.Value)
		}
		if got.ID == "" {
			t.Error("alert.ID is empty, expected a generated UUID")
		}
		if got.CreatedAt.IsZero() {
			t.Error("alert.CreatedAt is zero, expected it to be set")
		}
	})

	noFireCases := []struct {
		name    string
		reading *sensorv1.SensorReading
	}{
		{"room mismatch", reading("kitchen", sensorv1.SensorType_SENSOR_TYPE_TEMPERATURE, 31)},
		{"sensor type mismatch", reading("bedroom", sensorv1.SensorType_SENSOR_TYPE_HUMIDITY, 31)},
		{"condition false", reading("bedroom", sensorv1.SensorType_SENSOR_TYPE_TEMPERATURE, 25)},
	}
	for _, tc := range noFireCases {
		t.Run("no alert on "+tc.name, func(t *testing.T) {
			alerts := &fakeAlertRepo{}
			svc := NewEngineService(&fakeRuleRepo{enabled: []engine.Rule{tempRule()}}, alerts)

			if err := svc.Evaluate(context.Background(), tc.reading); err != nil {
				t.Fatalf("Evaluate returned error: %v", err)
			}
			if len(alerts.saved) != 0 {
				t.Fatalf("expected no alert, got %d", len(alerts.saved))
			}
		})
	}

	t.Run("skips rule with malformed payload without erroring", func(t *testing.T) {
		bad := tempRule()
		bad.Action.Payload = "not json"
		alerts := &fakeAlertRepo{}
		svc := NewEngineService(&fakeRuleRepo{enabled: []engine.Rule{bad}}, alerts)

		if err := svc.Evaluate(context.Background(), tempReading); err != nil {
			t.Fatalf("Evaluate returned error: %v", err)
		}
		if len(alerts.saved) != 0 {
			t.Fatalf("expected no alert, got %d", len(alerts.saved))
		}
	})

	t.Run("propagates GetEnabled error", func(t *testing.T) {
		wantErr := errors.New("db down")
		svc := NewEngineService(&fakeRuleRepo{enabledErr: wantErr}, &fakeAlertRepo{})

		if err := svc.Evaluate(context.Background(), tempReading); !errors.Is(err, wantErr) {
			t.Fatalf("Evaluate error = %v, want %v", err, wantErr)
		}
	})

	t.Run("propagates alert Save error", func(t *testing.T) {
		wantErr := errors.New("insert failed")
		svc := NewEngineService(
			&fakeRuleRepo{enabled: []engine.Rule{tempRule()}},
			&fakeAlertRepo{saveErr: wantErr},
		)

		if err := svc.Evaluate(context.Background(), tempReading); !errors.Is(err, wantErr) {
			t.Fatalf("Evaluate error = %v, want %v", err, wantErr)
		}
	})
}
