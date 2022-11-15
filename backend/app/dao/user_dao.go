package dao

import (
	"backend/app/model"
	"gorm.io/gorm"
)

func CreateUser(p *model.User) (tx *gorm.DB) {
	return db.Create(&p)
}

func GetAllUsersInWorkspace(u *model.Users, workspaceId string) (tx *gorm.DB) {
	return db.Where("workspace_id = ?", workspaceId).Find(&u)
}

func DeleteUser(p *model.User, userId string) (tx *gorm.DB) {
	return db.Where("id = ?", userId).Delete(&p)
}

func DeleteAllUsersInWorkspace(u *model.Users, workspaceId string) (tx *gorm.DB) {
	return db.Where("workspace_id = ?", workspaceId).Delete(&u)
}

func UpdateRoleOfUser(p *model.User, userId string, role string) (tx *gorm.DB) {
	return db.Model(&p).Where("id = ?", userId).Update("role", role)
}

func UpdateUserAttributes(p *model.User, userId string, userName string, description string, avatarUrl string) (tx *gorm.DB) {
	return db.Model(&p).Where("id = ?", userId).Updates(model.User{Name: userName, Description: description, AvatarUrl: avatarUrl})
}

func GetAllUsers(u *model.Users, accountId string) (tx *gorm.DB) {
	return db.Where("account_id = ?", accountId).Find(&u)
}

func FindUserById(p *model.User, workspaceId string, accountId string) (tx *gorm.DB) {
	return db.Where("workspace_id = ? AND account_id = ?", workspaceId, accountId).First(&p)
}

func FindUserByUserId(p *model.User, userId string) (tx *gorm.DB) {
	return db.Where("id = ?", userId).First(&p)
}
