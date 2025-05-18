package main

import (
	"budget-service/internal/budget"
	"budget-service/internal/database"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Inicializa o banco de dados
	database.ConnectDatabase()

	// Inicializa e configura o servidor
	router := gin.Default()

	// ✅ Middleware de CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Configura limite de upload de arquivos
	router.MaxMultipartMemory = 100 << 20 // 100 MB

	// Registra as rotas da API
	budget.SetupRoutes(router)

	log.Println("✅ Servidor rodando na porta 8082")
	if err := router.Run(":8082"); err != nil {
		log.Fatalf("❌ Erro ao iniciar o servidor: %v", err)
	}
}
