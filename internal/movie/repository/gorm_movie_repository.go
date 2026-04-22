package repository

import (
	"cinema-booking-api/internal/movie/domain"
	"cinema-booking-api/internal/movie/dto"
	"gorm.io/gorm"
)

type gormMovieRepository struct {
	db *gorm.DB
}

func NewMovieRepository(db *gorm.DB) domain.MovieRepository {
	return &gormMovieRepository{db: db}
}

func (r *gormMovieRepository) Create(m *domain.Movie) error {
	return r.db.Create(m).Error
}

func (r *gormMovieRepository) GetAll(filter dto.MovieFilter) ([]domain.Movie, error) {
	var movies []domain.Movie
	query := r.db.Model(&domain.Movie{})

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	if filter.Title != "" {
		query = query.Where("title ILIKE ?", "%"+filter.Title+"%")
	}

	if filter.Genre != "" {
		query = query.Where("genre = ?", filter.Genre)
	}

	err := query.Find(&movies).Error
	return movies, err
}

func (r *gormMovieRepository) GetByID(id string) (*domain.Movie, error) {
	var m domain.Movie
	err := r.db.First(&m, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func (r *gormMovieRepository) Update(m *domain.Movie) error {
	return r.db.Save(m).Error
}

func (r *gormMovieRepository) Delete(id string) error {
	return r.db.Delete(&domain.Movie{}, "id = ?", id).Error
}
