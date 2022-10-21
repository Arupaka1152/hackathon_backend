package model

import (
	"gorm.io/gorm"
	"time"
)

type Contribution struct {
	Id          string    `json:"id" gorm:"primaryKey"`
	WorkspaceId string    `json:"workspace_id" gorm:"not null"`
	From        string    `json:"from" gorm:"not null"`
	To          string    `json:"to" gorm:"not null"`
	Points      int       `json:"points" gorm:"not null"`
	Message     string    `json:"message"`
	Reaction    int       `json:"reaction"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Contributions []Contribution

func (p *Contribution) CreateContribution() (db *gorm.DB) {
	return db.Create(&p)
}

func (p *Contribution) DeleteContribution(contributionId string) (db *gorm.DB) {
	return db.Where("id = ?", contributionId).Delete(&p)
}

func (p *Contribution) EditContribution(contributionId string, points int, message string) (db *gorm.DB) {
	return db.Model(&p).Where("id = ?", contributionId).Updates(Contribution{Message: message, Points: points})
}

func (u *Contributions) FetchAllContributionInWorkspace(workspaceId string) (db *gorm.DB) {
	return db.Where("workspace_id = ?", workspaceId).Find(&u)
}

func (u *Contributions) FetchAllContributionSent(workspaceId string, from string) (db *gorm.DB) {
	return db.Where(Contribution{WorkspaceId: workspaceId, From: from}).Find(&u)
}

func (u *Contributions) FetchAllContributionReceived(workspaceId string, to string) (db *gorm.DB) {
	return db.Where(Contribution{WorkspaceId: workspaceId, To: to}).Find(&u)
}

func (p *Contribution) SendReaction(contributionId string, reactions int) (db *gorm.DB) {
	return db.Model(&p).Where("id = ?", contributionId).Update("reaction", reactions)
}
