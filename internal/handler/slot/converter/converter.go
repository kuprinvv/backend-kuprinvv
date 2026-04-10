package converter

import (
	"test-backend-1-kuprinvv/internal/handler/slot/dto"
	"test-backend-1-kuprinvv/internal/model"
)

func ServiceToGetSlotsResponse(slots []model.Slot) dto.GetSlotsResponse {
	resp := dto.GetSlotsResponse{Slots: make([]dto.Slot, 0, len(slots))}
	for _, slot := range slots {
		resp.Slots = append(resp.Slots, slotServiceToSlotDto(slot))
	}

	return resp
}

func slotServiceToSlotDto(slots model.Slot) dto.Slot {
	return dto.Slot{
		ID:     slots.ID,
		RoomID: slots.RoomID,
		Start:  slots.StartTime,
		End:    slots.EndTime,
	}
}
