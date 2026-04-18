package repositories

import (
	"hiliriset_ecoprint_golang/config"
	"hiliriset_ecoprint_golang/models"
)

type KomporRepository interface {
	GetKompor(user_id int) ([]models.Kompors, error)
	AddKompor(kompor *models.Kompors) error
}

type KomporRepositoryImpl struct {
}

func NewKomporRepository() KomporRepository {
	return &KomporRepositoryImpl{}
}

func (r *KomporRepositoryImpl) GetKompor(userID int) ([]models.Kompors, error){
	var kompors []models.Kompors

	if err := config.DB.Where("user_id = ?", userID).Find(&kompors).Error; err != nil {
		return nil, err
	}

	return kompors, nil
}


func (r *KomporRepositoryImpl) AddKompor(kompor *models.Kompors) error{
	return config.DB.Create(kompor).Error
}