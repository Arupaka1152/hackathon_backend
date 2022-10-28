package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func RespondWithError(c *gin.Context, status int, message interface{}) {
	c.AbortWithStatusJSON(status, gin.H{"error": message})
}

func GetValueFromContext(c *gin.Context, key string) (valueString string) {
	value, exists := c.Get(key)
	if !exists {
		RespondWithError(c, http.StatusInternalServerError, "cant get value from gin.Context")
	}
	return value.(string)
}
