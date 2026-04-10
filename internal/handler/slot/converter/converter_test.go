package converter_test

import (
	"test-backend-1-kuprinvv/internal/handler/slot/converter"
	"test-backend-1-kuprinvv/internal/model"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServiceToGetSlotsResponse(t *testing.T) {
	now := time.Now()
	s1 := model.Slot{ID: uuid.New(), RoomID: uuid.New(), StartTime: now, EndTime: now.Add(30 * time.Minute)}
	s2 := model.Slot{ID: uuid.New(), RoomID: uuid.New(), StartTime: now.Add(time.Hour), EndTime: now.Add(90 * time.Minute)}

	t.Run("два слота", func(t *testing.T) {
		resp := converter.ServiceToGetSlotsResponse([]model.Slot{s1, s2})
		require.Len(t, resp.Slots, 2)
		assert.Equal(t, s1.ID, resp.Slots[0].ID)
		assert.Equal(t, s1.RoomID, resp.Slots[0].RoomID)
		assert.Equal(t, s1.StartTime, resp.Slots[0].Start)
		assert.Equal(t, s1.EndTime, resp.Slots[0].End)
	})

	t.Run("пустой список", func(t *testing.T) {
		resp := converter.ServiceToGetSlotsResponse([]model.Slot{})
		assert.NotNil(t, resp.Slots)
		assert.Empty(t, resp.Slots)
	})
}
