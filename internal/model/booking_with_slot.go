package model

import (
	"time"

	"github.com/google/uuid"
)

type BookingWithSlot struct {
	ID        uuid.UUID
	Status    string
	StartTime time.Time
	EndTime   time.Time
}
