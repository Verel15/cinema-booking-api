package delivery

import (
	"cinema-booking-api/internal/user/domain"
	"cinema-booking-api/internal/user/dto"
	"cinema-booking-api/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	usecase domain.UserUsecase
}

func NewUserHandler(u domain.UserUsecase) *UserHandler {
	return &UserHandler{usecase: u}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	res, err := h.usecase.CreateUser(req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusCreated, res)
}

func (h *UserHandler) GetAllUsers(c *gin.Context) {
	var filter dto.UserFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	users, err := h.usecase.GetAllUsers(filter)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusOK, users)
}


func (h *UserHandler) GetUserByUsername(c *gin.Context) {
	username := c.Param("username")
	user, err := h.usecase.GetUserByUsername(username)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	if user == nil {
		response.Error(c, http.StatusNotFound, "User not found")
		return
	}

	res := dto.UserResponse{
		ID: user.ID,
		Username: user.Username,
		Email: user.Email,
		Role: user.Role,
		Status: user.Status,
	}

	response.Success(c, http.StatusOK, res)
}