package repositories

import (
	"hiliriset_ecoprint_golang/config"
	"hiliriset_ecoprint_golang/models"
)

type FabricTypeRepository interface {
    GetByName(name string) (*models.FabricType, error)
    GetById(id int64) (*models.FabricType, error)
}

type FabricTypeRepositoryImpl struct{}

func NewFabricTypeRepository() FabricTypeRepository {
    return &FabricTypeRepositoryImpl{}
}

func (r *FabricTypeRepositoryImpl) GetByName(name string) (*models.FabricType, error) {
    var ft models.FabricType
    if err := config.DB.Where("name = ?", name).First(&ft).Error; err != nil {
        return nil, err
    }
    return &ft, nil
}

func (r *FabricTypeRepositoryImpl) GetById(internalID int64) (*models.FabricType, error) {
    var ft models.FabricType
    if err := config.DB.Where("internal_id = ?", internalID).First(&ft).Error; err != nil {
        return nil, err
    }
    return &ft, nil
}
