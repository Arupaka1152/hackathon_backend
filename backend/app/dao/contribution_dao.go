package dao

import (
	"backend/app/model"
	"gorm.io/gorm"
)

func FindContribution(p *model.Contribution, contributionId string) (tx *gorm.DB) {
	return db.Where("id = ?", contributionId).First(&p)
}

func CreateContribution(p *model.Contribution) (tx *gorm.DB) {
	return db.Create(&p)
}

func DeleteContribution(p *model.Contribution, contributionId string) (tx *gorm.DB) {
	return db.Where("id = ?", contributionId).Delete(&p)
}

func UpdateContribution(p *model.Contribution, contributionId string, points int, message string) (tx *gorm.DB) {
	return db.Model(&p).Where("id = ?", contributionId).Updates(model.Contribution{Message: message, Points: points})
}

func GetAllContributionInWorkspace(u *model.Contributions, workspaceId string) (tx *gorm.DB) {
	return db.Where("workspace_id = ?", workspaceId).Find(&u)
}

func GetAllContributionSent(u *model.Contributions, workspaceId string, from string) (tx *gorm.DB) {
	return db.Where(model.Contribution{WorkspaceId: workspaceId, From: from}).Find(&u)
}

func GetDesignatedContributionSent(u *model.Contributions, workspaceId string, from string, startDate string, endDate string) (tx *gorm.DB) {
	return db.Where(model.Contribution{WorkspaceId: workspaceId, From: from}).Where("created_at BETWEEN ? AND ?", startDate, endDate).Find(&u)
}

func GetAllContributionReceived(u *model.Contributions, workspaceId string, to string) (tx *gorm.DB) {
	return db.Where(model.Contribution{WorkspaceId: workspaceId, To: to}).Find(&u)
}

func GetDesignatedContributionReceived(u *model.Contributions, workspaceId string, to string, startDate string, endDate string) (tx *gorm.DB) {
	return db.Where(model.Contribution{WorkspaceId: workspaceId, To: to}).Where("created_at BETWEEN ? AND ?", startDate, endDate).Find(&u)
}

func UpdateReaction(p *model.Contribution, contributionId string, reactions int) (tx *gorm.DB) {
	return db.Model(&p).Where("id = ?", contributionId).Update("reaction", reactions)
}
