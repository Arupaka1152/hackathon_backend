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

type CreateContributionReq struct {
	ReceiverId string `json:"receiver_id" binding:"required"`
	Points     int    `json:"points" binding:"required"`
	Message    string `json:"message" binding:"required"`
}

type DeleteContributionReq struct {
	ContributionId string `json:"contribution_id" binding:"required"`
}

type EditContributionReq struct {
	ContributionId string `json:"contribution_id" binding:"required"`
	Points         int    `json:"points" binding:"required"`
	Message        string `json:"message" binding:"required"`
}

type SendReactionReq struct {
	ContributionId string `json:"contribution_id" binding:"required"`
}

type ContributionRes struct {
	Id          string `json:"contribution_id"`
	WorkspaceId string `json:"workspace_id"`
	From        string `json:"sender_id"`
	To          string `json:"receiver_id"`
	Points      int    `json:"points"`
	Message     string `json:"message"`
	Reaction    int    `json:"reaction"`
}

type ContributionsRes []ContributionRes

type EditContributionRes struct {
	Id      string `json:"contribution_id"`
	To      string `json:"receiver_id"`
	Points  int    `json:"points"`
	Message string `json:"message"`
}

type SendReactionRes struct {
	Id       string `json:"contribution_id"`
	Reaction int    `json:"reaction"`
}

func CreateContribution(c *gin.Context) {
	req := new(CreateContributionReq)
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	workspaceId := utils.GetValueFromContext(c, "workspaceId")
	userId := utils.GetValueFromContext(c, "userId")

	if userId != req.ReceiverId {
		c.JSON(http.StatusBadRequest, gin.H{"error": "you cant send contribution to yourself"})
		return
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
		return
	}

	res := &ContributionRes{
		newContribution.Id,
		newContribution.WorkspaceId,
		newContribution.From,
		newContribution.To,
		newContribution.Points,
		newContribution.Message,
		newContribution.Reaction,
	}

	c.JSON(http.StatusOK, res)
}

func DeleteContribution(c *gin.Context) {
	req := new(DeleteContributionReq)
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	userId := utils.GetValueFromContext(c, "userId")
	if err := auth.ContributionAuth(req.ContributionId, userId); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "not permitted"})
		return
	}

	targetContribution := model.Contribution{}
	if err := dao.DeleteContribution(&targetContribution, req.ContributionId).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "contribution deleted"})
}

func EditContribution(c *gin.Context) {
	req := new(EditContributionReq)
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	userId := utils.GetValueFromContext(c, "userId")
	if err := auth.ContributionAuth(req.ContributionId, userId); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "not permitted"})
		return
	}

	targetContribution := model.Contribution{}
	if err := dao.UpdateContribution(&targetContribution, req.ContributionId, req.Points, req.Message).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	res := &EditContributionRes{
		req.ContributionId,
		targetContribution.To,
		targetContribution.Points,
		targetContribution.Message,
	}

	c.JSON(http.StatusOK, res)
}

func FetchAllContributionInWorkspace(c *gin.Context) {
	workspaceId := utils.GetValueFromContext(c, "workspaceId")

	targetContributions := model.Contributions{}
	if err := dao.GetAllContributionInWorkspace(&targetContributions, workspaceId).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if targetContributions[0].Id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "contributions not found"})
		return
	}

	res := make(ContributionsRes, 0)
	for i := 0; i < len(targetContributions); i++ {
		res[i].Id = targetContributions[i].Id
		res[i].WorkspaceId = targetContributions[i].WorkspaceId
		res[i].From = targetContributions[i].From
		res[i].To = targetContributions[i].To
		res[i].Points = targetContributions[i].Points
		res[i].Message = targetContributions[i].Message
		res[i].Reaction = targetContributions[i].Reaction
	}

	c.JSON(http.StatusOK, res)
}

func FetchAllContributionSent(c *gin.Context) {
	workspaceId := utils.GetValueFromContext(c, "workspaceId")
	userId := utils.GetValueFromContext(c, "userId")

	targetContributions := model.Contributions{}
	if err := dao.GetAllContributionSent(&targetContributions, workspaceId, userId).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if targetContributions[0].Id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "contributions not found"})
		return
	}

	res := make(ContributionsRes, 0)
	for i := 0; i < len(targetContributions); i++ {
		res[i].Id = targetContributions[i].Id
		res[i].WorkspaceId = targetContributions[i].WorkspaceId
		res[i].From = targetContributions[i].From
		res[i].To = targetContributions[i].To
		res[i].Points = targetContributions[i].Points
		res[i].Message = targetContributions[i].Message
		res[i].Reaction = targetContributions[i].Reaction
	}

	c.JSON(http.StatusOK, res)
}

func FetchAllContributionReceived(c *gin.Context) {
	workspaceId := utils.GetValueFromContext(c, "workspaceId")
	userId := utils.GetValueFromContext(c, "userId")

	targetContributions := model.Contributions{}
	if err := dao.GetAllContributionReceived(&targetContributions, workspaceId, userId).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if targetContributions[0].Id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "contributions not found"})
		return
	}

	res := make(ContributionsRes, 0)
	for i := 0; i < len(targetContributions); i++ {
		res[i].Id = targetContributions[i].Id
		res[i].WorkspaceId = targetContributions[i].WorkspaceId
		res[i].From = targetContributions[i].From
		res[i].To = targetContributions[i].To
		res[i].Points = targetContributions[i].Points
		res[i].Message = targetContributions[i].Message
		res[i].Reaction = targetContributions[i].Reaction
	}

	c.JSON(http.StatusOK, res)
}

func SendReaction(c *gin.Context) {
	req := new(SendReactionReq)
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	targetContribution := model.Contribution{}
	if err := dao.FindContribution(&targetContribution, req.ContributionId).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalReaction := targetContribution.Reaction + 1
	newContribution := model.Contribution{}
	if err := dao.UpdateReaction(&newContribution, req.ContributionId, totalReaction).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	res := &SendReactionRes{
		req.ContributionId,
		newContribution.Reaction,
	}

	c.JSON(http.StatusOK, res)
}
