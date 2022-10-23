package controller

import (
	"backend/app/dao"
	"backend/app/model"
	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"
	"net/http"
)

type CreateContributionReq struct {
	ReceiverId string `json:"receiver_id"`
	Points     int    `json:"points"`
	Message    string `json:"message"`
}

type DeleteContributionReq struct {
	ContributionId string `json:"contribution_id"`
}

type EditContributionReq struct {
	ContributionId string `json:"contribution_id"`
	Points         int    `json:"points"`
	Message        string `json:"message"`
}

type FetchAllContributionSentReq struct {
	ReceiverId string `json:"receiver_id"`
}

type SendReactionReq struct {
	ContributionId string `json:"contribution_id"`
}

func CreateContribution(c *gin.Context) {
	workspaceId := c.Request.Header.Get("workspace_id")
	//ここでアクセストークンをデコードしてアカウントIDを取得する accountId := ...
	//取得したworkspaceIdとaccountIdを使ってデータベースを参照しuserIdを取得する、できなかったら認証エラーを吐くようにする
	userId := "faksdjlkajlsdfkljfads"

	r := new(CreateContributionReq)
	if err := c.Bind(&r); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	contributionId := ulid.Make().String()
	newContribution := model.Contribution{
		Id:          contributionId,
		WorkspaceId: workspaceId,
		From:        userId,
		To:          r.ReceiverId,
		Points:      r.Points,
		Message:     r.Message,
		Reaction:    0,
	}

	if err := dao.CreateContribution(&newContribution).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, newContribution)
}

func DeleteContribution(c *gin.Context) {
	workspaceId := c.Request.Header.Get("workspace_id")
	//ここでアクセストークンをデコードしてアカウントIDを取得する accountId := ...
	//取得したworkspaceIdとaccountIdを使ってデータベースを参照しuserIdを取得する、できなかったら認証エラーを吐く
	//もしContributionのFromとuserIdが異なれば認可エラーを吐く

	r := new(DeleteContributionReq)
	if err := c.Bind(&r); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	targetContribution := model.Contribution{}

	if err := dao.DeleteContribution(&targetContribution, r.ContributionId).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"message": "contribution deleted"})
}

func EditContribution(c *gin.Context) {
	workspaceId := c.Request.Header.Get("workspace_id")
	//ここでアクセストークンをデコードしてアカウントIDを取得する accountId := ...
	//取得したworkspaceIdとaccountIdを使ってデータベースを参照しuserIdを取得する、できなかったら認証エラーを吐く
	//もしContributionのFromとuserIdが異なれば認可エラーを吐く

	r := new(EditContributionReq)
	if err := c.Bind(&r); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	targetContribution := model.Contribution{}

	if err := dao.EditContribution(&targetContribution, r.ContributionId, r.Points, r.Message).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, targetContribution)
}

func FetchAllContributionInWorkspace(c *gin.Context) {
	workspaceId := c.Request.Header.Get("workspace_id")
	//ここでアクセストークンをデコードしてアカウントIDを取得する accountId := ...
	//取得したworkspaceIdとaccountIdを使ってデータベースを参照しuserIdを取得する、できなかったら認証エラーを吐く

	targetContributions := model.Contributions{}

	if err := dao.FetchAllContributionInWorkspace(&targetContributions, workspaceId).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, targetContributions)
}

func FetchAllContributionSent(c *gin.Context) {
	workspaceId := c.Request.Header.Get("workspace_id")
	//ここでアクセストークンをデコードしてアカウントIDを取得する accountId := ...
	//取得したworkspaceIdとaccountIdを使ってデータベースを参照しuserIdを取得する、できなかったら認証エラーを吐く
	userId := "akfjdlksadjlfkajlfsd"

	targetContributions := model.Contributions{}

	if err := dao.FetchAllContributionSent(&targetContributions, workspaceId, userId).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, targetContributions)
}

func FetchAllContributionReceived(c *gin.Context) {
	workspaceId := c.Request.Header.Get("workspace_id")
	//ここでアクセストークンをデコードしてアカウントIDを取得する accountId := ...
	//取得したworkspaceIdとaccountIdを使ってデータベースを参照しuserIdを取得する、できなかったら認証エラーを吐く

	r := new(FetchAllContributionSentReq)
	if err := c.Bind(&r); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	targetContributions := model.Contributions{}

	if err := dao.FetchAllContributionReceived(&targetContributions, workspaceId, r.ReceiverId).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, targetContributions)
}

func SendReaction(c *gin.Context) {
	workspaceId := c.Request.Header.Get("workspace_id")
	//ここでアクセストークンをデコードしてアカウントIDを取得する accountId := ...
	//取得したworkspaceIdとaccountIdを使ってデータベースを参照しuserIdを取得する、できなかったら認証エラーを吐く

	r := new(SendReactionReq)
	if err := c.Bind(&r); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	targetContribution := model.Contribution{}

	if err := dao.FetchContribution(&targetContribution, r.ContributionId).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	totalReaction := targetContribution.Reaction + 1
	newContribution := model.Contribution{}

	if err := dao.SendReaction(&newContribution, r.ContributionId, totalReaction).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, newContribution)
}
