package models

import (
	"time"

	"github.com/google/uuid"
)

type BoilingSessionBase struct {
    InternalID    int64     `json:"internal_id"`
    PublicID      uuid.UUID `json:"public_id"`
    BoilingStatus string    `json:"boiling_status"`
    FabricType    string    `json:"fabric_type"`
    UserID        *int64    `json:"user_id"`
    KomporID      *int64    `json:"kompor_id"`
    EspID         *int64    `json:"esp_id"`
    CreatedAt     time.Time `json:"created_at"`
    FinishedAt    *time.Time `json:"finished_at"`
    UpdatedAt     time.Time `json:"updated_at"`
    DeletedAt     *time.Time `json:"deleted_at"`
}

type BoilingSession struct {
    InternalID    int64      `json:"internal_id"    db:"internal_id"    gorm:"column:internal_id;primaryKey;autoIncrement"`
    PublicID      uuid.UUID  `json:"public_id"      db:"public_id"      gorm:"column:public_id;unique;type:uuid"`
    BoilingStatus string     `json:"boiling_status" db:"boiling_status" gorm:"column:boiling_status;type:boiling_status_enum;default:boiling"`
    FabricType    string     `json:"fabric_type"    db:"fabric_type"    gorm:"column:fabric_type;type:fabric_type_enum;default:katun"`
    UserID        *int64     `json:"user_id"        db:"user_id"        gorm:"column:user_id"`
    KomporID      *int64     `json:"kompor_id"      db:"kompor_id"      gorm:"column:kompor_id"`
    EspID         *int64     `json:"esp_id"         db:"esp_id"         gorm:"column:esp_id"`
    CreatedAt     time.Time  `json:"created_at"     db:"created_at"     gorm:"column:created_at;autoCreateTime"`
    FinishedAt    *time.Time `json:"finished_at"    db:"finished_at"    gorm:"column:finished_at"`
    UpdatedAt     time.Time  `json:"updated_at"     db:"updated_at"     gorm:"column:updated_at;autoUpdateTime"`
    DeletedAt     *time.Time `json:"deleted_at"     db:"deleted_at"     gorm:"column:deleted_at;index"`
}

func (BoilingSession) TableName() string {
    return "boiling_sessions"
}

func (bs BoilingSession) ToBase() BoilingSessionBase {
    return BoilingSessionBase{
        InternalID:    bs.InternalID,
        PublicID:      bs.PublicID,
        BoilingStatus: bs.BoilingStatus,
        FabricType:    bs.FabricType,
        UserID:        bs.UserID,
        KomporID:      bs.KomporID,
        EspID:         bs.EspID,
        CreatedAt:     bs.CreatedAt,
        FinishedAt:    bs.FinishedAt,
        UpdatedAt:     bs.UpdatedAt,
        DeletedAt:     bs.DeletedAt,
    }
}

func (bs *BoilingSession) FromBase(b BoilingSessionBase) {
    bs.InternalID    = b.InternalID
    bs.PublicID      = b.PublicID
    bs.BoilingStatus = b.BoilingStatus
    bs.FabricType    = b.FabricType
    bs.UserID        = b.UserID
    bs.KomporID      = b.KomporID
    bs.EspID         = b.EspID
    bs.CreatedAt     = b.CreatedAt
    bs.FinishedAt    = b.FinishedAt
    bs.UpdatedAt     = b.UpdatedAt
    bs.DeletedAt     = b.DeletedAt
}