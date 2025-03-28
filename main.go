package main

import (
	"budget-service/database"
	"budget-service/routes"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	// Conectar ao banco de dados
	database.ConnectDatabase()

	// Iniciar o servidor
	router := gin.Default()

	// Definir rotas
	routes.SetupRoutes(router)
	fmt.Println("Server running on port: 8082")
	// Rodar a aplicação na porta 8082
	router.Run(":8082")
}
