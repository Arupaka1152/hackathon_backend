package dao

import (
	"backend/app/model"
	"gorm.io/gorm"
)

func CreateWorkspace(p *model.Workspace) (db *gorm.DB) {
	return db.Create(&p)
}

func ChangeWorkspaceAttributes(p *model.Workspace, workspaceId string, workspaceName string, avatarUrl string) (db *gorm.DB) {
	return db.Model(&p).Where("id = ?", workspaceId).Updates(model.Workspace{Name: workspaceName, AvatarUrl: avatarUrl})
}

func DeleteWorkspace(p *model.Workspace, workspaceId string) (db *gorm.DB) {
	return db.Where("id = ?", workspaceId).Delete(&p)
}

func FetchWorkspaceInfo(p *model.Workspace, workspaceId string) (db *gorm.DB) {
	return db.Where("id = ?", workspaceId).Find(&p)
}
