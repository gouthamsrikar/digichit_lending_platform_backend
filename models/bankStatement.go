package models

import "gorm.io/gorm"

type BankStatement struct {
	gorm.Model
	ID                     uint    `gorm:"primaryKey"`
	UserID                 uint     `json:"user_id"`
	MonthName              string  `json:"month_name"`
	NoOfDebitTransactions  int     `json:"no_of_debit_transactions"`
	NoOfCreditTransactions int     `json:"no_of_credit_transactions"`
	TotalCreditAmount      float64 `json:"total_credit_amount"`
	Year                   string  `json:"year"`
	TotalDebitAmount       float64 `json:"total_debit_amount"`
	AverageEodBalance      float64 `json:"average_eod_balance"`
}
