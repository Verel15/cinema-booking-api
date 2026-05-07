package domain

import "cinema-booking-api/internal/movie/dto"

type MovieRepository interface {
	Create(movie *Movie) error
	GetAll(filter dto.MovieFilter) ([]Movie, int64, error)
	GetByID(id string) (*Movie, error)
	Update(movie *Movie) error
	Delete(id string) error
}
