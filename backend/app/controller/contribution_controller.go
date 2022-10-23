package controller

import (
	"backend/app/auth"
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
	token := c.Request.Header.Get("authentication")

	accountId, err := auth.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	}

	userId, _, err := auth.UserAuth(workspaceId, accountId)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "not permitted"})
	}

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
	token := c.Request.Header.Get("authentication")

	accountId, err := auth.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	}

	userId, _, err := auth.UserAuth(workspaceId, accountId)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "not permitted"})
	}

	r := new(DeleteContributionReq)
	if err := c.Bind(&r); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	if err = auth.ContributionAuth(r.ContributionId, userId); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "not permitted"})
	}

	targetContribution := model.Contribution{}
	if err := dao.DeleteContribution(&targetContribution, r.ContributionId).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"message": "contribution deleted"})
}

func EditContribution(c *gin.Context) {
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

	r := new(EditContributionReq)
	if err := c.Bind(&r); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	if err = auth.ContributionAuth(r.ContributionId, userId); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "not permitted"})
	}

	targetContribution := model.Contribution{}
	if err := dao.EditContribution(&targetContribution, r.ContributionId, r.Points, r.Message).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, targetContribution)
}

func FetchAllContributionInWorkspace(c *gin.Context) {
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

	targetContributions := model.Contributions{}
	if err := dao.FetchAllContributionInWorkspace(&targetContributions, workspaceId).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, targetContributions)
}

func FetchAllContributionSent(c *gin.Context) {
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

	targetContributions := model.Contributions{}
	if err := dao.FetchAllContributionSent(&targetContributions, workspaceId, userId).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, targetContributions)
}

func FetchAllContributionReceived(c *gin.Context) {
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
	token := c.Request.Header.Get("authentication")

	accountId, err := auth.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	}

	_, _, err = auth.UserAuth(workspaceId, accountId)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "not permitted"})
	}

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
