package engine

import (
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
	return nil
}
