package controller

import (
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
	//作る時点ではワークスペースは存在していないのでヘッダーにはワークスペースIDはない
	//アクセストークンによる認証だけでOK アカウントIDを取り出しておく
	accountId := "fkdlasafkdsl"

	req := new(CreateWorkspaceReq)
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	workspaceId := ulid.Make().String()
	newWorkspace := model.Workspace{
		Id:        workspaceId,
		Name:      req.WorkspaceName,
		AvatarUrl: req.WorkspaceAvatarUrl,
	}

	if err := dao.CreateWorkspace(&newWorkspace).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	userId := ulid.Make().String()
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
	}

	c.JSON(http.StatusOK, newWorkspace)
	//この後はワークスペース選択画面にリダイレクトor選択画面へのリンクを貼るor作ったワークスペースへリダイレクト
}

func ChangeWorkspaceAttributes(c *gin.Context) {
	workspaceId := c.Request.Header.Get("workspace_id")
	//ここでアクセストークンをデコードしてアカウントIDを取得する accountId := ...
	//取得したworkspaceIdとaccountIdを使ってデータベースを参照しuserIdを取得する、できなかったら認証エラーを吐くようにする
	//ownerのみが変更できるようにする

	req := new(ChangeWorkspaceAttributesReq)
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	targetWorkspace := model.Workspace{}

	if err := dao.ChangeWorkspaceAttributes(&targetWorkspace, workspaceId, req.WorkspaceName, req.WorkspaceAvatarUrl).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, targetWorkspace)
}

func DeleteWorkspace(c *gin.Context) {
	workspaceId := c.Request.Header.Get("workspace_id")
	//ここでアクセストークンをデコードしてアカウントIDを取得する accountId := ...
	//取得したworkspaceIdとaccountIdを使ってデータベースを参照しuserIdを取得する、できなかったら認証エラーを吐くようにする
	//ownerのみが削除できるようにする

	targetWorkspace := model.Workspace{}

	if err := dao.DeleteWorkspace(&targetWorkspace, workspaceId).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"message": "contribution deleted"})
}

func FetchAllWorkSpaces(c *gin.Context) {
	//ここでアクセストークンをデコードしてアカウントIDを取得する accountId := ...
	//取得したworkspaceIdとaccountIdを使ってデータベースを参照しuserIdを取得する、できなかったら認証エラーを吐くようにする
	accountId := "fakdslsadfkl"

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
