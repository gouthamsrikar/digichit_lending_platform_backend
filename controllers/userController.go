package controllers

import (
	"chitfund/httpclient"
	"chitfund/models"
	"chitfund/services"
	"fmt"
	"math"
	"net/http"
	"time"

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
	// && user.BankStatementFetched
	if result.Error == nil && user.BankStatementFetched {

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

	if response.Data.KeyDetails.BANKSTATEMENT.KeyName == "" && len(response.Data.KeyDetails.BANKSTATEMENT.KeyData) == 0 {
		pan := response.Data.KeyDetails.PAN.KeyData[0]
		currentTime := time.Now()

		formattedDate := currentTime.Format("02/01/2006")
		user := models.User{
			Name:                 pan.Name,
			PanNo:                string(pan.KeyID[0]) + "XXXXXXX" + string(pan.KeyID[len(pan.KeyID)-2:]),
			PhoneNumber:          request.PhoneNumber,
			JoinDate:             formattedDate,
			BankStatementFetched: false,
		}

		err = services.CreateUser(db, &user)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"user": user, "bank_statement": nil})

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

	currentTime := time.Now()

	// Format the date to DD/MM/YYYY
	formattedDate := currentTime.Format("02/01/2006")

	user = models.User{
		Name:                 pan.Name,
		PanNo:                string(pan.KeyID[0]) + "XXXXXXX" + string(pan.KeyID[len(pan.KeyID)-2:]),
		BankAccountNo:        bankStatement.AccountNumber,
		CreditToDebitRatio:   creditToDebitRatio,
		CommunityLevel:       "STARTER",
		DigitScore:           0,
		MonthlyIncome:        monthlyIncome,
		AvgBalance:           avgBalance,
		EmiToIncomeRatio:     0,
		PhoneNumber:          request.PhoneNumber,
		JoinDate:             formattedDate,
		BankStatementFetched: true,
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

	digitScore := calculateCombinedFinancialScore(bankStatements)
	var communityLevel string

	if digitScore >= 80 {
		communityLevel = "PREMIUM"
	} else if digitScore >= 65 {
		communityLevel = "STANDARD"
	} else {
		communityLevel = "STARTER"
	}

	err = db.Model(&models.User{}).Where("id = ?", user.ID).Updates(
		map[string]interface{}{
			"digit_score":     digitScore,
			"community_level": communityLevel,
		}).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to update"})
		return
	}

	user.DigitScore = digitScore
	user.CommunityLevel = communityLevel

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

func calculateCombinedFinancialScore(statements []models.BankStatement) int {
	var totalCreditSum, totalDebitSum, totalAverageEodBalance float64
	var totalDebitTransactions, totalCreditTransactions int

	for _, statement := range statements {
		totalCreditSum += statement.TotalCreditAmount
		totalDebitSum += statement.TotalDebitAmount
		totalDebitTransactions += statement.NoOfDebitTransactions
		totalCreditTransactions += statement.NoOfCreditTransactions
		totalAverageEodBalance += statement.AverageEodBalance
	}

	entryCount := float64(len(statements))
	if entryCount == 0 {
		return 0
	}

	avgCredit := totalCreditSum / entryCount
	avgDebit := totalDebitSum / entryCount
	avgEodBalance := totalAverageEodBalance / entryCount
	avgTransactions := float64(totalDebitTransactions+totalCreditTransactions) / entryCount

	incomeScore := math.Min((avgCredit/100000)*20, 20)

	balanceScore := math.Min((avgEodBalance/(avgCredit+1))*15, 15)

	creditDebitRatio := 1.0
	if avgDebit > 0 {
		creditDebitRatio = avgCredit / avgDebit
	}
	creditDebitScore := 15.0
	if creditDebitRatio < 0.95 {
		creditDebitScore = math.Min(creditDebitRatio*15, 15)
	}

	activityScore := math.Min((avgTransactions/50)*10, 10)

	combinedScore := (incomeScore + balanceScore + creditDebitScore + activityScore) / 50 * 100

	return int(math.Round(combinedScore))
}
