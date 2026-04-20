package repositories

import (
	"hiliriset_ecoprint_golang/config"
	"hiliriset_ecoprint_golang/models"

	"github.com/google/uuid"
)

type BoSeRepository interface {
    CreateSession(userID int64, komporID int64, espID int64, fabricType string) (*models.BoilingSessionBase, error)
    GetSessions(userID int64) ([]models.BoilingSessionBase, error)
    GetSessionByPublicID(publicID uuid.UUID) (*models.BoilingSessionBase, error)
    UpdateSessionStatus(publicID uuid.UUID, status string) error
    FinishSession(publicID uuid.UUID) error
}

type BoSeRepositoryImpl struct{}

func NewBoSeRepository() BoSeRepository {
    return &BoSeRepositoryImpl{}
}

func (r *BoSeRepositoryImpl) CreateSession(userID int64, komporID int64, espID int64, fabricType string) (*models.BoilingSessionBase, error) {
    gormModel := models.BoilingSession{
        UserID:     &userID,
        KomporID:   &komporID,
        EspID:      &espID,
        FabricType: fabricType,
		PublicID: uuid.New(),
    }

    if err := config.DB.Create(&gormModel).Error; err != nil {
        return nil, err
    }

    base := gormModel.ToBase()
    return &base, nil
}

func (r *BoSeRepositoryImpl) GetSessions(userID int64) ([]models.BoilingSessionBase, error) {
    var gormModels []models.BoilingSession

    if err := config.DB.Where("user_id = ?", userID).Find(&gormModels).Error; err != nil {
        return nil, err
    }

    result := make([]models.BoilingSessionBase, len(gormModels))
    for i, g := range gormModels {
        result[i] = g.ToBase()
    }

    return result, nil
}

func (r *BoSeRepositoryImpl) GetSessionByPublicID(publicID uuid.UUID) (*models.BoilingSessionBase, error) {
    var gormModel models.BoilingSession

    if err := config.DB.Where("public_id = ?", publicID).First(&gormModel).Error; err != nil {
        return nil, err
    }

    base := gormModel.ToBase()
    return &base, nil
}

func (r *BoSeRepositoryImpl) UpdateSessionStatus(publicID uuid.UUID, status string) error {
    return config.DB.Model(&models.BoilingSession{}).
        Where("public_id = ?", publicID).
        Update("boiling_status", status).
        Error
}

func (r *BoSeRepositoryImpl) FinishSession(publicID uuid.UUID) error {
    return config.DB.Model(&models.BoilingSession{}).
        Where("public_id = ?", publicID).
        Updates(map[string]any{
            "boiling_status": "finished",
            "finished_at":    config.DB.NowFunc(),
        }).
        Error
}