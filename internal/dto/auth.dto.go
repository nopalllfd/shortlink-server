package dto

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

type UserLoginResponse struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
}

type AuthResponse struct {
	User  UserLoginResponse `json:"user"`
	Token string            `json:"token,omitempty"`
}
