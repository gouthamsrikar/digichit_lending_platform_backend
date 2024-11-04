package models

import "gorm.io/gorm"

type LoanLedger struct {
	gorm.Model
	ID               uint    `gorm:"primaryKey"`
	LoanAmount       float64 `json:"loan_amount"`
	CommunityID      uint    `json:"community_id"`
	LoanAmountRepaid float64 `json:"loan_amount_repaid"`
	UserID           uint    `json:"user_id"` // New field for user ID
	LedgerState      string  `json:"ledger_state"`
}
