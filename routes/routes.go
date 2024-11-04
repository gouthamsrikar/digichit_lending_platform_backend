package routes

import (
	"chitfund/controllers"
	"chitfund/httpclient"
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(router *gin.Engine, db *gorm.DB, httpService *httpclient.Service) {

	transactionrGroup := router.Group("/init")
	{
		transactionrGroup.POST("/", func(c *gin.Context) {
			controllers.InitTransaction(c, db, httpService)
		})
	}

	userGroup := router.Group("/user")
	{
		userGroup.POST("/", func(c *gin.Context) {
			controllers.CreateUser(c, db)
		})
		userGroup.POST("/month_analysis", func(c *gin.Context) {
			controllers.AddMonthAnalysis(c, db)
		})

		userGroup.GET("/:phone_number", func(c *gin.Context) {
			fmt.Print("get user request")
			controllers.GetUserByPhoneNumber(c, db)
		})

		userGroup.POST("/idg", func(c *gin.Context) {
			controllers.CreateUserWithIdempotencyId(c, db, httpService)
		})

	}
	communityGroup := router.Group("/community")
	{
		communityGroup.GET("/:user_id", func(c *gin.Context) {
			controllers.FetchAllCommunities(c, db)
		})

		communityGroup.POST("/request", func(c *gin.Context) {
			controllers.RequestToJoinCommunity(c, db)
		})

		communityGroup.POST("/approve/:ledger_id", func(c *gin.Context) {
			controllers.ApproveToJoinCommunity(c, db)
		})
	}

	loanGroup := router.Group("/loan")
	{
		loanGroup.POST("/request", func(c *gin.Context) {
			controllers.RequestLoan(c, db)
		})
		loanGroup.GET("/:user_id", func(c *gin.Context) {
			controllers.GetLoanLedgersByUserID(c, db)
		})

		loanGroup.POST("approve/:loan_id", func(c *gin.Context) {
			controllers.ApproveLoan(c, db)
		})
	}
}
