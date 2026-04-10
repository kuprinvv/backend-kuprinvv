package httpx

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
)

func DecodeBody[T any](body io.ReadCloser) (T, error) {
	var payload T
	if err := json.NewDecoder(body).Decode(&payload); err != nil {
		return payload, err
	}
	return payload, nil
}

func Validate[T any](payload T) error {
	return validator.New().Struct(payload)
}

func HandleBody[T any](r *http.Request) (*T, error) {
	body, err := DecodeBody[T](r.Body)
	if err != nil {
		return nil, err
	}
	if err = Validate(body); err != nil {
		return nil, err
	}
	return &body, nil
}

func QueryParam[T any](r *http.Request, key string, parse func(string) (T, error), defaultVal ...T) (T, error) {
	str := r.URL.Query().Get(key)
	if str == "" && len(defaultVal) > 0 {
		return defaultVal[0], nil
	}
	return parse(str)
}

func ParseUUIDParam(r *http.Request, param string) (uuid.UUID, error) {
	raw := chi.URLParam(r, param)
	id, err := uuid.Parse(raw)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid %s: %w", param, err)
	}
	return id, nil
}

func JSON(w http.ResponseWriter, data any, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}

type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

func Error(w http.ResponseWriter, code, message string, statusCode int) {
	JSON(w, ErrorResponse{Error: ErrorDetail{Code: code, Message: message}}, statusCode)
}
