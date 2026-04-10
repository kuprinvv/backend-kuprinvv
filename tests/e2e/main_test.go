package e2e

import (
	"os"
	"testing"
)

var baseURL string

func TestMain(m *testing.M) {
	baseURL = os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}
	os.Exit(m.Run())
}
