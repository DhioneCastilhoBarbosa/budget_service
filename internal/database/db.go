package database

import (
	"budget-service/internal/budget/models"
	"log"
	"os"

	"github.com/joho/godotenv" // ✅ importação do godotenv
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDatabase() {
	// Carrega o .env local se existir
	if os.Getenv("ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Println("⚠️  Arquivo .env não encontrado, seguindo com variáveis do sistema.")
		} else {
			log.Println("✅ Variáveis carregadas do .env")
		}
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("❌ DATABASE_URL não encontrada. Configure a variável de ambiente.")
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("❌ Falha ao conectar ao banco de dados:", err)
	}

	err = DB.AutoMigrate(&models.Budget{})
	if err != nil {
		log.Fatalf("❌ Erro ao realizar AutoMigrate: %v", err)
	}

	log.Println("✅ Banco de dados conectado com sucesso")
}
