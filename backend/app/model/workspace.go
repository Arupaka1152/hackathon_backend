package model

import "time"

type Workspace struct {
	Id         string    `json:"id" gorm:"primaryKey"`
	Name       string    `json:"name" gorm:"not null"`
	PictureUrl string    `json:"picture_url"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	DeletedAt  time.Time `json:"deleted_at"`
}
