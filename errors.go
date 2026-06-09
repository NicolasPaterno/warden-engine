package engine

import "errors"

// Domain-level sentinel errors. Transport layers map these to protocol codes
// (e.g. the HTTP layer turns ErrNotFound into 404, ErrInvalid into 400) without
// knowing anything about the underlying infrastructure.
var (
	// ErrNotFound is returned by repositories when a requested entity does not exist.
	ErrNotFound = errors.New("not found")
	// ErrInvalid is returned when an entity fails validation. Wrap it with detail:
	// fmt.Errorf("%w: operator %q", ErrInvalid, op).
	ErrInvalid = errors.New("invalid")
)
