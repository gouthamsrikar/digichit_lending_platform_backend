package controllers

import (
	"chitfund/httpclient"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func InitTransaction(c *gin.Context, db *gorm.DB, httpService *httpclient.Service) {
	var request struct {
		PhoneNumber string `json:"phone_number"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	idempotency_uuid := uuid.New()

	fmt.Printf(idempotency_uuid.String())

	response, err := httpService.InitTransaction(idempotency_uuid.String(), request.PhoneNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"token": response.Token, "idempotency_id": idempotency_uuid})
}
