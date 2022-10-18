package model

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	Id          string         `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"not null"`
	AccountId   string         `json:"account_id" gorm:"not null"`
	WorkspaceId string         `json:"workspace_id" gorm:"not null"`
	Role        string         `json:"role" gorm:"not null"`
	AvatarUrl   string         `json:"avatar_url"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at"`
}

type Users []User

func (p *User) CreateUser() (db *gorm.DB) {
	return db.Create(&p)
}

func (u *Users) FetchAllUsersInWorkspace(workspaceId string) (db *gorm.DB) {
	return db.Where("workspaceId = ?", workspaceId).Find(&u)
}

func (p *User) DeleteUserFromWorkspace(userId string) (db *gorm.DB) {
	return db.Where("Id = ?", userId).Delete(&p)
}

func (p *User) GrantRoleToUser(userId string, role string) (db *gorm.DB) {
	return db.Where("Id = ?", userId).Update("role", role)
}

func (p *User) ChangeUserName(userId string, userName string) (db *gorm.DB) {
	return db.Where("Id = ?", userId).Update("name", userName)
}

func (p *User) ChangeUserAvatarUrl(userId string, avatarUrl string) (db *gorm.DB) {
	return db.Where("Id = ?", userId).Update("avatarUrl", avatarUrl)
}
