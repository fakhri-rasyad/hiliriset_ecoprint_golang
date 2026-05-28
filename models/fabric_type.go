package models

import (
	"time"

	"github.com/google/uuid"
)

type FabricType struct {
    InternalID     int64      `json:"internal_id"     gorm:"column:internal_id;primaryKey;autoIncrement"`
    PublicID       uuid.UUID  `json:"public_id"       gorm:"column:public_id;unique;type:uuid"`
    Name           string     `json:"name"            gorm:"column:name;type:fabric_type_enum"`
    BoilingMinutes int        `json:"boiling_minutes" gorm:"column:boiling_minutes"`
    CreatedAt      time.Time  `json:"created_at"      gorm:"column:created_at;autoCreateTime"`
    UpdatedAt      time.Time  `json:"updated_at"      gorm:"column:updated_at;autoUpdateTime"`
}

func (FabricType) TableName() string { return "fabric_types" }
