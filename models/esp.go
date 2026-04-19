package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CreateEspRequest struct {
	MacAddress string `json:"mac_address" validate:"required"`
}

type EspBase struct{
	PublicID uuid.UUID `json:"public_id" db:"public_id" gorm:"column:public_id"`
	MacAddress string `json:"mac_address" db:"mac_address" gorm:"column:mac_address"`
	DeviceStatus string `json:"device_status" db:"device_status" gorm:"column:device_status"`
}

type EspGorm struct {
	EspBase 
	InternalID int64 `json:"internal_id" db:"internal_id" gorm:"column:internal_id;primaryKey"`
	UserID int64 `json:"user_id" db:"user_id" gorm:"column:user_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at" gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" db:"deleted_at" gorm:"column:deleted_at"`
}

func (EspGorm) TableName() string {
	return "esps"
}