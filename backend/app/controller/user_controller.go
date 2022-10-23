package controller

import (
	"backend/app/auth"
	"backend/app/dao"
	"backend/app/model"
	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"
	"net/http"
)

type CreateUserReq struct {
	Email     string `json:"email"`
	Name      string `json:"name"`
	Role      string `json:"role"`
	AvatarUrl string `json:"avatar_url"`
}

type GrantRoleToUserReq struct {
	UserId string `json:"user_id"`
	Role   string `json:"role"`
}

type ChangeUserAttributesReq struct {
	UserName      string `json:"workspace_name"`
	UserAvatarUrl string `json:"workspace_avatar_url"`
}

type RemoveUserFromWorkspaceReq struct {
	UserId string `json:"user_id"`
}

func CreateUser(c *gin.Context) {
	workspaceId := c.Request.Header.Get("workspace_id")
	token := c.Request.Header.Get("authentication")

	accountId, err := auth.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	}

	_, role, err := auth.UserAuth(workspaceId, accountId)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "not permitted"})
	}

	if role != "manager" && role != "owner" {
		c.JSON(http.StatusForbidden, gin.H{"message": "not permitted"})
	}

	r := new(CreateUserReq)
	if err := c.Bind(&r); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	targetAccount := model.Account{}
	if err := dao.FetchAccountByEmail(&targetAccount, r.Email).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "account not found"})
	}

	newUserId := ulid.Make().String()
	newUser := model.User{
		Id:          newUserId,
		Name:        r.Name,
		AccountId:   targetAccount.Id,
		WorkspaceId: workspaceId,
		Role:        r.Role,
		AvatarUrl:   r.AvatarUrl,
	}

	if err := dao.CreateUser(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, newUser)
}

func FetchAllUsersInWorkspace(c *gin.Context) {
	workspaceId := c.Request.Header.Get("workspace_id")
	token := c.Request.Header.Get("authentication")

	accountId, err := auth.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	}

	_, _, err = auth.UserAuth(workspaceId, accountId)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "not permitted"})
	}

	targetUsers := model.Users{}
	if err := dao.FetchAllUsersInWorkspace(&targetUsers, workspaceId).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, targetUsers)
}

func RemoveUserFromWorkspace(c *gin.Context) {
	workspaceId := c.Request.Header.Get("workspace_id")
	token := c.Request.Header.Get("authentication")

	accountId, err := auth.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	}

	userId, role, err := auth.UserAuth(workspaceId, accountId)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "not permitted"})
	}

	if role != "manager" && role != "owner" {
		c.JSON(http.StatusForbidden, gin.H{"message": "not permitted"})
	}

	r := new(RemoveUserFromWorkspaceReq)
	if err := c.Bind(&r); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	if userId == r.UserId {
		c.JSON(http.StatusBadRequest, gin.H{"error": "you cant delete yourself"})
	}

	targetUser := model.User{}
	if err := dao.DeleteUser(&targetUser, r.UserId).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}

func DeleteUser(c *gin.Context) {
	workspaceId := c.Request.Header.Get("workspace_id")
	token := c.Request.Header.Get("authentication")

	accountId, err := auth.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	}

	userId, role, err := auth.UserAuth(workspaceId, accountId)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "not permitted"})
	}

	if role == "owner" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "you cant delete yourself"})
	}

	targetUser := model.User{}
	if err := dao.DeleteUser(&targetUser, userId).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}

func GrantRoleToUser(c *gin.Context) {
	workspaceId := c.Request.Header.Get("workspace_id")
	token := c.Request.Header.Get("authentication")

	accountId, err := auth.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	}

	_, role, err := auth.UserAuth(workspaceId, accountId)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "not permitted"})
	}

	if role != "manager" && role != "owner" {
		c.JSON(http.StatusForbidden, gin.H{"message": "not permitted"})
	}

	r := new(GrantRoleToUserReq)
	if err := c.Bind(&r); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	if r.Role != "manager" && r.Role != "general" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "choose manager or general"})
	}

	targetUser := model.User{}
	if err := dao.GrantRoleToUser(&targetUser, r.UserId, r.Role).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, targetUser)
}

func ChangeUserAttributes(c *gin.Context) {
	workspaceId := c.Request.Header.Get("workspace_id")
	token := c.Request.Header.Get("authentication")

	accountId, err := auth.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	}

	userId, _, err := auth.UserAuth(workspaceId, accountId)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "not permitted"})
	}

	r := new(ChangeUserAttributesReq)
	if err := c.Bind(&r); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	targetUser := model.User{}

	if err := dao.ChangeUserAttributes(&targetUser, userId, r.UserName, r.UserAvatarUrl).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, targetUser)
}
