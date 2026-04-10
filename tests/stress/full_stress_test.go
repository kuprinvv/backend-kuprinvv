package stress

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

func TestFullLoadScenario(t *testing.T) {
	if testing.Short() {
		t.Skip("stress test skipped in short mode")
	}

	adminToken := mustGetToken(t, "admin")
	userToken := mustGetToken(t, "user")

	roomIDs := make([]string, 3)
	for i := range roomIDs {
		roomIDs[i] = mustCreateRoom(t, adminToken)
		mustCreateSchedule(t, roomIDs[i], adminToken)
		t.Logf("created room %d: %s", i+1, roomIDs[i])
	}

	date := nextWeekday()

	authHeader := http.Header{
		"Authorization": []string{"Bearer " + userToken},
		"Content-Type":  []string{"application/json"},
	}

	for _, roomID := range roomIDs {
		url := fmt.Sprintf("%s/rooms/%s/slots/list?date=%s", baseURL, roomID, date)
		req, _ := http.NewRequest(http.MethodGet, url, nil)
		req.Header.Set("Authorization", "Bearer "+userToken)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return
		}
		resp.Body.Close()
	}
	t.Log("warmup done")

	type namedAttack struct {
		name    string
		weight  int
		targets []vegeta.Target
	}

	attacks := []namedAttack{
		{
			name:   "получение слотов",
			weight: 70,
			targets: func() []vegeta.Target {
				var tt []vegeta.Target
				for _, roomID := range roomIDs {
					tt = append(tt, vegeta.Target{
						Method: http.MethodGet,
						URL:    fmt.Sprintf("%s/rooms/%s/slots/list?date=%s", baseURL, roomID, date),
						Header: authHeader,
					})
				}
				return tt
			}(),
		},
		{
			name:   "список переговорок",
			weight: 15,
			targets: []vegeta.Target{
				{
					Method: http.MethodGet,
					URL:    baseURL + "/rooms/list",
					Header: authHeader,
				},
			},
		},
		{
			name:   "мои брони",
			weight: 10,
			targets: []vegeta.Target{
				{
					Method: http.MethodGet,
					URL:    baseURL + "/bookings/my",
					Header: authHeader,
				},
			},
		},
		{
			name:   "создание брони",
			weight: 5,
			targets: func() []vegeta.Target {
				url := fmt.Sprintf("%s/rooms/%s/slots/list?date=%s", baseURL, roomIDs[0], date)
				req, _ := http.NewRequest(http.MethodGet, url, nil)
				req.Header.Set("Authorization", "Bearer "+userToken)
				resp, err := http.DefaultClient.Do(req)
				if err != nil {
					return nil
				}
				defer resp.Body.Close()

				var out struct {
					Slots []struct {
						ID string `json:"id"`
					} `json:"slots"`
				}
				err = json.NewDecoder(resp.Body).Decode(&out)
				if err != nil {
					return nil
				}

				var tt []vegeta.Target
				for i, s := range out.Slots {
					if i >= 5 {
						break
					}
					body, _ := json.Marshal(map[string]any{"slotId": s.ID})
					tt = append(tt, vegeta.Target{
						Method: http.MethodPost,
						URL:    baseURL + "/bookings/create",
						Header: authHeader,
						Body:   body,
					})
				}
				return tt
			}(),
		},
	}

	duration := 30 * time.Second
	t.Log("=== Results by endpoint ===")

	allPassed := true

	for _, attack := range attacks {
		if len(attack.targets) == 0 {
			t.Logf("[%s] skipped — no targets", attack.name)
			continue
		}

		attackRate := vegeta.Rate{
			Freq: attack.weight,
			Per:  time.Second,
		}

		targeter := vegeta.NewStaticTargeter(attack.targets...)
		attacker := vegeta.NewAttacker()
		var metrics vegeta.Metrics

		for res := range attacker.Attack(targeter, attackRate, duration, attack.name) {
			metrics.Add(res)
		}
		metrics.Close()

		t.Logf("[%s] requests: %d | success: %.2f%% | p50: %s | p95: %s | p99: %s",
			attack.name,
			metrics.Requests,
			metrics.Success*100,
			metrics.Latencies.P50,
			metrics.Latencies.P95,
			metrics.Latencies.P99,
		)

		if attack.name == "получение слотов" {
			if metrics.Success < 0.999 {
				t.Errorf("[получение слотов] success %.2f%% < 99.9%%", metrics.Success*100)
				allPassed = false
			}
			if metrics.Latencies.P95 > 200*time.Millisecond {
				t.Errorf("[получение слотов] p95 %s > 200ms SLI", metrics.Latencies.P95)
				allPassed = false
			}
		}
	}

	if allPassed {
		t.Log("all SLI checks passed")
	}
}
