package controller

import (
	"backend/app/auth"
	"backend/app/dao"
	"backend/app/model"
	"backend/app/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

const layout = "2006-01-02 15:04:05"

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
	CreatedAt   string `json:"created_at"`
	UpdateAt    string `json:"update_at"`
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

	if userId == req.ReceiverId {
		c.JSON(http.StatusBadRequest, gin.H{"error": "you cant send contribution to yourself"})
		return
	}

	if req.Points > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "you cant send more than 100 points"})
		return
	}

	contributionId := utils.GenerateId()
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
		newContribution.CreatedAt.Format(layout),
		newContribution.UpdatedAt.Format(layout),
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

	if req.Points > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "you cant send more than 100 points"})
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

	startDate := c.Param("startDate")
	endDate := c.Param("endDate")

	targetContributions := model.Contributions{}

	if startDate != "" && endDate != "" {
		if err := dao.GetDesignatedContributionInWorkspace(&targetContributions, workspaceId, startDate, endDate).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if targetContributions[0].Id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "contributions not found"})
			return
		}
	} else {
		if err := dao.GetAllContributionInWorkspace(&targetContributions, workspaceId).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if targetContributions[0].Id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "contributions not found"})
			return
		}
	}

	res := make(ContributionsRes, 0)
	for i := 0; i < len(targetContributions); i++ {
		res = append(res, ContributionRes{
			targetContributions[i].Id,
			targetContributions[i].WorkspaceId,
			targetContributions[i].From,
			targetContributions[i].To,
			targetContributions[i].Points,
			targetContributions[i].Message,
			targetContributions[i].Reaction,
			targetContributions[i].CreatedAt.Format(layout),
			targetContributions[i].UpdatedAt.Format(layout),
		})
	}

	c.JSON(http.StatusOK, res)
}

func FetchAllContributionSent(c *gin.Context) {
	workspaceId := utils.GetValueFromContext(c, "workspaceId")
	userId := utils.GetValueFromContext(c, "userId")

	startDate := c.Param("startDate")
	endDate := c.Param("endDate")

	targetContributions := model.Contributions{}

	if startDate != "" && endDate != "" {
		if err := dao.GetDesignatedContributionSent(&targetContributions, workspaceId, userId, startDate, endDate).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if targetContributions[0].Id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "contributions not found"})
			return
		}
	} else {
		if err := dao.GetAllContributionSent(&targetContributions, workspaceId, userId).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if targetContributions[0].Id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "contributions not found"})
			return
		}
	}

	res := make(ContributionsRes, 0)
	for i := 0; i < len(targetContributions); i++ {
		res = append(res, ContributionRes{
			targetContributions[i].Id,
			targetContributions[i].WorkspaceId,
			targetContributions[i].From,
			targetContributions[i].To,
			targetContributions[i].Points,
			targetContributions[i].Message,
			targetContributions[i].Reaction,
			targetContributions[i].CreatedAt.Format(layout),
			targetContributions[i].UpdatedAt.Format(layout),
		})
	}

	c.JSON(http.StatusOK, res)
}

func FetchAllContributionReceived(c *gin.Context) {
	workspaceId := utils.GetValueFromContext(c, "workspaceId")
	userId := utils.GetValueFromContext(c, "userId")

	startDate := c.Param("startDate")
	endDate := c.Param("endDate")

	targetContributions := model.Contributions{}

	if startDate != "" && endDate != "" {
		if err := dao.GetDesignatedContributionReceived(&targetContributions, workspaceId, userId, startDate, endDate).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if targetContributions[0].Id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "contributions not found"})
			return
		}
	} else {
		if err := dao.GetAllContributionReceived(&targetContributions, workspaceId, userId).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if targetContributions[0].Id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "contributions not found"})
			return
		}
	}

	res := make(ContributionsRes, 0)
	for i := 0; i < len(targetContributions); i++ {
		res = append(res, ContributionRes{
			targetContributions[i].Id,
			targetContributions[i].WorkspaceId,
			targetContributions[i].From,
			targetContributions[i].To,
			targetContributions[i].Points,
			targetContributions[i].Message,
			targetContributions[i].Reaction,
			targetContributions[i].CreatedAt.Format(layout),
			targetContributions[i].UpdatedAt.Format(layout),
		})
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
