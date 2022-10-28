package dao

import (
	"backend/app/model"
	"gorm.io/gorm"
)

func CreateAccount(p *model.Account) (tx *gorm.DB) {
	return db.Create(&p)
}

func FetchAccountByEmail(p *model.Account, email string) (tx *gorm.DB) {
	return db.Where("email = ?", email).Find(&p)
}
