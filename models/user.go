package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID                 uint    `gorm:"primaryKey"`
	Name               string  `json:"name"`
	JoinDate           string  `json:"join_date"`
	PanNo              string  `json:"pan_no"`
	BankAccountNo      string  `json:"bank_account_no"`
	CommunityLevel     string  `json:"community_level"`
	DigitScore         int     `json:"digit_score"`
	MonthlyIncome      float64 `json:"monthly_income"`
	AvgBalance         float64 `json:"avg_balance"`
	EmiToIncomeRatio   float64 `json:"emi_to_income_ratio"`
	CreditToDebitRatio float64 `json:"credit_to_debit_ratio"`
	PhoneNumber        string  `json:"phone_number"`
}

type MonthAnalysis struct {
	gorm.Model
	UserID int    `json:"user_id"`
	Month  string `json:"month"`
	Data   string `json:"data"`
}
