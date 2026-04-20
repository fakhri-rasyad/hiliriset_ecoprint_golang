package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)


type CreateEspRequest struct {
    MacAddress string `json:"mac_address" validate:"required"`
}

type EspResponse struct {
    PublicID     uuid.UUID `json:"public_id"`
    MacAddress   string    `json:"mac_address"`
    DeviceStatus string    `json:"device_status"`
}


type EspBase struct {
    InternalID   int64      `json:"internal_id"`
    PublicID     uuid.UUID  `json:"public_id"`
    MacAddress   string     `json:"mac_address"`
    DeviceStatus string     `json:"device_status"`
    UserID       *int64     `json:"user_id"`
    CreatedAt    time.Time  `json:"created_at"`
    UpdatedAt    time.Time  `json:"updated_at"`
    DeletedAt    *time.Time `json:"deleted_at"`
}


type EspGorm struct {
    InternalID   int64          `json:"internal_id"   db:"internal_id"   gorm:"column:internal_id;primaryKey;autoIncrement"`
    PublicID     uuid.UUID      `json:"public_id"     db:"public_id"     gorm:"column:public_id;unique;type:uuid"`
    MacAddress   string         `json:"mac_address"   db:"mac_address"   gorm:"column:mac_address;unique"`
    DeviceStatus string         `json:"device_status" db:"device_status" gorm:"column:device_status;type:device_status_enum;default:offline"`
    UserID       *int64         `json:"user_id"       db:"user_id"       gorm:"column:user_id"`
    CreatedAt    time.Time      `json:"created_at"    db:"created_at"    gorm:"column:created_at;autoCreateTime"`
    UpdatedAt    time.Time      `json:"updated_at"    db:"updated_at"    gorm:"column:updated_at;autoUpdateTime"`
    DeletedAt    gorm.DeletedAt `json:"-"             db:"deleted_at"    gorm:"column:deleted_at;index"`
}

func (EspGorm) TableName() string { return "esps" }

func (e EspGorm) ToBase() EspBase {
    var deletedAt *time.Time
    if e.DeletedAt.Valid {
        deletedAt = &e.DeletedAt.Time
    }
    return EspBase{
        InternalID:   e.InternalID,
        PublicID:     e.PublicID,
        MacAddress:   e.MacAddress,
        DeviceStatus: e.DeviceStatus,
        UserID:       e.UserID,
        CreatedAt:    e.CreatedAt,
        UpdatedAt:    e.UpdatedAt,
        DeletedAt:    deletedAt,
    }
}

func (e EspBase) ToResponse() EspResponse {
    return EspResponse{
        PublicID:     e.PublicID,
        MacAddress:   e.MacAddress,
        DeviceStatus: e.DeviceStatus,
    }
}