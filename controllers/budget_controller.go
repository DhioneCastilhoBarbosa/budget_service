package controllers

import (
	"budget-service/database"
	"budget-service/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Criar um orçamento (chamado pelo Chat Service)
func CreateBudget(c *gin.Context) {
	var budget models.Budget

	if err := c.ShouldBindJSON(&budget); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database.DB.Create(&budget)
	c.JSON(http.StatusOK, budget)
}

// Buscar orçamentos do usuário autenticado
func GetUserBudgets(c *gin.Context) {
	userID := c.Param("user_id")
	var budgets []models.Budget

	database.DB.Where("user_id = ?", userID).Find(&budgets)
	c.JSON(http.StatusOK, budgets)
}

// Vincular orçamentos ao usuário após login
func LinkBudgetsToUser(c *gin.Context) {
	var req struct {
		SessionID string `json:"session_id"`
		UserID    uint   `json:"user_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Atualiza os orçamentos com base no SessionID
	database.DB.Model(&models.Budget{}).
		Where("session_id = ?", req.SessionID).
		Update("user_id", req.UserID)

	c.JSON(http.StatusOK, gin.H{"message": "Orçamentos vinculados ao usuário"})
}
