package controller

import (
	"backend/app/auth"
	"backend/app/dao"
	"backend/app/model"
	"backend/app/utils"
	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"
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

func CreateWorkspace(c *gin.Context) {
	token := c.Request.Header.Get("authentication")

	accountId, err := auth.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	r := new(CreateWorkspaceReq)
	if err := c.Bind(&r); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	workspaceId := ulid.Make().String()
	newWorkspace := model.Workspace{
		Id:        workspaceId,
		Name:      r.WorkspaceName,
		AvatarUrl: r.WorkspaceAvatarUrl,
	}

	if err := dao.CreateWorkspace(&newWorkspace).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	userId := ulid.Make().String()
	newUser := model.User{
		Id:          userId,
		Name:        r.UserName,
		AccountId:   accountId,
		WorkspaceId: workspaceId,
		Role:        "owner",
		AvatarUrl:   r.UserAvatarUrl,
	}

	if err := dao.CreateUser(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, newWorkspace)
	//この後はワークスペース選択画面にリダイレクトor選択画面へのリンクを貼るor作ったワークスペースへリダイレクト
}

func ChangeWorkspaceAttributes(c *gin.Context) {
	workspaceId := utils.GetValueFromContext(c, "workspaceId")
	role := utils.GetValueFromContext(c, "role")

	if role != "owner" {
		c.JSON(http.StatusForbidden, gin.H{"message": "not permitted"})
		return
	}

	r := new(ChangeWorkspaceAttributesReq)
	if err := c.Bind(&r); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	targetWorkspace := model.Workspace{}

	if err := dao.ChangeWorkspaceAttributes(&targetWorkspace, workspaceId, r.WorkspaceName, r.WorkspaceAvatarUrl).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, targetWorkspace)
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

	c.JSON(http.StatusOK, gin.H{"message": "contribution deleted"})
}

func FetchAllWorkSpaces(c *gin.Context) {
	token := c.Request.Header.Get("authentication")

	accountId, err := auth.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	targetUsers := model.Users{}

	if err := dao.FetchAllUsers(&targetUsers, accountId).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	targetWorkspaces := model.Workspaces{}

	for i := 0; i < len(targetUsers); i++ {
		workspace := model.Workspace{}
		if err := dao.FetchWorkspaceInfo(&workspace, targetUsers[i].WorkspaceId).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		targetWorkspaces = append(targetWorkspaces, workspace)
	}

	c.JSON(http.StatusOK, targetWorkspaces)
}
