package domain

import (
	"cinema-booking-api/pkg/enums"
	"time"
)

type Movie struct {
	ID          string            `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Title       string            `gorm:"not null" json:"title"`
	Description string            `json:"description"`
	Duration    int               `json:"duration"`
	Genre       string            `json:"genre"`
	ReleaseDate time.Time         `json:"release_date"`
	PosterURL   string            `json:"poster_url"`
	Status      enums.MovieStatus `gorm:"default:active" json:"status"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}
