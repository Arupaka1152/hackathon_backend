package dao

import (
	"backend/app/model"
	"gorm.io/gorm"
)

func FetchContribution(p *model.Contribution, contributionId string) (db *gorm.DB) {
	return db.Where("id = ?", contributionId).Find(&p)
}

func CreateContribution(p *model.Contribution) (db *gorm.DB) {
	return db.Create(&p)
}

func DeleteContribution(p *model.Contribution, contributionId string) (db *gorm.DB) {
	return db.Where("id = ?", contributionId).Delete(&p)
}

func EditContribution(p *model.Contribution, contributionId string, points int, message string) (db *gorm.DB) {
	return db.Model(&p).Where("id = ?", contributionId).Updates(model.Contribution{Message: message, Points: points})
}

func FetchAllContributionInWorkspace(u *model.Contributions, workspaceId string) (db *gorm.DB) {
	return db.Where("workspace_id = ?", workspaceId).Find(&u)
}

func FetchAllContributionSent(u *model.Contributions, workspaceId string, from string) (db *gorm.DB) {
	return db.Where(model.Contribution{WorkspaceId: workspaceId, From: from}).Find(&u)
}

func FetchAllContributionReceived(u *model.Contributions, workspaceId string, to string) (db *gorm.DB) {
	return db.Where(model.Contribution{WorkspaceId: workspaceId, To: to}).Find(&u)
}

func SendReaction(p *model.Contribution, contributionId string, reactions int) (db *gorm.DB) {
	return db.Model(&p).Where("id = ?", contributionId).Update("reaction", reactions)
}
