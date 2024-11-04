package models

import "gorm.io/gorm"

// User struct

type Community struct {
	gorm.Model
	ID                      uint    `gorm:"primaryKey"`
	CommunityName           string  `json:"community_name"`
	MonthlyDeposit          float64 `json:"monthly_deposit"`
	TotalFund               float64 `json:"total_fund"`
	InterestRate            float64 `json:"interest_rate"`
	AdminName               string  `json:"admin_name"`
	AdminUserID             string  `json:"admin_user_id"`
	RepaymentPeriodInMonths int     `json:"repayment_period_in_months"`
	CommunityDescription    string  `json:"community_description"`
	UserCount               int     `json:"user_count"`
	MaxCount                int     `json:"max_count"`
}
