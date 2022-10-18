package model

import (
	"gorm.io/gorm"
	"time"
)

type Account struct {
	Id        string    `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"not null"`
	Email     string    `json:"email" gorm:"unique;not null"`
	Password  string    `json:"password" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}

func (p *Account) CreateAccount() (db *gorm.DB) {
	return db.Create(&p)
}
