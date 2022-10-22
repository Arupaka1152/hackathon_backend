package dao

import (
	"backend/app/model"
	"gorm.io/gorm"
)

func CreateAccount(p *model.Account) (db *gorm.DB) {
	return db.Create(&p)
}

func FetchAccountByEmail(p *model.Account, email string) (db *gorm.DB) {
	return db.Where("email = ?", email).Find(&p)
}
