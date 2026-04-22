package http

import (
	"cinema-booking-api/internal/movie/domain"
	"cinema-booking-api/internal/movie/dto"
	"cinema-booking-api/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MovieHandler struct {
	usecase domain.MovieUsecase
}

func NewMovieHandler(u domain.MovieUsecase) *MovieHandler {
	return &MovieHandler{usecase: u}
}

func (h *MovieHandler) CreateMovie(c *gin.Context) {
	var req dto.CreateMovieRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.usecase.CreateMovie(req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusCreated, res)
}

func (h *MovieHandler) GetAllMovies(c *gin.Context) {
	var filter dto.MovieFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	movies, err := h.usecase.GetAllMovies(filter)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusOK, movies)
}

func (h *MovieHandler) GetMovieByID(c *gin.Context) {
	id := c.Param("id")

	res, err := h.usecase.GetMovieByID(id)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	if res == nil {
		response.Error(c, http.StatusNotFound, "movie not found")
		return
	}

	response.Success(c, http.StatusOK, res)
}

func (h *MovieHandler) UpdateMovie(c *gin.Context) {
	id := c.Param("id")

	var req dto.UpdateMovieRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.usecase.UpdateMovie(id, req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	if res == nil {
		response.Error(c, http.StatusNotFound, "movie not found")
		return
	}

	response.Success(c, http.StatusOK, res)
}

func (h *MovieHandler) DeleteMovie(c *gin.Context) {
	id := c.Param("id")

	if err := h.usecase.DeleteMovie(id); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "movie deleted successfully"})
}
