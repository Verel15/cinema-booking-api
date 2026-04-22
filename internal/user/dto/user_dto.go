package dto

import (
	"cinema-booking-api/pkg/enums"
	"time"
)

type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=8,max=32"`
	Email    string `json:"email" binding:"required,email"`
	Role     enums.UserRole `json:"role" binding:"required"`
}

type UserResponse struct {
	ID string `json:"id"`
	Username string `json:"username"`
	Email string `json:"email"`
	Role enums.UserRole `json:"role"`
	Status enums.UserStatus `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserFilter struct {
	Status string `form:"status"`
	Role string `form:"role"`
	Username string `form:"username"`
	Email string `form:"email"`
}