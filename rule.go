package engine

import (
	"encoding/json"
	"fmt"
	"time"
)

type SensorType string

const (
	Temperature SensorType = "temperature"
	Humidity    SensorType = "humidity"
	Motion      SensorType = "motion"
	CO2         SensorType = "co2"
)

type Operator string

const (
	OperatorGT  Operator = "gt"
	OperatorLT  Operator = "lt"
	OperatorGTE Operator = "gte"
	OperatorLTE Operator = "lte"
	OperatorEQ  Operator = "eq"
)

type ActionType string

const (
	ActionAlert ActionType = "alert"
)

type Condition struct {
	SensorType SensorType
	Operator   Operator
	Threshold  float64
}

type Action struct {
	Type    ActionType
	Payload string
}

// AlertPayload is the JSON shape carried by an Action of type ActionAlert.
type AlertPayload struct {
	Message  string        `json:"message"`
	Severity AlertSeverity `json:"severity"`
}

// validatePayload checks that the action's Payload is well-formed for its Type.
// It returns an error wrapping ErrInvalid so a bad payload is rejected at create
// time (400) instead of silently failing later, at evaluation.
func (a Action) validatePayload() error {
	switch a.Type {
	case ActionAlert:
		var p AlertPayload
		if err := json.Unmarshal([]byte(a.Payload), &p); err != nil {
			return fmt.Errorf("%w: action payload is not valid JSON: %v", ErrInvalid, err)
		}
		if p.Message == "" {
			return fmt.Errorf("%w: action payload message is required", ErrInvalid)
		}
		if !p.Severity.Valid() {
			return fmt.Errorf("%w: action payload severity %q", ErrInvalid, p.Severity)
		}
	}
	return nil
}

type Rule struct {
	ID        string
	Name      string
	Room      string
	Condition Condition
	Action    Action
	Enabled   bool
	CreatedAt time.Time
}

// Valid reports whether the sensor type is one the engine understands.
func (s SensorType) Valid() bool {
	switch s {
	case Temperature, Humidity, Motion, CO2:
		return true
	default:
		return false
	}
}

// Valid reports whether the operator is a known comparison.
func (o Operator) Valid() bool {
	switch o {
	case OperatorGT, OperatorLT, OperatorGTE, OperatorLTE, OperatorEQ:
		return true
	default:
		return false
	}
}

// Valid reports whether the action type has a handler.
func (a ActionType) Valid() bool {
	switch a {
	case ActionAlert:
		return true
	default:
		return false
	}
}

// Validate checks that a rule is well-formed before it is persisted. It returns
// an error wrapping ErrInvalid so transport layers can map it to a 400.
func (r Rule) Validate() error {
	if r.Name == "" {
		return fmt.Errorf("%w: name is required", ErrInvalid)
	}
	if r.Room == "" {
		return fmt.Errorf("%w: room is required", ErrInvalid)
	}
	if !r.Condition.SensorType.Valid() {
		return fmt.Errorf("%w: sensor_type %q", ErrInvalid, r.Condition.SensorType)
	}
	if !r.Condition.Operator.Valid() {
		return fmt.Errorf("%w: operator %q", ErrInvalid, r.Condition.Operator)
	}
	if !r.Action.Type.Valid() {
		return fmt.Errorf("%w: action type %q", ErrInvalid, r.Action.Type)
	}
	if err := r.Action.validatePayload(); err != nil {
		return err
	}
	return nil
}
