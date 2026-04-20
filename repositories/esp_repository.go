package repositories

import (
	"hiliriset_ecoprint_golang/config"
	"hiliriset_ecoprint_golang/models"

	"github.com/google/uuid"
)

type EspRepository interface {
    AddEsp(req *models.CreateEspRequest, userID int64) (*models.EspBase, error)
    GetEsps(userID int64) ([]models.EspBase, error)
    GetEspByPublicID(publicID uuid.UUID) (*models.EspBase, error)
    DeleteEsp(publicID uuid.UUID) error
}

type EspRepositoryImpl struct{}

func NewEspRepository() EspRepository {
    return &EspRepositoryImpl{}
}

func (r *EspRepositoryImpl) AddEsp(req *models.CreateEspRequest, userID int64) (*models.EspBase, error) {
    gormModel := models.EspGorm{
        MacAddress: req.MacAddress,
        UserID:     &userID,
		PublicID: uuid.New(),
    }

    if err := config.DB.Create(&gormModel).Error; err != nil {
        return nil, err
    }

    base := gormModel.ToBase()
    return &base, nil
}

func (r *EspRepositoryImpl) GetEsps(userID int64) ([]models.EspBase, error) {
    var gormModels []models.EspGorm

    if err := config.DB.Where("user_id = ?", userID).Find(&gormModels).Error; err != nil {
        return nil, err
    }

    result := make([]models.EspBase, len(gormModels))
    for i, g := range gormModels {
        result[i] = g.ToBase()
    }

    return result, nil
}

func (r *EspRepositoryImpl) GetEspByPublicID(publicID uuid.UUID) (*models.EspBase, error) {
    var gormModel models.EspGorm

    if err := config.DB.Where("public_id = ?", publicID).First(&gormModel).Error; err != nil {
        return nil, err
    }

    base := gormModel.ToBase()
    return &base, nil
}

func (r *EspRepositoryImpl) DeleteEsp(publicID uuid.UUID) error {
    return config.DB.
        Where("public_id = ?", publicID).
        Delete(&models.EspGorm{}).
        Error
}