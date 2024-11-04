package services

import (
	"chitfund/models"
	"gorm.io/gorm"
)

func CreateUser(db *gorm.DB, user *models.User) error {
	return db.Create(user).Error
}

func AddMonthAnalysis(db *gorm.DB, analysis *models.MonthAnalysis) error {
	return db.Create(analysis).Error
}


