package controllers

import (
	"chitfund/models"
	"gorm.io/gorm"
)

func GetBankStatementsByUserID(db *gorm.DB, userID uint) ([]models.BankStatement, error) {
	var statements []models.BankStatement
	if err := db.Where("user_id = ?", userID).Find(&statements).Error; err != nil {
		return nil, err
	}
	return statements, nil
}

func CreateBankStatements(db *gorm.DB, statements []models.BankStatement) error {
	if err := db.Create(&statements).Error; err != nil {
		return err
	}
	return nil
}
