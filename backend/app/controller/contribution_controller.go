package controller

import (
	"backend/app/auth"
	"backend/app/dao"
	"backend/app/model"
	"backend/app/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
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
	Id       string `json:"contribution_id"`
	To       string `json:"receiver_id"`
	Points   int    `json:"points"`
	Message  string `json:"message"`
	UpdateAt string `json:"update_at"`
}

type SendReactionRes struct {
	Id       string `json:"contribution_id"`
	Reaction int    `json:"reaction"`
	UpdateAt string `json:"update_at"`
}

type ContributionReport struct {
	UserId               string `json:"user_id"`
	Name                 string `json:"name"`
	AvatarUrl            string `json:"avatar_url"`
	ContributionSent     int    `json:"contribution_sent"`
	PointsSent           int    `json:"points_sent"`
	ReactionSent         int    `json:"reaction_sent"`
	ContributionReceived int    `json:"contribution_received"`
	PointsReceived       int    `json:"points_received"`
	ReactionReceived     int    `json:"reaction_received"`
}

type ContributionReportRes []ContributionReport

func FetchContributionReport(c *gin.Context) {
	workspaceId := utils.GetValueFromContext(c, "workspaceId")

	endDate := time.Now()
	startDate := endDate.Add(-7 * 24 * time.Hour)

	targetContributions := model.Contributions{}
	if err := dao.GetDesignatedContributionInWorkspace(&targetContributions, workspaceId, startDate, endDate).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	targetUsers := model.Users{}
	if err := dao.GetAllUsersInWorkspace(&targetUsers, workspaceId).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	res := make(ContributionReportRes, 0)

	for i := 0; i < len(targetUsers); i++ {
		contributionSent := 0
		pointsSent := 0
		reactionSent := 0
		contributionReceived := 0
		pointsReceived := 0
		reactionReceived := 0
		for j := 0; j < len(targetContributions); j++ {
			if targetUsers[i].Id == targetContributions[j].From {
				contributionSent++
				pointsSent += targetContributions[j].Points
				reactionSent += targetContributions[j].Reaction
			}

			if targetUsers[i].Id == targetContributions[j].To {
				contributionReceived++
				pointsReceived += targetContributions[j].Points
				reactionReceived += targetContributions[j].Reaction
			}
		}
		res = append(res, ContributionReport{
			targetUsers[i].Id,
			targetUsers[i].Name,
			targetUsers[i].AvatarUrl,
			contributionSent,
			pointsSent,
			reactionSent,
			contributionReceived,
			pointsReceived,
			reactionReceived,
		})
	}

	c.JSON(http.StatusOK, res)
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

	if req.Points < 1 && req.Points > 100 {
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

	jst := time.FixedZone("Asia/Tokyo", 9*60*60)

	res := &ContributionRes{
		newContribution.Id,
		newContribution.WorkspaceId,
		newContribution.From,
		newContribution.To,
		newContribution.Points,
		newContribution.Message,
		newContribution.Reaction,
		newContribution.CreatedAt.In(jst).Format(layout),
		newContribution.UpdatedAt.In(jst).Format(layout),
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

	jst := time.FixedZone("Asia/Tokyo", 9*60*60)

	res := &EditContributionRes{
		req.ContributionId,
		targetContribution.To,
		targetContribution.Points,
		targetContribution.Message,
		targetContribution.UpdatedAt.In(jst).Format(layout),
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

	jst := time.FixedZone("Asia/Tokyo", 9*60*60)

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
			targetContributions[i].CreatedAt.In(jst).Format(layout),
			targetContributions[i].UpdatedAt.In(jst).Format(layout),
		})
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

	jst := time.FixedZone("Asia/Tokyo", 9*60*60)

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
			targetContributions[i].CreatedAt.In(jst).Format(layout),
			targetContributions[i].UpdatedAt.In(jst).Format(layout),
		})
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

	jst := time.FixedZone("Asia/Tokyo", 9*60*60)

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
			targetContributions[i].CreatedAt.In(jst).Format(layout),
			targetContributions[i].UpdatedAt.In(jst).Format(layout),
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

	jst := time.FixedZone("Asia/Tokyo", 9*60*60)

	res := &SendReactionRes{
		req.ContributionId,
		newContribution.Reaction,
		targetContribution.UpdatedAt.In(jst).Format(layout),
	}

	c.JSON(http.StatusOK, res)
}
