package model

import (
	"time"

	"github.com/google/uuid"
)

type Booking struct {
	ID             uuid.UUID
	SlotID         uuid.UUID
	UserID         uuid.UUID
	Status         string
	ConferenceLink *string
	CreatedAt      *time.Time
}
