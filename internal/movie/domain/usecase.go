package domain

import "cinema-booking-api/internal/movie/dto"

type MovieUsecase interface {
	CreateMovie(req dto.CreateMovieRequest) (*dto.MovieResponse, error)
	GetAllMovies(filter dto.MovieFilter) ([]dto.MovieResponse, error)
	GetMovieByID(id string) (*dto.MovieResponse, error)
	UpdateMovie(id string, req dto.UpdateMovieRequest) (*dto.MovieResponse, error)
	DeleteMovie(id string) error
}
