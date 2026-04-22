package repository

import (
	userDomain "cinema-booking-api/internal/user/domain"

	"gorm.io/gorm"
)

type gormAuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) userDomain.AuthRepository {
	return &gormAuthRepository{db: db}
}

func (r *gormAuthRepository) FindByEmail(email string) (*userDomain.User, error) {
	var user userDomain.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *gormAuthRepository) FindByProviderID(provider, providerID string) (*userDomain.User, error) {
	var user userDomain.User
	err := r.db.Where("provider = ? AND provider_id = ?", provider, providerID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *gormAuthRepository) Create(user *userDomain.User) error {
	return r.db.Create(user).Error
}

func (r *gormAuthRepository) Update(user *userDomain.User) error {
	return r.db.Save(user).Error
}

func (r *gormAuthRepository) FindByID(id string) (*userDomain.User, error) {
	var user userDomain.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
