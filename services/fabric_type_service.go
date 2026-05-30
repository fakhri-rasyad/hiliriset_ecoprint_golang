package services

import (
	"hiliriset_ecoprint_golang/models"
	"hiliriset_ecoprint_golang/repositories"

	"github.com/google/uuid"
)

type FabricTypeService interface {
	GetFabricByUUID(uuid uuid.UUID) (*models.FabricType, error)
  GetAllFabric()([]models.FabricType, error)
}

type FabricTypeServiceImpl struct {
	fabricTypeRepo repositories.FabricTypeRepository
}


func NewFabricTypeService(fabricTypeRepo repositories.FabricTypeRepository) FabricTypeService {
	return &FabricTypeServiceImpl{fabricTypeRepo: fabricTypeRepo}
}

func (f *FabricTypeServiceImpl) GetFabricByUUID(uuid uuid.UUID) (*models.FabricType, error) {
	return f.fabricTypeRepo.GetByUUID(uuid)
}

func (f *FabricTypeServiceImpl) GetAllFabric()([]models.FabricType, error){
  return f.fabricTypeRepo.GetAllFabric()
}

