package httpx_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"test-backend-1-kuprinvv/pkg/httpx"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecodeBody(t *testing.T) {
	type payload struct {
		Name string `json:"name"`
	}

	t.Run("валидный JSON", func(t *testing.T) {
		body := strings.NewReader(`{"name":"test"}`)
		got, err := httpx.DecodeBody[payload](io.NopCloser(body))
		require.NoError(t, err)
		assert.Equal(t, "test", got.Name)
	})

	t.Run("невалидный JSON", func(t *testing.T) {
		body := strings.NewReader(`not json`)
		_, err := httpx.DecodeBody[payload](io.NopCloser(body))
		assert.Error(t, err)
	})
}

func TestValidate(t *testing.T) {
	type payload struct {
		Name string `validate:"required"`
		Age  int    `validate:"min=1"`
	}

	t.Run("валидная структура", func(t *testing.T) {
		assert.NoError(t, httpx.Validate(payload{Name: "Alice", Age: 25}))
	})

	t.Run("пустое required поле", func(t *testing.T) {
		assert.Error(t, httpx.Validate(payload{Name: "", Age: 25}))
	})

	t.Run("нарушение min", func(t *testing.T) {
		assert.Error(t, httpx.Validate(payload{Name: "Alice", Age: 0}))
	})
}

func TestHandleBody(t *testing.T) {
	type payload struct {
		Role string `json:"role" validate:"required,oneof=user admin"`
	}

	makeReq := func(body string) *http.Request {
		r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
		r.Header.Set("Content-Type", "application/json")
		return r
	}

	t.Run("успешно", func(t *testing.T) {
		p, err := httpx.HandleBody[payload](makeReq(`{"role":"admin"}`))
		require.NoError(t, err)
		assert.Equal(t, "admin", p.Role)
	})

	t.Run("невалидный JSON", func(t *testing.T) {
		_, err := httpx.HandleBody[payload](makeReq(`not json`))
		assert.Error(t, err)
	})

	t.Run("не проходит валидацию", func(t *testing.T) {
		_, err := httpx.HandleBody[payload](makeReq(`{"role":"superuser"}`))
		assert.Error(t, err)
	})
}

func TestQueryParam(t *testing.T) {
	atoi := func(s string) (int, error) {
		var n int
		_, err := fmt.Sscan(s, &n)
		return n, err
	}

	makeReq := func(query string) *http.Request {
		return httptest.NewRequest(http.MethodGet, "/?"+query, nil)
	}

	t.Run("параметр присутствует", func(t *testing.T) {
		v, err := httpx.QueryParam(makeReq("page=3"), "page", atoi)
		require.NoError(t, err)
		assert.Equal(t, 3, v)
	})

	t.Run("параметр отсутствует — default", func(t *testing.T) {
		v, err := httpx.QueryParam(makeReq(""), "page", atoi, 1)
		require.NoError(t, err)
		assert.Equal(t, 1, v)
	})

	t.Run("параметр отсутствует — без default", func(t *testing.T) {
		_, err := httpx.QueryParam(makeReq(""), "page", atoi)
		assert.Error(t, err)
	})

	t.Run("невалидное значение", func(t *testing.T) {
		_, err := httpx.QueryParam(makeReq("page=abc"), "page", atoi)
		assert.Error(t, err)
	})
}

func withChiParam(r *http.Request, key, value string) *http.Request {
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add(key, value)
	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
}

func TestParseUUIDParam(t *testing.T) {
	id := uuid.New()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	t.Run("валидный UUID", func(t *testing.T) {
		req := withChiParam(r, "id", id.String())
		got, err := httpx.ParseUUIDParam(req, "id")
		require.NoError(t, err)
		assert.Equal(t, id, got)
	})

	t.Run("невалидный UUID", func(t *testing.T) {
		req := withChiParam(r, "id", "not-a-uuid")
		_, err := httpx.ParseUUIDParam(req, "id")
		assert.Error(t, err)
	})
}

func TestJSON(t *testing.T) {
	w := httptest.NewRecorder()
	httpx.JSON(w, map[string]string{"key": "value"}, http.StatusCreated)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var got map[string]string
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &got))
	assert.Equal(t, "value", got["key"])
}

func TestError(t *testing.T) {
	w := httptest.NewRecorder()
	httpx.Error(w, "NOT_FOUND", "resource not found", http.StatusNotFound)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var resp httpx.ErrorResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "NOT_FOUND", resp.Error.Code)
	assert.Equal(t, "resource not found", resp.Error.Message)
}
