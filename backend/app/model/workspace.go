package model

import (
	"gorm.io/gorm"
	"time"
)

type Workspace struct {
	Id        string         `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"not null"`
	AvatarUrl string         `json:"avatar_url"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

type Workspaces []Workspace

func (p *Workspace) CreateWorkspace() (db *gorm.DB) {
	return db.Create(&p)
}

func (p *Workspace) ChangeWorkspaceAttributes(workspaceId string, workspaceName string, avatarUrl string) (db *gorm.DB) {
	return db.Model(&p).Where("id = ?", workspaceId).Updates(Workspace{Name: workspaceName, AvatarUrl: avatarUrl})
}

func (p *Workspace) DeleteWorkspace(workspaceId string) (db *gorm.DB) {
	return db.Where("id = ?", workspaceId).Delete(&p)
}

func (p *Workspace) FetchWorkspaceInfo(workspaceId string) (db *gorm.DB) {
	return db.Where("id = ?", workspaceId).Find(&p)
}
