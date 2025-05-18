package budget

import (
	"budget-service/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	v1 := router.Group("/api/v1/budget")
	{
		// Cria um novo orçamento
		v1.POST("/", CreateBudget)

		// Busca orçamentos por user_id (via query string ou param)
		v1.GET("/", GetUserBudgets)

		// Busca todos os orçamentos
		v1.GET("/all", GetAllBudgets)

		// Vincula orçamentos a um usuário após login
		v1.PUT("/link", LinkBudgetsToUser)

		protected := v1.Group("/", middlewares.AuthMiddleware())
		{
			protected.PUT("/:id/value", UpdateBudgetValue)
			protected.PUT("/:id/status", UpdateBudgetStatus)
			protected.PUT("/:id/dates", UpdateBudgetDates)
			protected.PUT("/:id/payment", UpdatePaymentStatus)
			protected.PUT("/:id/confirm", ConfirmExecution)
		}
	}
}
