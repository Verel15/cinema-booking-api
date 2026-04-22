package domain

import "cinema-booking-api/internal/user/dto"

type UserUsecase interface {
	CreateUser(req dto.CreateUserRequest) (*dto.UserResponse,error)
	GetAllUsers(filter dto.UserFilter) ([]dto.UserResponse,error)
	GetUserByUsername(username string) (*User, error)
}