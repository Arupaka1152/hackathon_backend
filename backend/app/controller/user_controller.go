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

type ChangeUserAttributesReq struct {
	UserName      string `json:"workspace_name"`
	UserAvatarUrl string `json:"workspace_avatar_url"`
}

func CreateUser(c *gin.Context) {
	workspaceId := c.Request.Header.Get("workspace_id")
	//ここでアクセストークンをデコードしてアカウントIDを取得する accountId := ...
	//取得したworkspaceIdとaccountIdを使ってデータベースを参照しuserIdを取得する、できなかったら認証エラーを吐くようにする
	//ownerかmanagerが招待できるようにする

	r := new(CreateUserReq)
	if err := c.Bind(&r); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	targetAccount := model.Account{}
	if err := dao.FetchAccountByEmail(&targetAccount, r.Email).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "account not found"})
	}

	userId := ulid.Make().String()
	newUser := model.User{
		Id:          userId,
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
	//ここでアクセストークンをデコードしてアカウントIDを取得する accountId := ...
	//取得したworkspaceIdとaccountIdを使ってデータベースを参照しuserIdを取得する、できなかったら認証エラーを吐くようにする
	userId := "fdksjlakflafdj"

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
