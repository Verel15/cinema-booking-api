package dto

import (
	"cinema-booking-api/pkg/enums"
	"time"
)

type CreateMovieRequest struct {
	Title       string            `json:"title" binding:"required"`
	Description string            `json:"description"`
	Duration    int               `json:"duration" binding:"required,gt=0"`
	Genre       string            `json:"genre" binding:"required"`
	ReleaseDate time.Time         `json:"release_date" binding:"required"`
	PosterURL   string            `json:"poster_url" binding:"url"`
	Status      enums.MovieStatus `json:"status" binding:"omitempty,oneof=active inactive deleted"`
}

type UpdateMovieRequest struct {
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Duration    int               `json:"duration" binding:"omitempty,gt=0"`
	Genre       string            `json:"genre"`
	ReleaseDate time.Time         `json:"release_date"`
	PosterURL   string            `json:"poster_url" binding:"omitempty,url"`
	Status      enums.MovieStatus `json:"status" binding:"omitempty,oneof=active inactive deleted"`
}

type MovieResponse struct {
	ID          string            `json:"id"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Duration    int               `json:"duration"`
	Genre       string            `json:"genre"`
	ReleaseDate time.Time         `json:"release_date"`
	PosterURL   string            `json:"poster_url"`
	Status      enums.MovieStatus `json:"status"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

type MovieFilter struct {
	Status string `form:"status"`
	Title  string `form:"title"`
	Genre  string `form:"genre"`
}
