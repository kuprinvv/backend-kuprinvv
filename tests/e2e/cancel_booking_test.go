package e2e

import (
	"fmt"
	"net/http"
	"testing"
)

// создать бронь → отменить
func TestCancelBookingFlow(t *testing.T) {
	adminToken := getToken(t, "admin")
	userToken := getToken(t, "user")

	slotID := setupRoomWithSlot(t, adminToken)
	bookingID := createBooking(t, slotID, userToken)

	resp := doPost(t, fmt.Sprintf("/bookings/%s/cancel", bookingID), nil, userToken)
	defer resp.Body.Close()
	assertStatus(t, resp, http.StatusOK)

	var out bookingResponse
	mustDecode(t, resp, &out)

	if out.Booking.Status != "cancelled" {
		t.Errorf("want status 'cancelled', got %q", out.Booking.Status)
	}
}

// Повторная отмена
func TestCancelBooking_Idempotent(t *testing.T) {
	adminToken := getToken(t, "admin")
	userToken := getToken(t, "user")

	slotID := setupRoomWithSlot(t, adminToken)
	bookingID := createBooking(t, slotID, userToken)

	resp := doPost(t, fmt.Sprintf("/bookings/%s/cancel", bookingID), nil, userToken)
	defer resp.Body.Close()
	assertStatus(t, resp, http.StatusOK)

	resp2 := doPost(t, fmt.Sprintf("/bookings/%s/cancel", bookingID), nil, userToken)
	defer resp2.Body.Close()
	assertStatus(t, resp2, http.StatusOK)

	var out bookingResponse
	mustDecode(t, resp2, &out)

	if out.Booking.Status != "cancelled" {
		t.Errorf("want status 'cancelled', got %q", out.Booking.Status)
	}
}

// После отмены слот снова свободен — можно создать новую бронь
func TestCancelBooking_SlotFreedAfterCancel(t *testing.T) {
	adminToken := getToken(t, "admin")
	userToken := getToken(t, "user")

	slotID := setupRoomWithSlot(t, adminToken)
	bookingID := createBooking(t, slotID, userToken)

	resp := doPost(t, fmt.Sprintf("/bookings/%s/cancel", bookingID), nil, userToken)
	defer resp.Body.Close()
	assertStatus(t, resp, http.StatusOK)

	// Слот снова свободен — новая бронь должна пройти
	createBooking(t, slotID, userToken)
}

// Отмена чужой брони
func TestCancelBooking_OtherUserForbidden(t *testing.T) {
	adminToken := getToken(t, "admin")
	userToken := getToken(t, "user")

	slotID := setupRoomWithSlot(t, adminToken)
	bookingID := createBooking(t, slotID, userToken)

	resp := doPost(t, fmt.Sprintf("/bookings/%s/cancel", bookingID), nil, adminToken)
	defer resp.Body.Close()
	assertStatus(t, resp, http.StatusForbidden)
}
