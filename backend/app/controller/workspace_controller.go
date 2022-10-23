package controller

import (
	"backend/app/auth"
	"backend/app/dao"
	"backend/app/model"
	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"
	"net/http"
)

type CreateWorkspaceReq struct {
	WorkspaceName      string `json:"workspace_name"`
	WorkspaceAvatarUrl string `json:"workspace_avatar_url"`
	UserName           string `json:"user_name"`
	UserAvatarUrl      string `json:"user_avatar_url"`
}

type ChangeWorkspaceAttributesReq struct {
	WorkspaceName      string `json:"workspace_name"`
	WorkspaceAvatarUrl string `json:"workspace_avatar_url"`
}

func CreateWorkspace(c *gin.Context) {
	token := c.Request.Header.Get("authentication")

	accountId, err := auth.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	}

	r := new(CreateWorkspaceReq)
	if err := c.Bind(&r); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	workspaceId := ulid.Make().String()
	newWorkspace := model.Workspace{
		Id:        workspaceId,
		Name:      r.WorkspaceName,
		AvatarUrl: r.WorkspaceAvatarUrl,
	}

	if err := dao.CreateWorkspace(&newWorkspace).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
	}

	c.JSON(http.StatusOK, newWorkspace)
	//この後はワークスペース選択画面にリダイレクトor選択画面へのリンクを貼るor作ったワークスペースへリダイレクト
}

func ChangeWorkspaceAttributes(c *gin.Context) {
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

	if role != "owner" {
		c.JSON(http.StatusForbidden, gin.H{"message": "not permitted"})
	}

	r := new(ChangeWorkspaceAttributesReq)
	if err := c.Bind(&r); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	targetWorkspace := model.Workspace{}

	if err := dao.ChangeWorkspaceAttributes(&targetWorkspace, workspaceId, r.WorkspaceName, r.WorkspaceAvatarUrl).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, targetWorkspace)
}

func DeleteWorkspace(c *gin.Context) {
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

	if role != "owner" {
		c.JSON(http.StatusForbidden, gin.H{"message": "not permitted"})
	}

	targetWorkspace := model.Workspace{}

	if err := dao.DeleteWorkspace(&targetWorkspace, workspaceId).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"message": "contribution deleted"})
}

func FetchAllWorkSpaces(c *gin.Context) {
	token := c.Request.Header.Get("authentication")

	accountId, err := auth.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	}

	targetUsers := model.Users{}

	if err := dao.FetchAllUsers(&targetUsers, accountId).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	targetWorkspaces := model.Workspaces{}

	for i := 0; i < len(targetUsers); i++ {
		workspace := model.Workspace{}
		if err := dao.FetchWorkspaceInfo(&workspace, targetUsers[i].WorkspaceId).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		targetWorkspaces = append(targetWorkspaces, workspace)
	}

	c.JSON(http.StatusOK, targetWorkspaces)
}
