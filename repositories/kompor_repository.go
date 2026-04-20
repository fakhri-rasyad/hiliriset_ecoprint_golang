package repositories

import (
	"hiliriset_ecoprint_golang/config"
	"hiliriset_ecoprint_golang/models"

	"github.com/google/uuid"
)

type KomporRepository interface {
    GetKompors(userID int64) ([]models.KomporBase, error)
    GetKomporByPublicID(publicID uuid.UUID) (*models.KomporBase, error)
    AddKompor(req *models.KomporRequest, userID int64)(*models.KomporBase, error)
    DeleteKompor(publicID uuid.UUID) error
}

type KomporRepositoryImpl struct{}

func NewKomporRepository() KomporRepository {
    return &KomporRepositoryImpl{}
}

func (r *KomporRepositoryImpl) GetKompors(userID int64) ([]models.KomporBase, error) {
    var gormModels []models.KomporGorm

    if err := config.DB.Where("user_id = ?", userID).Find(&gormModels).Error; err != nil {
        return nil, err
    }

    result := make([]models.KomporBase, len(gormModels))
    for i, g := range gormModels {
        result[i] = g.ToBase()
    }

    return result, nil
}

func (r *KomporRepositoryImpl) GetKomporByPublicID(publicID uuid.UUID) (*models.KomporBase, error) {
    var gormModel models.KomporGorm

    if err := config.DB.Where("public_id = ?", publicID).First(&gormModel).Error; err != nil {
        return nil, err
    }

    base := gormModel.ToBase()
    return &base, nil
}

func (r *KomporRepositoryImpl) AddKompor(req *models.KomporRequest, userID int64) (*models.KomporBase, error) {
    gormModel := models.KomporGorm{
        KomporName: req.KomporName,
        UserID:     &userID,
		PublicID: uuid.New(),
    }

	if err := config.DB.Create(&gormModel).Error; err != nil {
		return nil, err
	}

	baseModel := gormModel.ToBase()

    return &baseModel, nil
}

func (r *KomporRepositoryImpl) DeleteKompor(publicID uuid.UUID) error {
    return config.DB.
        Where("public_id = ?", publicID).
        Delete(&models.KomporGorm{}).
        Error
}