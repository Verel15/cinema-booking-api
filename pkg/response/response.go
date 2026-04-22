package response

import (
	"github.com/gin-gonic/gin"
)

type JSONResponse struct {
	Data       interface{} `json:"data"`
	Success    bool        `json:"success"`
	HttpStatus int         `json:"httpStatus"`
}

func Success(c *gin.Context, status int, data interface{}) {
	c.JSON(status, JSONResponse{
		Data:       data,
		Success:    true,
		HttpStatus: status,
	})
}

func Error(c *gin.Context, status int, message string) {
	c.JSON(status, JSONResponse{
		Data:       gin.H{"message": message},
		Success:    false,
		HttpStatus: status,
	})
}
