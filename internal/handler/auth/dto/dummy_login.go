package dto

type DummyLoginRequest struct {
	Role string `json:"role" validate:"required,oneof=user admin"`
}

type DummyLoginResponse struct {
	Token string `json:"token"`
}
