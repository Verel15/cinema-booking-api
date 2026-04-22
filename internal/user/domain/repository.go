package domain

import "cinema-booking-api/internal/user/dto"

type UserRepository interface {
	Create(user *User) error
	GetAll(filter dto.UserFilter) ([]User, error)
	GetByUsername(username string) (*User, error)
}

type AuthRepository interface {
	FindByEmail(email string) (*User, error)
	FindByProviderID(provider, providerID string) (*User, error)
	Create(user *User) error
	Update(user *User) error
	FindByID(id string) (*User, error)
}

type AuthUsecase interface {
	Register(req interface{}) (interface{}, error)
	Login(req interface{}) (interface{}, error)
	LoginWithGoogle(code string) (interface{}, error)
	RefreshToken(refreshToken string) (interface{}, error)
	ValidateToken(token string) (*User, error)
	GetGoogleAuthURL() string
}
