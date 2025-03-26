package main

import (
	"budget-service/database"
	"budget-service/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Conectar ao banco de dados
	database.ConnectDatabase()

	// Iniciar o servidor
	router := gin.Default()

	// Definir rotas
	routes.SetupRoutes(router)

	// Rodar a aplicação na porta 8082
	router.Run(":8082")
}
