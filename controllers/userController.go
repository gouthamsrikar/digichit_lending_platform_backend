package controllers

import (
	"chitfund/models"
	"chitfund/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateUser(c *gin.Context, db *gorm.DB) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := services.CreateUser(db, &user)
	if err != nil {
		fmt.Printf(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"user": user})
}

func AddMonthAnalysis(c *gin.Context, db *gorm.DB) {
	var analysis models.MonthAnalysis
	if err := c.ShouldBindJSON(&analysis); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := services.AddMonthAnalysis(db, &analysis)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Month analysis added"})
}

func GetUserByPhoneNumber(c *gin.Context, db *gorm.DB) {
	phoneNumber := c.Param("phone_number")
	var user models.User

	fmt.Println(`user phone no: ` + phoneNumber)

	result := db.Where("phone_number = ?", "7893016461").First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	// Return the found user
	c.JSON(http.StatusOK, gin.H{"user": user})
}
