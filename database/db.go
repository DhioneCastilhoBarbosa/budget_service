package database

import (
	"budget-service/models"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	_ "modernc.org/sqlite" // Driver SQLite sem CGO
)

var DB *gorm.DB

func ConnectDatabase() {
	var err error
	// Especificando o driver do SQLite explicitamente
	DB, err = gorm.Open(sqlite.Dialector{
		DSN:        "budgets.db",
		DriverName: "sqlite",
	}, &gorm.Config{})

	if err != nil {
		log.Fatal("Falha ao conectar ao banco de dados:", err)
	}

	// AutoMigrate para criar a tabela
	DB.AutoMigrate(&models.Budget{})
}
