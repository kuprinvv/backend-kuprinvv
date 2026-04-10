package converter_test

import (
	"test-backend-1-kuprinvv/internal/handler/room/converter"
	"test-backend-1-kuprinvv/internal/handler/room/dto"
	"test-backend-1-kuprinvv/internal/model"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateRoomRequestToService(t *testing.T) {
	desc := "большая комната"
	cap := 10
	body := dto.CreateRoomRequest{Name: "Room 1", Description: &desc, Capacity: &cap}

	room := converter.CreateRoomRequestToService(body)
	assert.Equal(t, "Room 1", room.Name)
	assert.Equal(t, &desc, room.Description)
	assert.Equal(t, &cap, room.Capacity)
}

func TestServiceToListRoomsResponse(t *testing.T) {
	r1 := model.Room{ID: uuid.New(), Name: "A"}
	r2 := model.Room{ID: uuid.New(), Name: "B"}

	t.Run("два элемента", func(t *testing.T) {
		resp := converter.ServiceToListRoomsResponse([]model.Room{r1, r2})
		require.Len(t, resp.Rooms, 2)
		assert.Equal(t, r1.ID, resp.Rooms[0].ID)
	})

	t.Run("пустой список", func(t *testing.T) {
		resp := converter.ServiceToListRoomsResponse([]model.Room{})
		assert.NotNil(t, resp.Rooms)
		assert.Empty(t, resp.Rooms)
	})
}

func TestServiceToCreateRoomResponse(t *testing.T) {
	r := model.Room{ID: uuid.New(), Name: "Room X"}
	resp := converter.ServiceToCreateRoomResponse(r)
	assert.Equal(t, r.ID, resp.Room.ID)
	assert.Equal(t, r.Name, resp.Room.Name)
}
