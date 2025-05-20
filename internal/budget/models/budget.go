package models

import "time"

type Budget struct {
	ID               uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID           *string    `json:"user_id" gorm:"index"` // Agora aceita NULL
	SessionID        string     `json:"session_id" gorm:"size:100;not null"`
	InstallerID      string     `json:"installer_id" gorm:"not null"`
	Name             string     `json:"name" gorm:"size:100"`
	Email            *string    `json:"email" gorm:"size:100"`
	Phone            *string    `json:"phone" gorm:"size:20"`
	CEP              *string    `json:"cep" gorm:"size:20"`
	Street           *string    `json:"street" gorm:"size:100"`
	Number           *string    `json:"number" gorm:"size:20"`
	Neighborhood     *string    `json:"neighborhood" gorm:"size:100"`
	City             *string    `json:"city" gorm:"size:100"`
	State            *string    `json:"state" gorm:"size:2"`
	Complement       *string    `json:"complement" gorm:"size:100"`
	StationCount     uint       `json:"station_count"`
	LocationType     *string    `json:"location_type" gorm:"size:100"`
	Photo1           *string    `json:"photo1" gorm:"size:255"`
	Photo2           *string    `json:"photo2" gorm:"size:255"`
	Distance         *string    `json:"distance" gorm:"size:100"`
	NetworkType      *string    `json:"network_type" gorm:"size:100"`
	StructureType    *string    `json:"structure_type" gorm:"size:100"`
	ChargerType      *string    `json:"charger_type" gorm:"size:100"`
	Power            *string    `json:"power" gorm:"size:100"`
	Protection       *string    `json:"protection" gorm:"size:255"`
	Notes            *string    `json:"notes" gorm:"type:text"`
	InstallerName    *string    `json:"installer_name" gorm:"size:100"`
	Value            float64    `json:"value"`
	Status           string     `json:"status" gorm:"default:aguardando orçamento"`
	ExecutionDate    *time.Time `json:"execution_date"`
	FinishDate       *time.Time `json:"finish_date"`
	PaymentStatus    *string    `json:"payment_status"`
	InstallerConfirm bool       `json:"installer_confirm"`
	ClientConfirm    bool       `json:"client_confirm"`
	CreatedAt        time.Time  `json:"created_at" gorm:"autoCreateTime"`
}
