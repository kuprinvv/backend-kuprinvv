package e2e

import (
	"net/http"
	"testing"
)

// Переговорка → расписание → бронь
func TestCreateBookingFlow(t *testing.T) {
	adminToken := getToken(t, "admin")
	userToken := getToken(t, "user")

	slotID := setupRoomWithSlot(t, adminToken)

	resp := doPost(t, "/bookings/create", map[string]any{
		"slotId": slotID,
	}, userToken)
	defer resp.Body.Close()
	assertStatus(t, resp, http.StatusCreated)

	var out bookingResponse
	mustDecode(t, resp, &out)

	if out.Booking.Status != "active" {
		t.Errorf("want status 'active', got %q", out.Booking.Status)
	}
	if out.Booking.ID == "" {
		t.Error("booking ID must not be empty")
	}
	if out.Booking.SlotID != slotID {
		t.Errorf("want slotId %q, got %q", slotID, out.Booking.SlotID)
	}
}

// Двойное бронирование одного слота
func TestCreateBooking_DoubleBooking(t *testing.T) {
	adminToken := getToken(t, "admin")
	userToken := getToken(t, "user")

	slotID := setupRoomWithSlot(t, adminToken)
	createBooking(t, slotID, userToken)

	resp := doPost(t, "/bookings/create", map[string]any{
		"slotId": slotID,
	}, userToken)
	defer resp.Body.Close()
	assertStatus(t, resp, http.StatusConflict)
}

// Admin не может создавать брони
func TestCreateBooking_AdminForbidden(t *testing.T) {
	adminToken := getToken(t, "admin")

	slotID := setupRoomWithSlot(t, adminToken)

	resp := doPost(t, "/bookings/create", map[string]any{
		"slotId": slotID,
	}, adminToken)
	defer resp.Body.Close()
	assertStatus(t, resp, http.StatusForbidden)
}
