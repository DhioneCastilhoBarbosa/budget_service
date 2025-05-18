package main

import (
	"budget-service/internal/budget"
	"budget-service/internal/database"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Inicializa o banco de dados
	database.ConnectDatabase()

	// Inicializa e configura o servidor
	router := gin.Default()
	router.MaxMultipartMemory = 100 << 20 // 64MB
	budget.SetupRoutes(router)

	log.Println("✅ Servidor rodando na porta 8082")
	err := router.Run(":8082")
	if err != nil {
		log.Fatalf("❌ Erro ao iniciar o servidor: %v", err)
	}
}
