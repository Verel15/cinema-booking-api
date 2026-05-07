package repository

import (
	"cinema-booking-api/internal/user/domain"
	"cinema-booking-api/internal/user/dto"

	"gorm.io/gorm"
)

type gormUserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &gormUserRepository{db: db}
}

func (r *gormUserRepository) Create(user *domain.User) error {
	return r.db.Create(user).Error

}

func (r *gormUserRepository) GetAll(filter dto.UserFilter) ([]domain.User, int64, error) {
	var users []domain.User
	var total int64

	query := r.db.Model(&domain.User{})

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.Role != "" {
		query = query.Where("role = ?", filter.Role)
	}
	if filter.Username != "" {
		query = query.Where("username ILIKE ?", "%"+filter.Username+"%")
	}
	if filter.Email != "" {
		query = query.Where("email ILIKE ?", "%"+filter.Email+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Offset(filter.Offset()).Limit(filter.Limit).Find(&users).Error
	return users, total, err
}

func (r *gormUserRepository) GetByUsername(username string) (*domain.User, error) {
	var user domain.User
	err := r.db.First(&user, "username = ?", username).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
