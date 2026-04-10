package converter

import (
	"test-backend-1-kuprinvv/internal/model"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestServiceToRegisterResponse(t *testing.T) {
	now := time.Now()
	user := model.User{
		ID:        uuid.New(),
		Email:     "test@example.com",
		Role:      "user",
		CreatedAt: &now,
	}

	resp := ServiceToRegisterResponse(user)

	assert.Equal(t, user.ID, resp.User.ID)
	assert.Equal(t, user.Email, resp.User.Email)
	assert.Equal(t, user.Role, resp.User.Role)
	assert.Equal(t, user.CreatedAt, resp.User.CreatedAt)
}
