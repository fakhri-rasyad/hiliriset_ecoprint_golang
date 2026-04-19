package repositories

import (
	"hiliriset_ecoprint_golang/config"
	"hiliriset_ecoprint_golang/models"
)

type EspRepository interface{
	AddEsp(esps *models.EspGorm) error
	GetEsps(userId int64) ([]models.EspBase, error)
	GetEspDetail(espPublicID string)(*models.EspBase, error)
}

type EspRepositoryImpl struct {
}

func NewEspRepository() EspRepository {
	return &EspRepositoryImpl{}
}

func (r *EspRepositoryImpl) AddEsp(esps *models.EspGorm) error {
	return config.DB.Create(esps).Error
}

func (r *EspRepositoryImpl) GetEsps(userID int64) ([]models.EspBase ,error){
	var esps []models.EspGorm
	if err := config.DB.Where("user_id = ?", userID).Find(&esps).Error; err != nil{
		return nil, err
	}

	baseEsps := make([]models.EspBase, len(esps))
	for i, base := range esps{
		baseEsps[i] = base.EspBase
	}

	return baseEsps, nil
}

func (r *EspRepositoryImpl) GetEspDetail(espPublicId string) (*models.EspBase, error){
	var esp models.EspGorm
	if err := config.DB.Where("public_id = ?", espPublicId).First(&esp).Error; err != nil {
		return nil, err
	}
	return &esp.EspBase, nil
}

//Selesaikan