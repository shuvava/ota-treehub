package data

import "github.com/google/uuid"

// CorrelationID wrapper on the top of github.com/google/uuid
type CorrelationID uuid.UUID

func (c CorrelationID) String() string {
	return uuid.UUID(c).String()
}

// NewCorrelationID creates a new CorrelationID
func NewCorrelationID() CorrelationID {
	id := uuid.New()
	return CorrelationID(id)
}
