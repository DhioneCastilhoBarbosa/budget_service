package routes

import (
	"budget-service/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	budgetGroup := router.Group("/api/budget")
	{
		budgetGroup.POST("/", controllers.CreateBudget)
		budgetGroup.GET("/:user_id", controllers.GetUserBudgets)
		budgetGroup.PUT("/link", controllers.LinkBudgetsToUser)
	}
}
