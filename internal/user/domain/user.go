package domain

import (
	"cinema-booking-api/pkg/enums"
	"time"
)

type User struct {
	ID         string           `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Username   string           `gorm:"not null;unique" json:"username"`
	Password   string           `gorm:"not null" json:"-"`
	Email      string           `gorm:"not null;unique" json:"email"`
	Role       enums.UserRole   `gorm:"default:user" json:"role"`
	Status     enums.UserStatus `gorm:"default:active" json:"status"`
	Provider   string           `gorm:"default:email" json:"provider"` // "email" or "google"
	ProviderID string           `gorm:"column:provider_id" json:"provider_id"`
	AvatarURL  string           `gorm:"column:avatar_url" json:"avatar_url"`
	CreatedAt  time.Time        `json:"created_at"`
	UpdatedAt  time.Time        `json:"updated_at"`
}
