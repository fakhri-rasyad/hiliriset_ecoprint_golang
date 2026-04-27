package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type KomporRequest struct {
    KomporName string `json:"kompor_name" validate:"required"`
}

type KomporResponse struct {
    PublicID   uuid.UUID `json:"public_id"`
    KomporName string    `json:"kompor_name"`
}


type KomporBase struct {
    InternalID int64      `json:"internal_id"`
    PublicID   uuid.UUID  `json:"public_id"`
    KomporName string     `json:"kompor_name"`
    UserID     *int64     `json:"user_id"`
    IsActive   bool       `json:"is_active"`
    CreatedAt  time.Time  `json:"created_at"`
    UpdatedAt  time.Time  `json:"updated_at"`
    DeletedAt  *time.Time `json:"deleted_at"`
}


type KomporGorm struct {
    InternalID int64          `json:"internal_id" db:"internal_id" gorm:"column:internal_id;primaryKey;autoIncrement"`
    PublicID   uuid.UUID      `json:"public_id"   db:"public_id"   gorm:"column:public_id;unique;type:uuid"`
    KomporName string         `json:"kompor_name" db:"kompor_name" gorm:"column:kompor_name"`
    UserID     *int64         `json:"user_id"     db:"user_id"     gorm:"column:user_id"`
    IsActive   bool           `json:"is_active"   db:"is_active"   gorm:"column:is_active;default:false"`
    CreatedAt  time.Time      `json:"created_at"  db:"created_at"  gorm:"column:created_at;autoCreateTime"`
    UpdatedAt  time.Time      `json:"updated_at"  db:"updated_at"  gorm:"column:updated_at;autoUpdateTime"`
    DeletedAt  gorm.DeletedAt `json:"-"           db:"deleted_at"  gorm:"column:deleted_at;index"`
}

func (KomporGorm) TableName() string { return "kompors" }

func (k KomporGorm) ToBase() KomporBase {
    var deletedAt *time.Time
    if k.DeletedAt.Valid {
        deletedAt = &k.DeletedAt.Time
    }
    return KomporBase{
        InternalID: k.InternalID,
        PublicID:   k.PublicID,
        KomporName: k.KomporName,
        UserID:     k.UserID,
        IsActive:   k.IsActive,
        CreatedAt:  k.CreatedAt,
        UpdatedAt:  k.UpdatedAt,
        DeletedAt:  deletedAt,
    }
}

func (k KomporBase) ToResponse() KomporResponse {
    return KomporResponse{
        PublicID:   k.PublicID,
        KomporName: k.KomporName,
    }
}
