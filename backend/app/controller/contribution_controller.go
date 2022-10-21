package controller

import "C"
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

func CreateContribution(c *gin.Context) {
	workspaceId := c.Request.Header.Get("workspace_id")
	//ここでアクセストークンをデコードしてアカウントIDを取得する accountId := ...
	//取得したworkspaceIdとaccountIdを使ってデータベースを参照しuserIdを取得する、できなかったら認証エラーを吐くようにする
	userId := "John Smith"

	req := new(CreateContributionReq)
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	contributionId := ulid.Make().String()

	newContribution := model.Contribution{
		Id:          contributionId,
		WorkspaceId: workspaceId,
		From:        userId,
		To:          req.ReceiverId,
		Points:      req.Points,
		Message:     req.Message,
		Reaction:    0,
	}

	if err := dao.CreateContribution(&newContribution).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"contribution_id": contributionId})
}

func DeleteContribution(c *gin.Context) {
	workspaceId := c.Request.Header.Get("workspace_id")
	//ここでアクセストークンをデコードしてアカウントIDを取得する accountId := ...
	//取得したworkspaceIdとaccountIdを使ってデータベースを参照しuserIdを取得する、できなかったら認証エラーを吐く
	//もしContributionのFromとuserIdが異なれば認可エラーを吐く

	req := new(DeleteContributionReq)
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	if err := dao.CreateContribution(&newContribution).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

func EditContribution(c *gin.Context) {

}

func FetchAllContributionInWorkspace(c *gin.Context) {

}

func FetchAllContributionSent(c *gin.Context) {

}

func FetchAllContributionReceived(c *gin.Context) {

}

func SendReaction(c *gin.Context) {

}
