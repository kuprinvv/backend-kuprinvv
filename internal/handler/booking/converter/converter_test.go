package converter_test

import (
	"test-backend-1-kuprinvv/internal/handler/booking/converter"
	"test-backend-1-kuprinvv/internal/model"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeBooking() model.Booking {
	now := time.Now()
	link := "https://meet.example.com/abc"
	return model.Booking{
		ID:             uuid.New(),
		SlotID:         uuid.New(),
		UserID:         uuid.New(),
		Status:         "active",
		ConferenceLink: &link,
		CreatedAt:      &now,
	}
}

func TestServiceToBookingResponse(t *testing.T) {
	b := makeBooking()
	resp := converter.ServiceToBookingResponse(b)
	assert.Equal(t, b.ID, resp.Booking.ID)
	assert.Equal(t, b.Status, resp.Booking.Status)
	assert.Equal(t, b.ConferenceLink, resp.Booking.ConferenceLink)
}

func TestServiceToBookingsResponse(t *testing.T) {
	b1, b2 := makeBooking(), makeBooking()

	t.Run("два элемента", func(t *testing.T) {
		resp := converter.ServiceToBookingsResponse([]model.Booking{b1, b2})
		require.Len(t, resp.Bookings, 2)
		assert.Equal(t, b1.ID, resp.Bookings[0].ID)
		assert.Equal(t, b2.ID, resp.Bookings[1].ID)
	})

	t.Run("пустой список", func(t *testing.T) {
		resp := converter.ServiceToBookingsResponse([]model.Booking{})
		assert.NotNil(t, resp.Bookings)
		assert.Empty(t, resp.Bookings)
	})
}

func TestServiceToBookingsListResponse(t *testing.T) {
	b := makeBooking()
	pagination := model.Pagination{Page: 2, PageSize: 10, Total: 55}

	resp := converter.ServiceToBookingsListResponse([]model.Booking{b}, pagination)
	require.Len(t, resp.Bookings, 1)
	assert.Equal(t, b.ID, resp.Bookings[0].ID)
	assert.Equal(t, 2, resp.Pagination.Page)
	assert.Equal(t, 10, resp.Pagination.PageSize)
	assert.Equal(t, 55, resp.Pagination.Total)
}
