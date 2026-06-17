package dto

type Response struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
	Error   string `json:"error"`
}
