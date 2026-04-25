package repositories

import (
	"hiliriset_ecoprint_golang/config"
	"hiliriset_ecoprint_golang/models"

	"github.com/google/uuid"
)

type SeReRepository interface {
    CreateRecord(record *models.SessionRecordGorm) error
    GetRecords(sessionID int64) ([]models.SessionRecordBase, error)
    GetRecordByPubID(publicID uuid.UUID) (*models.SessionRecordBase, error)
}

type SeReRepositoryImpl struct{}

func NewSeReRepository() SeReRepository {
    return &SeReRepositoryImpl{}
}

func (r *SeReRepositoryImpl) CreateRecord(record *models.SessionRecordGorm) error {
    return config.DB.Create(record).Error
}

func (r *SeReRepositoryImpl) GetRecords(sessionID int64) ([]models.SessionRecordBase, error) {
    var gormModels []models.SessionRecordGorm

    if err := config.DB.Where("session_id = ?", sessionID).Find(&gormModels).Error; err != nil {
        return nil, err
    }

    result := make([]models.SessionRecordBase, len(gormModels))
    for i, gormModel := range gormModels {
        result[i] = gormModel.ToBase()
    }

    return result, nil
}

func (r *SeReRepositoryImpl) GetRecordByPubID(publicID uuid.UUID) (*models.SessionRecordBase, error) {
    var gormModel models.SessionRecordGorm

    if err := config.DB.Where("public_id = ?", publicID).First(&gormModel).Error; err != nil {
        return nil, err
    }

    base := gormModel.ToBase()
    return &base, nil
}
