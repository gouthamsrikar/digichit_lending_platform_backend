package routes

import (
	"chitfund/controllers"
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterUserRoutes(router *gin.Engine, db *gorm.DB) {
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
