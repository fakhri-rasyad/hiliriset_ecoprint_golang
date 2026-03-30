package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	InternalID int64 `json:"internal_id" db:"internal_id" gorm:"primaryKey"`
	PublicID   uuid.UUID `json:"public_id" db:"public_id" gorm:"column:public_id"`
	Username string `json:"username" db:"username" gorm:"column:username"`
	Email string `json:"email" db:"email" gorm:"column:email"`
	Password string `json:"password" db:"password" gorm:"column:password"`
	Role string `json:"role" db:"role" gorm:"column:role"`
	CreatedAt time.Time `json:"created_at" db:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at" gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}