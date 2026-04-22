package usecase

import (
	"cinema-booking-api/internal/user/domain"
	"cinema-booking-api/internal/user/dto"
	"cinema-booking-api/pkg/utils"
)

type userUsecase struct {
	repo domain.UserRepository
}

func NewUserUsecase (repo domain.UserRepository) domain.UserUsecase {
	return &userUsecase{repo: repo}
}

func (u *userUsecase) CreateUser(req dto.CreateUserRequest) (*dto.UserResponse, error) {
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}
	user := &domain.User{
		Username: req.Username,
		Password: hashedPassword,
		Email: req.Email,
		Role: req.Role,
	}

	if err := u.repo.Create(user); err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID: user.ID,
		Username: user.Username,
		Email: user.Email,
		Role: user.Role,
		Status: user.Status,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}


func (u *userUsecase) GetAllUsers(filter dto.UserFilter) ([]dto.UserResponse, error) {
	users, err := u.repo.GetAll(filter)
	if err != nil {
		return nil, err
	}

	var res []dto.UserResponse
	for _, user := range users {
		res = append(res, dto.UserResponse{
			ID: user.ID,
			Username: user.Username,
			Email: user.Email,
			Role: user.Role,
			Status: user.Status,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
	}
	return res, nil
}

func (u *userUsecase) GetUserByUsername(username string) (*domain.User, error) {
	return u.repo.GetByUsername(username)
}