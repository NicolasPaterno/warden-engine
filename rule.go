package engine

import "time"

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
