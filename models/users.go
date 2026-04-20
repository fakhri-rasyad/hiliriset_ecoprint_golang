package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)


type UserDataResponse struct {
    Username string `json:"username"`
    Email    string `json:"email"`
    Role     string `json:"role"`
}

type UserLoginResponse struct {
    Username    string `json:"username"`
    BearerToken string `json:"token"`
}

type UserRegisterRequest struct {
    Username string `json:"username" validate:"required"`
    Email    string `json:"email"    validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
}

type UserLoginRequest struct {
    Email    string `json:"email"    validate:"required,email"`
    Password string `json:"password" validate:"required"`
}


type UserBase struct {
    InternalID int64      `json:"internal_id"`
    PublicID   uuid.UUID  `json:"public_id"`
    Username   string     `json:"username"`
    Email      string     `json:"email"`
    Password   string     `json:"-"`         
    Role       string     `json:"role"`
    CreatedAt  time.Time  `json:"created_at"`
    UpdatedAt  time.Time  `json:"updated_at"`
    DeletedAt  *time.Time `json:"deleted_at"`
}


type UserGorm struct {
    InternalID int64          `json:"internal_id" db:"internal_id" gorm:"column:internal_id;primaryKey;autoIncrement"`
    PublicID   uuid.UUID      `json:"public_id"   db:"public_id"   gorm:"column:public_id;unique;type:uuid"`
    Username   string         `json:"username"    db:"username"    gorm:"column:username;unique"`
    Email      string         `json:"email"       db:"email"       gorm:"column:email;unique"`
    Password   string         `json:"-"           db:"password"    gorm:"column:password"`
    Role       string         `json:"role"        db:"role"        gorm:"column:role;type:user_role_enum;default:user"`
    CreatedAt  time.Time      `json:"created_at"  db:"created_at"  gorm:"column:created_at;autoCreateTime"`
    UpdatedAt  time.Time      `json:"updated_at"  db:"updated_at"  gorm:"column:updated_at;autoUpdateTime"`
    DeletedAt  gorm.DeletedAt `json:"-"           db:"deleted_at"  gorm:"column:deleted_at;index"`
}

func (UserGorm) TableName() string { return "users" }

func (u UserGorm) ToBase() UserBase {
    var deletedAt *time.Time
    if u.DeletedAt.Valid {
        deletedAt = &u.DeletedAt.Time
    }
    return UserBase{
        InternalID: u.InternalID,
        PublicID:   u.PublicID,
        Username:   u.Username,
        Email:      u.Email,
        Password:   u.Password,
        Role:       u.Role,
        CreatedAt:  u.CreatedAt,
        UpdatedAt:  u.UpdatedAt,
        DeletedAt:  deletedAt,
    }
}

func (u UserBase) ToDataResponse() UserDataResponse {
    return UserDataResponse{
        Username: u.Username,
        Email:    u.Email,
        Role:     u.Role,
    }
}