package controller

import (
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

func CreateUser(c *gin.Context) {
	workspaceId := c.Request.Header.Get("workspace_id")
	//ここでアクセストークンをデコードしてアカウントIDを取得する accountId := ...
	//取得したworkspaceIdとaccountIdを使ってデータベースを参照しuserIdを取得する、できなかったら認証エラーを吐くようにする
	//ownerかmanagerが招待できるようにする

	req := new(CreateUserReq)
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	targetAccount := model.Account{}
	if err := dao.FetchAccountByEmail(&targetAccount, req.Email).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "account not found"})
	}

	userId := ulid.Make().String()
	newUser := model.User{
		Id:          userId,
		Name:        req.Name,
		AccountId:   targetAccount.Id,
		WorkspaceId: workspaceId,
		Role:        req.Role,
		AvatarUrl:   req.AvatarUrl,
	}

	if err := dao.CreateUser(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, newUser)
}

func FetchAllUsersInWorkspace(c *gin.Context) {
	workspaceId := c.Request.Header.Get("workspace_id")
	//ここでアクセストークンをデコードしてアカウントIDを取得する accountId := ...
	//取得したworkspaceIdとaccountIdを使ってデータベースを参照しuserIdを取得する、できなかったら認証エラーを吐くようにする

	targetUsers := model.Users{}
	if err := dao.FetchAllUsersInWorkspace(&targetUsers, workspaceId).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, targetUsers)
}

func RemoveUserFromWorkspace(c *gin.Context) {
	workspaceId := c.Request.Header.Get("workspace_id")
	//ここでアクセストークンをデコードしてアカウントIDを取得する accountId := ...
	//取得したworkspaceIdとaccountIdを使ってデータベースを参照しuserIdを取得する、できなかったら認証エラーを吐くようにする
	//他のユーザーを削除するため、ownerかmanagerのみが削除できるようにする
	userId := "dsafkljdfkslja"

	targetUser := model.User{}
	if err := dao.DeleteUser(&targetUser, userId).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}

func DeleteUser(c *gin.Context) {
	workspaceId := c.Request.Header.Get("workspace_id")
	//ここでアクセストークンをデコードしてアカウントIDを取得する accountId := ...
	//取得したworkspaceIdとaccountIdを使ってデータベースを参照しuserIdを取得する、できなかったら認証エラーを吐くようにする
	//自身を削除するために使う、ownerは自身を削除できないようにする
	userId := "dsafkljdfkslja"

	targetUser := model.User{}
	if err := dao.DeleteUser(&targetUser, userId).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}

func GrantRoleToUser(c *gin.Context) {
	workspaceId := c.Request.Header.Get("workspace_id")
	//ここでアクセストークンをデコードしてアカウントIDを取得する accountId := ...
	//取得したworkspaceIdとaccountIdを使ってデータベースを参照しuserIdを取得する、できなかったら認証エラーを吐くようにする
	//ownerかmanagerのみが変更できるようにする
	//managerかgeneralしか選べないようにする

	req := new(GrantRoleToUserReq)
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	if req.Role != "manager" && req.Role != "general" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "choose manager or general"})
	}

	targetUser := model.User{}
	if err := dao.GrantRoleToUser(&targetUser, req.UserId, req.Role).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, targetUser)
}

func ChangeUserAttributes(c *gin.Context) {
	
}
