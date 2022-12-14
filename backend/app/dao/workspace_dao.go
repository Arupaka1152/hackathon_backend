package dao

import (
	"backend/app/model"
	"gorm.io/gorm"
)

func CreateWorkspace(p *model.Workspace) (tx *gorm.DB) {
	return db.Create(&p)
}

func UpdateWorkspaceAttributes(p *model.Workspace, workspaceId string, workspaceName string, description string, avatarUrl string) (tx *gorm.DB) {
	return db.Model(&p).Where("id = ?", workspaceId).Updates(model.Workspace{Name: workspaceName, Description: description, AvatarUrl: avatarUrl})
}

func DeleteWorkspace(p *model.Workspace, workspaceId string) (tx *gorm.DB) {
	return db.Where("id = ?", workspaceId).Delete(&p)
}

func FindWorkspaceInfo(p *model.Workspace, workspaceId string) (tx *gorm.DB) {
	return db.Where("id = ?", workspaceId).First(&p)
}
