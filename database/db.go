package database

import (
	"budget-service/models"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Função para conectar ao banco de dados PostgreSQL
func ConnectDatabase() {
	var err error
	// Configuração da conexão com o banco de dados PostgreSQL
	dsn := "postgresql://postgress:postgress@69.62.88.198:5435/buget-service"
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Falha ao conectar ao banco de dados:", err)
	}

	// Você pode escolher não usar AutoMigrate aqui, já que não quer migração automática
	// Se você precisar realizar uma migração, descomente a linha abaixo
	DB.AutoMigrate(&models.Budget{})
}
