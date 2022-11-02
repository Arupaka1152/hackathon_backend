package model

import (
	"gorm.io/gorm"
	"time"
)

type Contribution struct {
	Id          string         `json:"id" gorm:"primaryKey"`
	WorkspaceId string         `json:"workspace_id" gorm:"not null"`
	From        string         `json:"from" gorm:"not null"`
	To          string         `json:"to" gorm:"not null"`
	Points      int            `json:"points" gorm:"not null"`
	Message     string         `json:"message"`
	Reaction    int            `json:"reaction"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at"`
}

type Contributions []Contribution
