package database

import (
	"budget-service/models"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Função para conectar ao banco de dados PostgreSQL
func ConnectDatabase() {
	// Carregar as variáveis do arquivo .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Erro ao carregar o arquivo .env", err)
	}

	// Obter a URL do banco de dados do arquivo .env
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL não encontrada no arquivo .env")
	}

	// Conectar ao banco de dados PostgreSQL usando a variável de ambiente
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Falha ao conectar ao banco de dados:", err)
	}

	// Usando AutoMigrate para criar a tabela se ela não existir
	err = DB.AutoMigrate(&models.Budget{})
	if err != nil {
		log.Fatalf("Erro ao realizar AutoMigrate: %v", err)
	}
}
