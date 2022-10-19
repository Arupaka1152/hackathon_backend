package dao

import (
	"backend/app/model"
	"gorm.io/gorm"
)

func CreateAccount(p *model.Account) (db *gorm.DB) {
	return db.Create(&p)
}
