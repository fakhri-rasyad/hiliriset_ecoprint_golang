package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SessionRecordInput struct {
    SessionID  uuid.UUID  `json:"session_pub_id"`
    AirTemp    float32    `json:"air_temp"`
    WaterTemp  float32    `json:"water_temp"`
    Humidity   float32    `json:"humidity"`
}

type SessionRecordOutput struct {
    PublicID   uuid.UUID  `json:"public_id"`
    SessionID  int64      `json:"session_id"`
    AirTemp    float32    `json:"air_temp"`
    WaterTemp  float32    `json:"water_temp"`
    Humidity   float32    `json:"humidity"`
    CreatedAt  time.Time  `json:"created_at"`
}

type SessionRecordBase struct {
    InternalID int64      `json:"internal_id"`
    PublicID   uuid.UUID  `json:"public_id"`
    SessionID  int64      `json:"session_id"`
    AirTemp    float32    `json:"air_temp"`
    WaterTemp  float32    `json:"water_temp"`
    Humidity   float32    `json:"humidity"`
    CreatedAt  time.Time  `json:"created_at"`
    DeletedAt  *time.Time `json:"deleted_at"`
}

func (rb *SessionRecordBase) ToOutput() SessionRecordOutput {
  return SessionRecordOutput{
    PublicID: rb.PublicID,
    SessionID: rb.SessionID,
    AirTemp: rb.AirTemp,
    WaterTemp: rb.WaterTemp,
    Humidity: rb.Humidity,
    CreatedAt: rb.CreatedAt,
  }
}


type SessionRecordGorm struct {
    InternalID int64          `json:"internal_id" db:"internal_id" gorm:"column:internal_id;primaryKey;autoIncrement"`
    PublicID   uuid.UUID      `json:"public_id"   db:"public_id"   gorm:"column:public_id;unique;type:uuid"`
    SessionID  int64          `json:"session_id"  db:"session_id"  gorm:"column:session_id;index"`
    AirTemp    float32        `json:"air_temp"    db:"air_temp"    gorm:"column:air_temp"`
    WaterTemp  float32        `json:"water_temp"  db:"water_temp"  gorm:"column:water_temp"`
    Humidity   float32        `json:"humidity"    db:"humidity"    gorm:"column:humidity"`
    CreatedAt  time.Time      `json:"created_at"  db:"created_at"  gorm:"column:created_at;autoCreateTime"`
    DeletedAt  gorm.DeletedAt `json:"-"           db:"deleted_at"  gorm:"column:deleted_at;index"`
}

func (SessionRecordGorm) TableName() string { return "session_records" }

func (s SessionRecordGorm) ToBase() SessionRecordBase {
    var deletedAt *time.Time
    if s.DeletedAt.Valid {
        deletedAt = &s.DeletedAt.Time
    }
    return SessionRecordBase{
        InternalID: s.InternalID,
        PublicID:   s.PublicID,
        SessionID:  s.SessionID,
        AirTemp:    s.AirTemp,
        WaterTemp:  s.WaterTemp,
        Humidity:   s.Humidity,
        CreatedAt:  s.CreatedAt,
        DeletedAt:  deletedAt,
    }
}
