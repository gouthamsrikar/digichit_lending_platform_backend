package controllers

import (
	"chitfund/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateUser handles user creation
func FetchAllCommunities(c *gin.Context, db *gorm.DB) {
	userId := c.Param("user_id")
	var communities []models.Community
	if err := db.Find(&communities).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var joinedCommunities []models.JoinCommunity
	if err := db.Where("user_id = ?", userId).Find(&joinedCommunities).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"communties": communities, "communty_ledger": joinedCommunities})
}

func RequestToJoinCommunity(c *gin.Context, db *gorm.DB) {

	var request struct {
		CommunityID uint `json:"community_id"`
		UserID      uint `json:"user_id"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	JoinCommunity := models.JoinCommunity{
		CommunityID: request.CommunityID,
		UserID:      request.UserID,
		State:       "REQUESTED",
	}

	if err := db.Create(&JoinCommunity).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create loan entry"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"communties": JoinCommunity})
}

func ApproveToJoinCommunity(c *gin.Context, db *gorm.DB) {
	ledger_id := c.Param("ledger_id")

	err := db.Model(&models.JoinCommunity{}).Where("id = ?", ledger_id).Update("state", "APPROVED").Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "id does not exist"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"SUCCESS": "SUCCESS"})
}
