package response

import (
	"github.com/gin-gonic/gin"
)

type JSONResponse struct {
	Data       interface{} `json:"data"`
	Success    bool        `json:"success"`
	HttpStatus int         `json:"httpStatus"`
}

type Pagination struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

type PaginatedJSONResponse struct {
	Data       interface{} `json:"data"`
	Success    bool        `json:"success"`
	HttpStatus int         `json:"httpStatus"`
	Pagination Pagination  `json:"pagination"`
}

func Success(c *gin.Context, status int, data interface{}) {
	c.JSON(status, JSONResponse{
		Data:       data,
		Success:    true,
		HttpStatus: status,
	})
}

func SuccessPaginated(c *gin.Context, status int, data interface{}, pagination Pagination) {
	c.JSON(status, PaginatedJSONResponse{
		Data:       data,
		Success:    true,
		HttpStatus: status,
		Pagination: pagination,
	})
}

func Error(c *gin.Context, status int, message string) {
	c.JSON(status, JSONResponse{
		Data:       gin.H{"message": message},
		Success:    false,
		HttpStatus: status,
	})
}
