package models

import "time"

type Budget struct {
	ID               uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID           string     `json:"user_id" gorm:"index"`
	SessionID        string     `json:"session_id" gorm:"size:100;not null"`
	InstallerID      string     `json:"installer_id" gorm:"not null"`
	Name             string     `json:"name" gorm:"size:100"` // Nome do solicitante
	Email            string     `json:"email" gorm:"size:100"`
	Phone            string     `json:"phone" gorm:"size:20"`
	StationCount     uint       `json:"station_count"` // Número de estações
	LocationType     string     `json:"location_type" gorm:"size:100"`
	Photo1           string     `json:"photo1" gorm:"size:255"` // URL ou caminho para a imagem
	Photo2           string     `json:"photo2" gorm:"size:255"`
	Distance         string     `json:"distance" gorm:"size:100"`       // Distância até o ponto de energia, por exemplo
	NetworkType      string     `json:"network_type" gorm:"size:100"`   // Tipo de rede
	StructureType    string     `json:"structure_type" gorm:"size:100"` // Tipo de estrutura
	ChargerType      string     `json:"charger_type" gorm:"size:100"`
	Power            string     `json:"power" gorm:"size:100"`      // Ex: 7,4kW, 22kW
	Protection       string     `json:"protection" gorm:"size:255"` // Ex: disjuntor, DPS, DR
	Notes            string     `json:"notes" gorm:"type:text"`     // Observações
	InstallerName    string     `json:"installer_name" gorm:"size:100"`
	Value            float64    `json:"value"`
	Status           string     `json:"status" gorm:"default:aguardando orçamento"`
	ExecutionDate    *time.Time `json:"execution_date"`    // nova
	FinishDate       *time.Time `json:"finish_date"`       // nova
	PaymentStatus    string     `json:"payment_status"`    // "pendente" ou "pago"
	InstallerConfirm bool       `json:"installer_confirm"` // true/false
	ClientConfirm    bool       `json:"client_confirm"`    // true/false
	CreatedAt        time.Time  `json:"created_at" gorm:"autoCreateTime"`
}
