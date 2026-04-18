package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type KomporRequestBody struct{
	KomporName string
}

type Kompors struct {
	InternalID int64 `json:"internal_id" db:"internal_id" gorm:"primaryKey"`
	PublicID   uuid.UUID `json:"public_id" db:"public_id" gorm:"public_id"`
	KomporName string `json:"kompor_name" db:"kompor_name" gorm:"kompor_name"`
	CreatedAt time.Time `json:"created_at" db:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at" gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	UserId int `json:"user_id" db:"user_id"`
}