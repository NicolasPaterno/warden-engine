package http

import (
	"time"

	engine "github.com/NicolasPaterno/warden-engine"
)

// --- shared sub-objects ---

type conditionDTO struct {
	SensorType string  `json:"sensor_type"`
	Operator   string  `json:"operator"`
	Threshold  float64 `json:"threshold"`
}

type actionDTO struct {
	Type    string `json:"type"`
	Payload string `json:"payload"`
}

// --- rules ---

// ruleRequest is the body accepted on POST/PUT. It deliberately omits id and
// created_at: those are owned by the server, not the client.
type ruleRequest struct {
	Name      string       `json:"name"`
	Room      string       `json:"room"`
	Condition conditionDTO `json:"condition"`
	Action    actionDTO    `json:"action"`
	Enabled   bool         `json:"enabled"`
}

func (req ruleRequest) toEngine() engine.Rule {
	return engine.Rule{
		Name: req.Name,
		Room: req.Room,
		Condition: engine.Condition{
			SensorType: engine.SensorType(req.Condition.SensorType),
			Operator:   engine.Operator(req.Condition.Operator),
			Threshold:  req.Condition.Threshold,
		},
		Action: engine.Action{
			Type:    engine.ActionType(req.Action.Type),
			Payload: req.Action.Payload,
		},
		Enabled: req.Enabled,
	}
}

type ruleResponse struct {
	ID        string       `json:"id"`
	Name      string       `json:"name"`
	Room      string       `json:"room"`
	Condition conditionDTO `json:"condition"`
	Action    actionDTO    `json:"action"`
	Enabled   bool         `json:"enabled"`
	CreatedAt time.Time    `json:"created_at"`
}

func newRuleResponse(r engine.Rule) ruleResponse {
	return ruleResponse{
		ID:   r.ID,
		Name: r.Name,
		Room: r.Room,
		Condition: conditionDTO{
			SensorType: string(r.Condition.SensorType),
			Operator:   string(r.Condition.Operator),
			Threshold:  r.Condition.Threshold,
		},
		Action: actionDTO{
			Type:    string(r.Action.Type),
			Payload: r.Action.Payload,
		},
		Enabled:   r.Enabled,
		CreatedAt: r.CreatedAt,
	}
}

func newRuleResponses(rules []engine.Rule) []ruleResponse {
	out := make([]ruleResponse, 0, len(rules))
	for _, r := range rules {
		out = append(out, newRuleResponse(r))
	}
	return out
}

// --- alerts (read-only over REST) ---

type alertResponse struct {
	ID        string    `json:"id"`
	RuleID    string    `json:"rule_id"`
	Room      string    `json:"room"`
	Message   string    `json:"message"`
	Severity  string    `json:"severity"`
	Value     float64   `json:"value"`
	CreatedAt time.Time `json:"created_at"`
}

func newAlertResponse(a engine.Alert) alertResponse {
	return alertResponse{
		ID:        a.ID,
		RuleID:    a.RuleID,
		Room:      a.Room,
		Message:   a.Message,
		Severity:  string(a.Severity),
		Value:     a.Value,
		CreatedAt: a.CreatedAt,
	}
}

func newAlertResponses(alerts []engine.Alert) []alertResponse {
	out := make([]alertResponse, 0, len(alerts))
	for _, a := range alerts {
		out = append(out, newAlertResponse(a))
	}
	return out
}
