package errors

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func BadRequest(c *gin.Context) {
	c.JSON(c.Writer.Status(), gin.H{
		"status": http.StatusBadRequest,
		"message": "BadRequest",
	})
	c.Abort()
}

func InternalServerError(c *gin.Context) {
	c.JSON(c.Writer.Status(), gin.H{
		"status": http.StatusInternalServerError,
		"message": "InternalServerError",
	})
	c.Abort()
}

func Unauthorized(c *gin.Context) {
	c.JSON(c.Writer.Status(), gin.H{
		"status": http.StatusUnauthorized,
		"message": "Unauthorized",
	})
	c.Abort()
}