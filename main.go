package main

import (
	"chitfund/config"
	"chitfund/httpclient"
	"chitfund/models"

	"chitfund/routes"
	// "log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func main() {
	r := gin.Default()

	db := config.InitDB()

	cfg := config.LoadConfig()

	client := httpclient.NewHttpClient()

	httpService := httpclient.NewService(cfg, client)

	addSampleCommunityData(db)

	routes.RegisterRoutes(r, db, httpService)

	r.Run(":8080")
}

func addSampleCommunityData(db *gorm.DB) {
	sampleCommunities := []models.Community{
		{CommunityName: "Swiggy Driver Cyberabad Fund", MonthlyDeposit: 500.0, TotalFund: 35000.0, InterestRate: 0.1, AdminName: "Ashok Patel", AdminUserID: "", RepaymentPeriodInMonths: 6, CommunityDescription: "Community fund created for the aid of Swiggy Delivery Partners in Cyberabad area", UserCount: 3, MaxCount: 10},
		{CommunityName: "Swiggy Driver Hyderabad Fund", MonthlyDeposit: 750.0, TotalFund: 40000.0, InterestRate: 0.12, AdminName: "Rakesh Ahuja", AdminUserID: "", RepaymentPeriodInMonths: 6, CommunityDescription: "Community fund created for the aid of Swiggy Delivery Partners in old Hyderabad area", UserCount: 26, MaxCount: 30},
		{CommunityName: "Swiggy Driver Secunderabad Fund", MonthlyDeposit: 500.0, TotalFund: 25000.0, InterestRate: 0.08, AdminName: "Ajith Kumar", AdminUserID: "", RepaymentPeriodInMonths: 6, CommunityDescription: "Community fund created for the aid of Swiggy Delivery Partners in Secunderabad area", UserCount: 8, MaxCount: 15},
		{CommunityName: "Swiggy Driver Gachibowli Fund", MonthlyDeposit: 600.0, TotalFund: 35000.0, InterestRate: 0.1, AdminName: "Steve Martin", AdminUserID: "", RepaymentPeriodInMonths: 6, CommunityDescription: "Community fund created for the aid of Swiggy Delivery Partners in Gachibowli area", UserCount: 16, MaxCount: 25},
	}
	db.Create(&sampleCommunities)
}
