package slot

import "test-backend-1-kuprinvv/internal/service"

type Handler struct {
	slotServ service.SlotService
}

func NewSlotHandler(slotServ service.SlotService) *Handler {
	return &Handler{
		slotServ: slotServ,
	}
}
