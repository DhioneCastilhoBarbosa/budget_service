package models

import "time"

// Definindo o modelo de orçamento
type Budget struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	UserID      *uint     `json:"user_id"` // Pode ser NULL para usuários não autenticados
	SessionID   string    `json:"session_id"`
	InstallerID uint      `json:"installer_id"`
	Value       float64   `json:"value"`
	Status      string    `json:"status" gorm:"default:pendente"`
	CreatedAt   time.Time `json:"created_at"`
}
