package dto

type Response struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
	Data    any    `json:"data,omitempty"`
	Error   any    `json:"error,omitempty"`
}
