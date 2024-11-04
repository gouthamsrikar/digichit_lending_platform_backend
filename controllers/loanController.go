package controllers

import (
	"chitfund/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RequestLoan(c *gin.Context, db *gorm.DB) {
	var request struct {
		LoanAmount float64 `json:"loan_amount"`
		ID         uint    `json:"community_id"`
		UserID     uint    `json:"user_id"` // Accept user ID in the request
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var community models.Community
	if err := db.First(&community, request.ID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Community not found"})
		return
	}

	if community.TotalFund < request.LoanAmount {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient community fund"})
		return
	}

	loanLedger := models.LoanLedger{
		LoanAmount:       request.LoanAmount,
		CommunityID:      request.ID,
		LoanAmountRepaid: 0.0,
		UserID:           request.UserID, // Set user ID
		LedgerState:      "APPLIED",
	}
	if err := db.Create(&loanLedger).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create loan entry"})
		return
	}

	c.JSON(http.StatusOK, loanLedger)
}

func GetLoanLedgersByUserID(c *gin.Context, db *gorm.DB) {
	userID := c.Param("user_id") // Get the user ID from the URL parameter

	var loanLedgers []models.LoanLedger
	if err := db.Where("user_id = ?", userID).Find(&loanLedgers).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No loan ledgers found for this user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"loans": loanLedgers})
}

func ApproveLoan(c *gin.Context, db *gorm.DB) {
	loanID := c.Param("loan_id")

	var loanLedger models.LoanLedger
	if err := db.First(&loanLedger, loanID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Loan ledger not found"})
		return
	}

	loanLedger.LedgerState = "APPROVED"

	err := db.Model(&models.LoanLedger{}).Where("id = ?", loanLedger.ID).Update("ledger_state", "APPROVED").Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "id does not exist"})
		return
	}

	if err := db.First(&loanLedger, loanID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Loan ledger not found"})
		return
	}

	var community models.Community
	if err := db.First(&community, loanLedger.CommunityID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Community not found"})
		return
	}

	if community.TotalFund < loanLedger.LoanAmount {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient community fund"})
		return
	}

	community.TotalFund -= loanLedger.LoanAmount

	if err := db.Save(&community).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update community fund"})
		return
	}

	if err := db.Save(&loanLedger).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update loan ledger"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Loan approved", "loan": loanLedger, "community": community})
}
