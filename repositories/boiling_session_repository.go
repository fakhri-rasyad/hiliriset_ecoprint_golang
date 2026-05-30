// ========================================================
// repositories/bose_repository.go
// ========================================================
package repositories

import (
	"hiliriset_ecoprint_golang/config"
	"hiliriset_ecoprint_golang/models"
	"time"

	"github.com/google/uuid"
)

type BoSeRepository interface {
    CreateSession(userID int64, komporID int64, espID int64, espPubId uuid.UUID, fabricTypeID int64) (*models.BoilingSessionBase, error)
    GetSessions(userID int64) ([]models.BoilingSessionBase, error)
    GetSessionByPublicID(publicID uuid.UUID) (*models.BoilingSessionBase, error)
    UpdateSessionStatus(publicID uuid.UUID, status string) error
    FinishSession(publicID uuid.UUID) error
    CancelSession(publicID uuid.UUID) error
    SetFinishedAt(publicID uuid.UUID, finishedAt time.Time) error
    GetActiveSessionByEspPublicID(espPublicID uuid.UUID) (*models.BoilingSessionBase, error)
    GetActiveSessions() ([]models.BoilingSessionBase, error)
}

type BoSeRepositoryImpl struct{}

func NewBoSeRepository() BoSeRepository {
    return &BoSeRepositoryImpl{}
}

func (r *BoSeRepositoryImpl) CreateSession(
    userID int64,
    komporID int64,
    espID int64,
    espPubId uuid.UUID,
    fabricTypeID int64,
) (*models.BoilingSessionBase, error) {
    gormModel := models.BoilingSession{
        PublicID:   uuid.New(),
        UserID:     &userID,
        KomporID:   &komporID,
        EspID:      &espID,
        FabricType: fabricTypeID,
        FinishedAt: nil,
    }

    if err := config.DB.Create(&gormModel).Error; err != nil {
        return nil, err
    }

    base := gormModel.ToBase()
    base.EspPublicID = &espPubId
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

	if gormModel.EspID != nil {
		var esp models.EspGorm
		if err := config.DB.
			Select("public_id").
			Where("internal_id = ?", *gormModel.EspID).
			First(&esp).Error; err == nil {
			base.EspPublicID = &esp.PublicID
		}
	}

	return &base, nil
}

func (r *BoSeRepositoryImpl) UpdateSessionStatus(publicID uuid.UUID, status string) error {
    return config.DB.Model(&models.BoilingSession{}).
        Where("public_id = ?", publicID).
        Update("boiling_status", status).
        Error
}

func (r *BoSeRepositoryImpl) CancelSession(publicID uuid.UUID) error {
    return config.DB.Model(&models.BoilingSession{}).
        Where("public_id = ?", publicID).
        Updates(map[string]any{
            "boiling_status": "cancelled",
        }).
        Error
}

func (r *BoSeRepositoryImpl) SetFinishedAt(publicID uuid.UUID, finishedAt time.Time) error {
    return config.DB.Model(&models.BoilingSession{}).
        Where("public_id = ?", publicID).
        Updates(map[string]any{
            "finished_at": finishedAt,
            "boiling_status": "steaming",
        }).
        Error
}

func (r *BoSeRepositoryImpl) FinishSession(publicID uuid.UUID) error {
    return config.DB.Model(&models.BoilingSession{}).
        Where("public_id = ?", publicID).
        Updates(map[string]any{
            "boiling_status": "finished",
        }).
        Error
}

func (r *BoSeRepositoryImpl) GetActiveSessionByEspPublicID(espPublicID uuid.UUID) (*models.BoilingSessionBase, error) {
	var gormModel models.BoilingSession

	err := config.DB.
		Joins("JOIN esps ON esps.internal_id = boiling_sessions.esp_id").
		Where("esps.public_id = ? AND (boiling_sessions.boiling_status = ? OR boiling_sessions.boiling_status = ?)",
			espPublicID, "steaming", "preparing").
		First(&gormModel).Error

	if err != nil {
		return nil, err
	}

	base := gormModel.ToBase()
	base.EspPublicID = &espPublicID
	return &base, nil
}

func (r *BoSeRepositoryImpl) GetActiveSessions() ([]models.BoilingSessionBase, error) {
	var gormModels []models.BoilingSession

	err := config.DB.
		Where("boiling_status IN ?", []string{"preparing", "steaming"}).
		Find(&gormModels).Error

	if err != nil {
		return nil, err
	}

	result := make([]models.BoilingSessionBase, len(gormModels))
	for i, g := range gormModels {
		base := g.ToBase()

		if g.EspID != nil {
			var esp models.EspGorm
			if err := config.DB.
				Select("public_id").
				Where("internal_id = ?", *g.EspID).
				First(&esp).Error; err == nil {
				base.EspPublicID = &esp.PublicID
			}
		}

		result[i] = base
	}
	return result, nil
}
