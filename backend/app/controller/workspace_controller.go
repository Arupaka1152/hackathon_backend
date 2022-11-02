package controller

import (
	"backend/app/auth"
	"backend/app/dao"
	"backend/app/model"
	"backend/app/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type CreateWorkspaceReq struct {
	WorkspaceName      string `json:"workspace_name" binding:"required"`
	WorkspaceAvatarUrl string `json:"workspace_avatar_url"`
	UserName           string `json:"user_name" binding:"required"`
	UserAvatarUrl      string `json:"user_avatar_url"`
}

type ChangeWorkspaceAttributesReq struct {
	WorkspaceName      string `json:"workspace_name" binding:"required"`
	WorkspaceAvatarUrl string `json:"workspace_avatar_url"`
}

type WorkspaceRes struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	AvatarUrl string `json:"avatar_url"`
}

type WorkspacesRes []WorkspaceRes

func CreateWorkspace(c *gin.Context) {
	token := c.Request.Header.Get("authentication")

	accountId, err := auth.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	req := new(CreateWorkspaceReq)
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	workspaceId := utils.GenerateId()
	newWorkspace := model.Workspace{
		Id:        workspaceId,
		Name:      req.WorkspaceName,
		AvatarUrl: req.WorkspaceAvatarUrl,
	}

	if err := dao.CreateWorkspace(&newWorkspace).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	userId := utils.GenerateId()
	newUser := model.User{
		Id:          userId,
		Name:        req.UserName,
		AccountId:   accountId,
		WorkspaceId: workspaceId,
		Role:        "owner",
		AvatarUrl:   req.UserAvatarUrl,
	}

	if err := dao.CreateUser(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	res := &WorkspaceRes{
		workspaceId,
		newWorkspace.Name,
		newWorkspace.AvatarUrl,
	}

	c.JSON(http.StatusOK, res)
}

func ChangeWorkspaceAttributes(c *gin.Context) {
	workspaceId := utils.GetValueFromContext(c, "workspaceId")
	role := utils.GetValueFromContext(c, "role")

	if role != "owner" {
		c.JSON(http.StatusForbidden, gin.H{"message": "not permitted"})
		return
	}

	req := new(ChangeWorkspaceAttributesReq)
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	targetWorkspace := model.Workspace{}
	if err := dao.UpdateWorkspaceAttributes(&targetWorkspace, workspaceId, req.WorkspaceName, req.WorkspaceAvatarUrl).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	res := &WorkspaceRes{
		workspaceId,
		targetWorkspace.Name,
		targetWorkspace.AvatarUrl,
	}

	c.JSON(http.StatusOK, res)
}

func DeleteWorkspace(c *gin.Context) {
	workspaceId := utils.GetValueFromContext(c, "workspaceId")
	role := utils.GetValueFromContext(c, "role")

	if role != "owner" {
		c.JSON(http.StatusForbidden, gin.H{"message": "not permitted"})
		return
	}

	targetWorkspace := model.Workspace{}
	if err := dao.DeleteWorkspace(&targetWorkspace, workspaceId).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	targetUsers := model.Users{}
	if err := dao.DeleteAllUsersInWorkspace(&targetUsers, workspaceId).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "workspace deleted"})
}

func FetchAllWorkSpaces(c *gin.Context) {
	token := c.Request.Header.Get("authentication")

	accountId, err := auth.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	//ここのクエリ文を一つにしたい！！
	targetUsers := model.Users{}
	if err := dao.GetAllUsers(&targetUsers, accountId).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if targetUsers[0].Id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "users not found"})
		return
	}

	targetWorkspaces := model.Workspaces{}
	for i := 0; i < len(targetUsers); i++ {
		workspace := model.Workspace{}
		if err := dao.FindWorkspaceInfo(&workspace, targetUsers[i].WorkspaceId).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		targetWorkspaces = append(targetWorkspaces, workspace)
	}

	res := make(WorkspacesRes, 0)
	for i := 0; i < len(targetWorkspaces); i++ {
		res[i].Id = targetWorkspaces[i].Id
		res[i].Name = targetWorkspaces[i].Name
		res[i].AvatarUrl = targetWorkspaces[i].AvatarUrl
	}

	c.JSON(http.StatusOK, res)
}
