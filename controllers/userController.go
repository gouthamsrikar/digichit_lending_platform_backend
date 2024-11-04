package controllers

import (
	"chitfund/httpclient"
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

func CreateUserWithIdempotencyId(c *gin.Context, db *gorm.DB, httpService *httpclient.Service) {
	var request struct {
		IdempotencyId string `json:"idempotency_id"`
		PhoneNumber   string `json:"phone_number"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var user models.User

	result := db.Where("phone_number = ?", request.PhoneNumber).First(&user)
	if result.Error == nil {

		bankStatement, err := GetBankStatementsByUserID(db, user.ID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"user": user, "bank_statement": bankStatement})
		return
	}

	response, err := httpService.GetIdgUserData(request.IdempotencyId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	pan := response.Data.KeyDetails.PAN.KeyData[0]
	bankStatement := response.Data.KeyDetails.BANKSTATEMENT.KeyData[0]

	totalCreditAmount := 0.0
	totalDebitAmount := 0.0

	for _, month := range bankStatement.MonthWiseAnalysis {
		totalCreditAmount += month.TotalCreditAmount
		totalDebitAmount += month.TotalDebitAmount

	}

	creditToDebitRatio := 0.0
	monthlyIncome := 0.0
	avgBalance := 0.0
	if len(bankStatement.MonthWiseAnalysis) > 0 {
		monthlyIncome = totalCreditAmount / float64(len(bankStatement.MonthWiseAnalysis))
		creditToDebitRatio = totalDebitAmount / totalCreditAmount
		avgBalance = (totalCreditAmount - totalDebitAmount) / float64(len(bankStatement.MonthWiseAnalysis))
	}

	user = models.User{
		Name:               pan.Name,
		PanNo:              string(pan.KeyID[0]) + "XXXXXXX" + string(pan.KeyID[len(pan.KeyID)-2:]),
		BankAccountNo:      bankStatement.AccountNumber,
		CreditToDebitRatio: creditToDebitRatio,
		CommunityLevel:     "STARTER",
		DigitScore:         12,
		MonthlyIncome:      monthlyIncome,
		AvgBalance:         avgBalance,
		EmiToIncomeRatio:   12,
		PhoneNumber:        request.PhoneNumber,
	}

	err = services.CreateUser(db, &user)

	if err != nil {
		fmt.Printf(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var bankStatements []models.BankStatement
	for _, month := range bankStatement.MonthWiseAnalysis {

		statement := models.BankStatement{
			UserID:                 user.ID,
			MonthName:              month.MonthName,
			NoOfDebitTransactions:  month.NoOfDebitTransactions,
			NoOfCreditTransactions: month.NoOfCreditTransactions,
			TotalCreditAmount:      month.TotalCreditAmount,
			Year:                   month.Year,
			TotalDebitAmount:       month.TotalDebitAmount,
			AverageEodBalance:      month.AverageEodBalance,
		}
		bankStatements = append(bankStatements, statement)
	}

	err = CreateBankStatements(db, bankStatements)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"user": user, "bank_statement": bankStatements})
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

	result := db.Where("phone_number = ?", phoneNumber).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}
