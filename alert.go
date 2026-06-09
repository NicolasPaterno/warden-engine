package engine

import "time"

type AlertSeverity string

const (
	SeverityInfo     AlertSeverity = "info"
	SeverityWarning  AlertSeverity = "warning"
	SeverityCritical AlertSeverity = "critical"
)

type Alert struct {
	ID        string
	RuleID    string
	Room      string
	Message   string
	Severity  AlertSeverity
	Value     float64
	CreatedAt time.Time
}

// Valid reports whether the severity is one the engine understands.
func (s AlertSeverity) Valid() bool {
	switch s {
	case SeverityInfo, SeverityWarning, SeverityCritical:
		return true
	default:
		return false
	}
}
