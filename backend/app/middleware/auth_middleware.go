package middleware

import (
	"backend/app/auth"
	"backend/app/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		workspaceId := c.Request.Header.Get("workspace_id")
		token := c.Request.Header.Get("authentication")

		accountId, err := auth.ParseToken(token)
		if err != nil {
			utils.RespondWithError(c, http.StatusUnauthorized, "invalid token")
			return
		}

		userId, role, err := auth.UserAuth(workspaceId, accountId)
		if err != nil {
			utils.RespondWithError(c, http.StatusForbidden, "not permitted")
			return
		}

		c.Set("workspaceId", workspaceId)
		c.Set("accountId", accountId)
		c.Set("userId", userId)
		c.Set("role", role)
		c.Next()
	}
}
