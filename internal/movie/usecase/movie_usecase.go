package usecase

import (
	"cinema-booking-api/internal/movie/domain"
	"cinema-booking-api/internal/movie/dto"
)

type movieUsecase struct {
	repo domain.MovieRepository
}

func NewMovieUsecase(repo domain.MovieRepository) domain.MovieUsecase {
	return &movieUsecase{repo: repo}
}

func (u *movieUsecase) CreateMovie(req dto.CreateMovieRequest) (*dto.MovieResponse, error) {
	m := &domain.Movie{
		Title:       req.Title,
		Description: req.Description,
		Duration:    req.Duration,
		Genre:       req.Genre,
		ReleaseDate: req.ReleaseDate,
		PosterURL:   req.PosterURL,
		Status:      req.Status,
	}

	if m.Status == "" {
		m.Status = "active"
	}

	if err := u.repo.Create(m); err != nil {
		return nil, err
	}

	return u.mapToResponse(m), nil
}

func (u *movieUsecase) GetAllMovies(filter dto.MovieFilter) ([]dto.MovieResponse, error) {
	movies, err := u.repo.GetAll(filter)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.MovieResponse, len(movies))
	for i, m := range movies {
		responses[i] = *u.mapToResponse(&m)
	}

	return responses, nil
}

func (u *movieUsecase) GetMovieByID(id string) (*dto.MovieResponse, error) {
	m, err := u.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, nil
	}

	return u.mapToResponse(m), nil
}

func (u *movieUsecase) UpdateMovie(id string, req dto.UpdateMovieRequest) (*dto.MovieResponse, error) {
	m, err := u.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, nil
	}

	// Update fields if provided
	if req.Title != "" {
		m.Title = req.Title
	}
	if req.Description != "" {
		m.Description = req.Description
	}
	if req.Duration > 0 {
		m.Duration = req.Duration
	}
	if req.Genre != "" {
		m.Genre = req.Genre
	}
	if !req.ReleaseDate.IsZero() {
		m.ReleaseDate = req.ReleaseDate
	}
	if req.PosterURL != "" {
		m.PosterURL = req.PosterURL
	}
	if req.Status != "" {
		m.Status = req.Status
	}

	if err := u.repo.Update(m); err != nil {
		return nil, err
	}

	return u.mapToResponse(m), nil
}

func (u *movieUsecase) DeleteMovie(id string) error {
	return u.repo.Delete(id)
}

func (u *movieUsecase) mapToResponse(m *domain.Movie) *dto.MovieResponse {
	return &dto.MovieResponse{
		ID:          m.ID,
		Title:       m.Title,
		Description: m.Description,
		Duration:    m.Duration,
		Genre:       m.Genre,
		ReleaseDate: m.ReleaseDate,
		PosterURL:   m.PosterURL,
		Status:      m.Status,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}
