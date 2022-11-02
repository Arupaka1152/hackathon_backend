package controller

import (
	"backend/app/dao"
	"backend/app/model"
	"backend/app/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type CreateUserReq struct {
	Email     string `json:"email" binding:"required"`
	Name      string `json:"name" binding:"required"`
	Role      string `json:"role" binding:"required"`
	AvatarUrl string `json:"avatar_url"`
}

type GrantRoleToUserReq struct {
	UserId string `json:"user_id" binding:"required"`
	Role   string `json:"role" binding:"required"`
}

type ChangeUserAttributesReq struct {
	UserName      string `json:"user_name" binding:"required"`
	UserAvatarUrl string `json:"user_avatar_url"`
}

type RemoveUserFromWorkspaceReq struct {
	UserId string `json:"user_id" binding:"required"`
}

type UserRes struct {
	Id          string `json:"user_id"`
	Name        string `json:"name"`
	AccountId   string `json:"account_id"`
	WorkspaceId string `json:"workspace_id"`
	Role        string `json:"role"`
	AvatarUrl   string `json:"avatar_url"`
}

type UsersRes []UserRes

type GrantRoleToUserRes struct {
	Id   string `json:"user_id"`
	Role string `json:"role"`
}

type ChangeUserAttributesRes struct {
	Name      string `json:"name"`
	AvatarUrl string `json:"avatar_url"`
}

func CreateUser(c *gin.Context) {
	workspaceId := utils.GetValueFromContext(c, "workspaceId")
	role := utils.GetValueFromContext(c, "role")

	if role != "manager" && role != "owner" {
		c.JSON(http.StatusForbidden, gin.H{"message": "not permitted"})
		return
	}

	req := new(CreateUserReq)
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	targetAccount := model.Account{}
	if err := dao.FindAccountByEmail(&targetAccount, req.Email).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if targetAccount.Id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "account not found"})
		return
	}

	newUserId := utils.GenerateId()
	newUser := model.User{
		Id:          newUserId,
		Name:        req.Name,
		AccountId:   targetAccount.Id,
		WorkspaceId: workspaceId,
		Role:        req.Role,
		AvatarUrl:   req.AvatarUrl,
	}

	if err := dao.CreateUser(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	res := &UserRes{
		newUser.Id,
		newUser.Name,
		newUser.AccountId,
		newUser.WorkspaceId,
		newUser.Role,
		newUser.AvatarUrl,
	}

	c.JSON(http.StatusOK, res)
}

func FetchAllUsersInWorkspace(c *gin.Context) {
	workspaceId := utils.GetValueFromContext(c, "workspaceId")

	targetUsers := model.Users{}
	if err := dao.GetAllUsersInWorkspace(&targetUsers, workspaceId).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	res := make(UsersRes, 0)
	for i := 0; i < len(targetUsers); i++ {
		res[i].Id = targetUsers[i].Id
		res[i].Name = targetUsers[i].Name
		res[i].AccountId = targetUsers[i].AccountId
		res[i].WorkspaceId = targetUsers[i].WorkspaceId
		res[i].Role = targetUsers[i].Role
		res[i].AvatarUrl = targetUsers[i].AvatarUrl
	}

	c.JSON(http.StatusOK, res)
}

func RemoveUserFromWorkspace(c *gin.Context) {
	userId := utils.GetValueFromContext(c, "userId")
	role := utils.GetValueFromContext(c, "role")

	if role != "owner" {
		c.JSON(http.StatusForbidden, gin.H{"message": "not permitted"})
		return
	}

	req := new(RemoveUserFromWorkspaceReq)
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if userId == req.UserId {
		c.JSON(http.StatusBadRequest, gin.H{"error": "you cant remove yourself"})
		return
	}

	targetUser := model.User{}
	if err := dao.DeleteUser(&targetUser, req.UserId).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}

func DeleteUser(c *gin.Context) {
	userId := utils.GetValueFromContext(c, "userId")
	role := utils.GetValueFromContext(c, "role")

	if role == "owner" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "you cant delete yourself"})
		return
	}

	targetUser := model.User{}
	if err := dao.DeleteUser(&targetUser, userId).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}

func GrantRoleToUser(c *gin.Context) {
	userId := utils.GetValueFromContext(c, "userId")
	role := utils.GetValueFromContext(c, "role")

	if role != "manager" && role != "owner" {
		c.JSON(http.StatusForbidden, gin.H{"message": "not permitted"})
		return
	}

	req := new(GrantRoleToUserReq)
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if req.Role != "manager" && req.Role != "general" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "choose manager or general"})
		return
	}

	if userId == req.UserId {
		c.JSON(http.StatusBadRequest, gin.H{"error": "you cant change your role"})
		return
	}

	targetUser := model.User{}
	if err := dao.UpdateRoleOfUser(&targetUser, req.UserId, req.Role).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	res := &GrantRoleToUserRes{
		req.UserId,
		targetUser.Role,
	}

	c.JSON(http.StatusOK, res)
}

func ChangeUserAttributes(c *gin.Context) {
	userId := utils.GetValueFromContext(c, "userId")

	req := new(ChangeUserAttributesReq)
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	targetUser := model.User{}
	if err := dao.UpdateUserAttributes(&targetUser, userId, req.UserName, req.UserAvatarUrl).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	res := &ChangeUserAttributesRes{
		targetUser.Name,
		targetUser.AvatarUrl,
	}

	c.JSON(http.StatusOK, res)
}
