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
