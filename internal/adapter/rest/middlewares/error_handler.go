package middlewares

import (
	"errors"
	"log"
	"net/http"

	"go-api-boilerplate/internal/adapter/rest/restkit"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		err := c.Errors.Last()
		if err == nil {
			return
		}

		var restErr *restkit.Error
		if errors.As(err, &restErr) {
			c.JSON(restErr.Status, gin.H{
				"code":    restErr.Code,
				"message": restErr.Message,
			})
		} else {
			log.Printf("[INTERNAL_ERROR]: %+v\n", err)

			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    restkit.ErrCodeInternalError,
				"message": restkit.ErrorMessages["internal_error"],
			})
		}
	}
}
