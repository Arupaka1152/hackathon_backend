package model

import "time"

type User struct {
	Id          string    `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null"`
	AccountId   string    `json:"account_id" gorm:"not null"`
	WorkspaceId string    `json:"workspace_id" gorm:"not null"`
	Role        string    `json:"role" gorm:"not null"`
	PictureUrl  string    `json:"picture_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   time.Time `json:"deleted_at"`
}
