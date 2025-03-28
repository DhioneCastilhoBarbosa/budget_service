package database

import (
	"budget-service/models"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	// Obtém a URL do banco de dados da variável de ambiente
	dsn := os.Getenv("DATABASE_URL")

	// Verifica se a variável de ambiente está vazia
	if dsn == "" {
		log.Fatal("DATABASE_URL não encontrada. Certifique-se de que a variável de ambiente está configurada no servidor.")
	}

	// Conecta ao PostgreSQL
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Falha ao conectar ao banco de dados:", err)
	}

	// Realiza AutoMigrate
	err = DB.AutoMigrate(&models.Budget{})
	if err != nil {
		log.Fatalf("Erro ao realizar AutoMigrate: %v", err)
	}
}
