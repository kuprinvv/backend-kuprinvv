package dto

import (
	"time"

	"github.com/google/uuid"
)

type ListRoomsResponse struct {
	Rooms []Room `json:"rooms"`
}

type CreateRoomRequest struct {
	Name        string  `json:"name" validate:"required"`
	Description *string `json:"description"`
	Capacity    *int    `json:"capacity" validate:"omitempty,min=1"`
}

type CreateRoomResponse struct {
	Room Room `json:"room"`
}

type Room struct {
	ID          uuid.UUID  `json:"id"`
	Name        string     `json:"name"`
	Description *string    `json:"description,omitempty"`
	Capacity    *int       `json:"capacity,omitempty"`
	CreatedAt   *time.Time `json:"createdAt,omitempty"`
}
