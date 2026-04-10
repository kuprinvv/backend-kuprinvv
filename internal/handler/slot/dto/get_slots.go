package dto

import (
	"time"

	"github.com/google/uuid"
)

type GetSlotsResponse struct {
	Slots []Slot `json:"slots"`
}

type Slot struct {
	ID     uuid.UUID `json:"id"`
	RoomID uuid.UUID `json:"roomId"`
	Start  time.Time `json:"start"`
	End    time.Time `json:"end"`
}
