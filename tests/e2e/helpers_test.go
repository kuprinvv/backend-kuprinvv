package e2e

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"
)

var httpClient = &http.Client{Timeout: 10 * time.Second}

type bookingResponse struct {
	Booking struct {
		ID     string `json:"id"`
		Status string `json:"status"`
		SlotID string `json:"slotId"`
	} `json:"booking"`
}

func getToken(t *testing.T, role string) string {
	t.Helper()
	resp := doPost(t, "/dummyLogin", map[string]any{"role": role}, "")
	defer resp.Body.Close()
	var out struct {
		Token string `json:"token"`
	}
	mustDecode(t, resp, &out)
	return out.Token
}

func doPost(t *testing.T, path string, body any, token string) *http.Response {
	t.Helper()
	b, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal request body: %v", err)
	}
	req, err := http.NewRequest(http.MethodPost, baseURL+path, bytes.NewReader(b))
	if err != nil {
		t.Fatalf("create POST request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		t.Fatalf("POST %s: %v", path, err)
	}
	return resp
}

func doGet(t *testing.T, path string, token string) *http.Response {
	t.Helper()
	req, err := http.NewRequest(http.MethodGet, baseURL+path, nil)
	if err != nil {
		t.Fatalf("create GET request: %v", err)
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		t.Fatalf("GET %s: %v", path, err)
	}
	return resp
}

func mustDecode(t *testing.T, resp *http.Response, target any) {
	t.Helper()
	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		t.Fatalf("decode response: %v", err)
	}
}

func assertStatus(t *testing.T, resp *http.Response, expected int) {
	t.Helper()
	if resp.StatusCode != expected {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("expected HTTP %d, got %d, body: %s", expected, resp.StatusCode, string(body))
	}
}

func nextWeekday() string {
	d := time.Now().UTC().Add(24 * time.Hour)
	for d.Weekday() == time.Saturday || d.Weekday() == time.Sunday {
		d = d.Add(24 * time.Hour)
	}
	return d.Format("2006-01-02")
}

func setupRoom(t *testing.T, adminToken string) string {
	t.Helper()
	resp := doPost(t, "/rooms/create", map[string]any{
		"name":     fmt.Sprintf("TestRoom-%d", time.Now().UnixNano()),
		"capacity": 6,
	}, adminToken)
	defer resp.Body.Close()
	assertStatus(t, resp, http.StatusCreated)

	var out struct {
		Room struct {
			ID string `json:"id"`
		} `json:"room"`
	}
	mustDecode(t, resp, &out)

	if out.Room.ID == "" {
		t.Fatal("room ID is empty")
	}
	return out.Room.ID
}

func setupSchedule(t *testing.T, roomID string, adminToken string) {
	t.Helper()
	resp := doPost(t, fmt.Sprintf("/rooms/%s/schedule/create", roomID), map[string]any{
		"daysOfWeek": []int{1, 2, 3, 4, 5, 6, 7},
		"startTime":  "09:00",
		"endTime":    "18:00",
	}, adminToken)
	defer resp.Body.Close()
	assertStatus(t, resp, http.StatusCreated)
}

func setupRoomWithSlot(t *testing.T, adminToken string) string {
	t.Helper()
	roomID := setupRoom(t, adminToken)
	setupSchedule(t, roomID, adminToken)

	date := nextWeekday()
	resp := doGet(t, fmt.Sprintf("/rooms/%s/slots/list?date=%s", roomID, date), adminToken)
	defer resp.Body.Close()
	assertStatus(t, resp, http.StatusOK)

	var slots struct {
		Items []struct {
			ID string `json:"id"`
		} `json:"slots"`
	}
	mustDecode(t, resp, &slots)

	if len(slots.Items) == 0 {
		t.Fatal("no slots generated — check schedule or date logic")
	}

	return slots.Items[0].ID
}

func createBooking(t *testing.T, slotID string, userToken string) string {
	t.Helper()
	resp := doPost(t, "/bookings/create", map[string]any{
		"slotId": slotID,
	}, userToken)
	defer resp.Body.Close()
	assertStatus(t, resp, http.StatusCreated)

	var out bookingResponse
	mustDecode(t, resp, &out)

	if out.Booking.ID == "" {
		t.Fatal("booking ID must not be empty")
	}
	return out.Booking.ID
}
