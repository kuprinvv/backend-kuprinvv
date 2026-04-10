package converter

import (
	"test-backend-1-kuprinvv/internal/handler/room/dto"
	"test-backend-1-kuprinvv/internal/model"
)

func CreateRoomRequestToService(body dto.CreateRoomRequest) model.Room {
	return model.Room{
		Name:        body.Name,
		Description: body.Description,
		Capacity:    body.Capacity,
	}
}

func ServiceToListRoomsResponse(rooms []model.Room) dto.ListRoomsResponse {
	resp := dto.ListRoomsResponse{Rooms: make([]dto.Room, 0, len(rooms))}
	for _, room := range rooms {
		resp.Rooms = append(resp.Rooms, roomServiceToRoomDto(room))
	}

	return resp
}

func ServiceToCreateRoomResponse(room model.Room) dto.CreateRoomResponse {
	return dto.CreateRoomResponse{
		Room: roomServiceToRoomDto(room),
	}
}

func roomServiceToRoomDto(room model.Room) dto.Room {
	return dto.Room{
		ID:          room.ID,
		Name:        room.Name,
		Description: room.Description,
		Capacity:    room.Capacity,
		CreatedAt:   room.CreatedAt,
	}
}
