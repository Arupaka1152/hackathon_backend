package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"
	"math/rand"
	"net/http"
	"time"
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

func GenerateId() string {
	t := time.Now()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	id := ulid.MustNew(ulid.Timestamp(t), entropy)
	return id.String()
}
