package config

import (
	"chitfund/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	var err error
	DB, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Drop existing tables
	err = DB.Migrator().DropTable(&models.User{}, &models.MonthAnalysis{}, &models.Community{}, &models.LoanLedger{}, &models.JoinCommunity{}, &models.BankStatement{})
	if err != nil {
		log.Fatal("Failed to drop tables:", err)
	}

	err = DB.AutoMigrate(&models.User{}, &models.MonthAnalysis{}, &models.Community{}, &models.LoanLedger{}, &models.JoinCommunity{}, &models.BankStatement{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	return DB
}
