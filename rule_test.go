package engine

import (
	"errors"
	"testing"
)

// validRule returns a rule that passes Validate. Each test mutates one field so
// a failure points at exactly the constraint under test.
func validRule() Rule {
	return Rule{
		Name: "bedroom too hot",
		Room: "bedroom",
		Condition: Condition{
			SensorType: Temperature,
			Operator:   OperatorGT,
			Threshold:  30,
		},
		Action: Action{
			Type:    ActionAlert,
			Payload: `{"message":"temp alta","severity":"warning"}`,
		},
		Enabled: true,
	}
}

func TestRuleValidate(t *testing.T) {
	tests := []struct {
		name    string
		mutate  func(*Rule)
		wantErr bool
	}{
		{"valid rule", func(*Rule) {}, false},
		{"empty name", func(r *Rule) { r.Name = "" }, true},
		{"empty room", func(r *Rule) { r.Room = "" }, true},
		{"unknown sensor type", func(r *Rule) { r.Condition.SensorType = "pressure" }, true},
		{"unknown operator", func(r *Rule) { r.Condition.Operator = "between" }, true},
		{"unknown action type", func(r *Rule) { r.Action.Type = "smoke_signal" }, true},
		{"malformed payload JSON", func(r *Rule) { r.Action.Payload = "not json" }, true},
		{"payload missing message", func(r *Rule) { r.Action.Payload = `{"severity":"warning"}` }, true},
		{"payload invalid severity", func(r *Rule) { r.Action.Payload = `{"message":"x","severity":"banana"}` }, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := validRule()
			tt.mutate(&rule)

			err := rule.Validate()

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected an error, got nil")
				}
				// Every validation failure must wrap ErrInvalid so the HTTP
				// layer can map it to a 400.
				if !errors.Is(err, ErrInvalid) {
					t.Fatalf("error %q does not wrap ErrInvalid", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
		})
	}
}

func TestSensorTypeValid(t *testing.T) {
	tests := []struct {
		in   SensorType
		want bool
	}{
		{Temperature, true},
		{Humidity, true},
		{Motion, true},
		{CO2, true},
		{"pressure", false},
		{"", false},
	}
	for _, tt := range tests {
		if got := tt.in.Valid(); got != tt.want {
			t.Errorf("SensorType(%q).Valid() = %v, want %v", tt.in, got, tt.want)
		}
	}
}

func TestOperatorValid(t *testing.T) {
	tests := []struct {
		in   Operator
		want bool
	}{
		{OperatorGT, true},
		{OperatorLT, true},
		{OperatorGTE, true},
		{OperatorLTE, true},
		{OperatorEQ, true},
		{"between", false},
		{"", false},
	}
	for _, tt := range tests {
		if got := tt.in.Valid(); got != tt.want {
			t.Errorf("Operator(%q).Valid() = %v, want %v", tt.in, got, tt.want)
		}
	}
}

func TestAlertSeverityValid(t *testing.T) {
	tests := []struct {
		in   AlertSeverity
		want bool
	}{
		{SeverityInfo, true},
		{SeverityWarning, true},
		{SeverityCritical, true},
		{"banana", false},
		{"", false},
	}
	for _, tt := range tests {
		if got := tt.in.Valid(); got != tt.want {
			t.Errorf("AlertSeverity(%q).Valid() = %v, want %v", tt.in, got, tt.want)
		}
	}
}

func TestActionTypeValid(t *testing.T) {
	tests := []struct {
		in   ActionType
		want bool
	}{
		{ActionAlert, true},
		{"notify_email", false},
		{"", false},
	}
	for _, tt := range tests {
		if got := tt.in.Valid(); got != tt.want {
			t.Errorf("ActionType(%q).Valid() = %v, want %v", tt.in, got, tt.want)
		}
	}
}
