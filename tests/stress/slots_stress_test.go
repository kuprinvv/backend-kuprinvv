package stress

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

var baseURL = func() string {
	if v := os.Getenv("BASE_URL"); v != "" {
		return v
	}
	return "http://localhost:8080"
}()

// TestSlotsUnderLoad проверяет наиболее нагруженный эндпоинт в трёх фазах:
//
//  1. Cold — 10 комнат, первые запросы идут конкурентно (гонка при генерации слотов).
//  2. Ramp-up — ступенчатый рост нагрузки 100 → 300 → 500 → 1000 RPS (по 15 с).
//  3. Sustained — 500 RPS в течение 30 с, проверка SLI (p95 ≤ 200 мс, success ≥ 99.9%).
func TestSlotsUnderLoad(t *testing.T) {
	if testing.Short() {
		t.Skip("stress test skipped in short mode")
	}

	adminToken := mustGetToken(t, "admin")
	userToken := mustGetToken(t, "user")

	authHeader := http.Header{
		"Authorization": []string{"Bearer " + userToken},
	}

	const roomCount = 10
	date := nextWeekday()

	roomURLs := make([]string, roomCount)
	for i := range roomURLs {
		roomID := mustCreateRoom(t, adminToken)
		mustCreateSchedule(t, roomID, adminToken)
		roomURLs[i] = fmt.Sprintf("%s/rooms/%s/slots/list?date=%s", baseURL, roomID, date)
	}
	t.Logf("created %d rooms with schedules", roomCount)

	t.Log("=== Phase 1: Cold concurrent generation ===")

	const coldConcurrency = 10
	var (
		coldTotal  int64
		coldErrors int64
		coldMu     sync.Mutex
		statusMap  = make(map[int]int)
		wg         sync.WaitGroup
	)

	for _, u := range roomURLs {
		for range coldConcurrency {
			wg.Add(1)
			go func(url string) {
				defer wg.Done()
				req, _ := http.NewRequest(http.MethodGet, url, nil)
				req.Header.Set("Authorization", "Bearer "+userToken)

				resp, err := http.DefaultClient.Do(req)

				coldMu.Lock()
				coldTotal++
				if err != nil {
					coldErrors++
					statusMap[-1]++
				} else {
					statusMap[resp.StatusCode]++
					if resp.StatusCode != http.StatusOK {
						coldErrors++
					}
					resp.Body.Close()
				}
				coldMu.Unlock()
			}(u)
		}
	}
	wg.Wait()

	coldSuccess := float64(coldTotal-coldErrors) / float64(coldTotal) * 100
	t.Logf("requests: %d | success: %.2f%% | errors: %d | status codes: %v",
		coldTotal, coldSuccess, coldErrors, statusMap)
	if coldErrors > 0 {
		t.Errorf("cold phase: %d errors out of %d requests (status codes: %v)", coldErrors, coldTotal, statusMap)
	}

	t.Log("=== Phase 2: Ramp-up (warm) ===")

	targets := make([]vegeta.Target, len(roomURLs))
	for i, u := range roomURLs {
		targets[i] = vegeta.Target{
			Method: http.MethodGet,
			URL:    u,
			Header: authHeader,
		}
	}
	targeter := vegeta.NewStaticTargeter(targets...)

	steps := []int{100, 300, 500, 1000}
	stepDuration := 15 * time.Second

	for _, rps := range steps {
		rate := vegeta.Rate{Freq: rps, Per: time.Second}
		attacker := vegeta.NewAttacker()
		var m vegeta.Metrics

		for res := range attacker.Attack(targeter, rate, stepDuration, fmt.Sprintf("ramp-%d", rps)) {
			m.Add(res)
		}
		m.Close()

		t.Logf("%4d RPS | requests: %5d | success: %.2f%% | p50: %s | p95: %s | p99: %s",
			rps, m.Requests, m.Success*100,
			m.Latencies.P50, m.Latencies.P95, m.Latencies.P99,
		)
	}

	t.Log("=== Phase 3: Sustained 500 RPS (30s) ===")

	rate := vegeta.Rate{Freq: 500, Per: time.Second}
	attacker := vegeta.NewAttacker()
	var metrics vegeta.Metrics

	for res := range attacker.Attack(targeter, rate, 30*time.Second, "sustained") {
		metrics.Add(res)
	}
	metrics.Close()

	t.Logf("requests:    %d", metrics.Requests)
	t.Logf("success:     %.2f%%", metrics.Success*100)
	t.Logf("latency p50: %s", metrics.Latencies.P50)
	t.Logf("latency p95: %s", metrics.Latencies.P95)
	t.Logf("latency p99: %s", metrics.Latencies.P99)
	t.Logf("latency max: %s", metrics.Latencies.Max)
	t.Logf("throughput:  %.2f req/s", metrics.Throughput)

	if metrics.Success < 0.999 {
		t.Errorf("success rate %.2f%% < 99.9%%", metrics.Success*100)
	}
	if metrics.Latencies.P95 > 200*time.Millisecond {
		t.Errorf("p95 latency %s exceeds 200ms SLI", metrics.Latencies.P95)
	}
}

func mustGetToken(t *testing.T, role string) string {
	t.Helper()
	body, _ := json.Marshal(map[string]any{"role": role})
	resp, err := http.Post(baseURL+"/dummyLogin", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("dummyLogin failed: %v", err)
	}
	defer resp.Body.Close()

	var out struct {
		Token string `json:"token"`
	}
	err = json.NewDecoder(resp.Body).Decode(&out)
	if err != nil {
		return ""
	}
	if out.Token == "" {
		t.Fatal("empty token from dummyLogin")
	}
	return out.Token
}

func mustCreateRoom(t *testing.T, adminToken string) string {
	t.Helper()
	body, _ := json.Marshal(map[string]any{
		"name":     fmt.Sprintf("StressRoom-%d", time.Now().UnixNano()),
		"capacity": 10,
	})
	req, _ := http.NewRequest(http.MethodPost, baseURL+"/rooms/create", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+adminToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("create room failed: %v", err)
	}
	defer resp.Body.Close()

	var out struct {
		Room struct {
			ID string `json:"id"`
		} `json:"room"`
	}
	err = json.NewDecoder(resp.Body).Decode(&out)
	if err != nil {
		return ""
	}
	if out.Room.ID == "" {
		t.Fatal("empty room ID")
	}
	return out.Room.ID
}

func mustCreateSchedule(t *testing.T, roomID, adminToken string) {
	t.Helper()
	body, _ := json.Marshal(map[string]any{
		"daysOfWeek": []int{1, 2, 3, 4, 5, 6, 7},
		"startTime":  "09:00",
		"endTime":    "18:00",
	})
	req, _ := http.NewRequest(http.MethodPost,
		fmt.Sprintf("%s/rooms/%s/schedule/create", baseURL, roomID),
		bytes.NewReader(body),
	)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+adminToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("create schedule failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("create schedule: expected 201, got %d", resp.StatusCode)
	}
}

func nextWeekday() string {
	d := time.Now().UTC().Add(24 * time.Hour)
	for d.Weekday() == time.Saturday || d.Weekday() == time.Sunday {
		d = d.Add(24 * time.Hour)
	}
	return d.Format("2006-01-02")
}
