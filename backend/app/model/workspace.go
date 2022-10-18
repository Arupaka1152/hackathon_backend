package model

import (
	"gorm.io/gorm"
	"time"
)

type Workspace struct {
	Id        string    `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"not null"`
	AvatarUrl string    `json:"avatar_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

func (p *Workspace) CreateWorkspace() (db *gorm.DB) {
	return db.Create(&p)
}

func (p *Workspace) ChangeWorkspaceAttributes(workspaceId string, avatarUrl string) (db *gorm.DB) {
	return db.Where("Id = ?", workspaceId).Update("avatarUrl", avatarUrl)
}

func (p *Workspace) DeleteWorkspace(workspaceId string) (db *gorm.DB) {
	return db.Where("Id = ?", workspaceId).Delete(&p)
}
