package dao

import (
	"backend/app/model"
	"gorm.io/gorm"
)

func CreateUser(p *model.User) (db *gorm.DB) {
	return db.Create(&p)
}

func FetchAllUsersInWorkspace(u *model.Users, workspaceId string) (db *gorm.DB) {
	return db.Where("workspaceId = ?", workspaceId).Find(&u)
}

func DeleteUserFromWorkspace(p *model.User, userId string) (db *gorm.DB) {
	return db.Where("id = ?", userId).Delete(&p)
}

func GrantRoleToUser(p *model.User, userId string, role string) (db *gorm.DB) {
	return db.Model(&p).Where("id = ?", userId).Update("role", role)
}

func ChangeUserAttributes(p *model.User, userId string, userName string, avatarUrl string) (db *gorm.DB) {
	return db.Model(&p).Where("id = ?", userId).Updates(model.User{Name: userName, AvatarUrl: avatarUrl})
}
