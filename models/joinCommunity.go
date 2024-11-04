package models

import "gorm.io/gorm"

type JoinCommunity struct {
	gorm.Model
	ID          uint   `gorm:"primaryKey"`
	CommunityID uint   `json:"community_id"`
	UserID      uint   `json:"user_id"`
	State       string `json:"state"`
}
