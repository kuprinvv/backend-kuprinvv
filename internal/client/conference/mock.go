package conference

import (
	"context"
	"fmt"
	"test-backend-1-kuprinvv/internal/client"

	"github.com/google/uuid"
)

var _ client.ConferenceClient = (*mock)(nil)

type mock struct{}

func NewMock() *mock {
	return &mock{}
}

func (m *mock) CreateLink(_ context.Context) (string, error) {
	return fmt.Sprintf("https://meet.example.com/%s", uuid.New().String()), nil
}
