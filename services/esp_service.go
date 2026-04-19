package services

import (
	"hiliriset_ecoprint_golang/models"
	"hiliriset_ecoprint_golang/repositories"

	"github.com/google/uuid"
)

type EspService interface{
	AddNewEsps(userEmail string, espData *models.CreateEspRequest) (*models.EspBase, error)
	GetEspDetail(espPubID string) (*models.EspBase, error)
	GetEsps(userEmail string)([]models.EspBase, error)
}

type EspServiceImpl struct {
	er repositories.EspRepository
	ur repositories.UserRepository
}

func NewEspService(er repositories.EspRepository, ur repositories.UserRepository ) EspService{
	return &EspServiceImpl{er, ur}
}

func (s *EspServiceImpl) AddNewEsps(userEmail string, espData *models.CreateEspRequest) (*models.EspBase, error) {
	existingUser, err := s.ur.FindByEmail(userEmail)
	if err != nil {
		return nil, err
	} 
	
	newEsp := &models.EspGorm{
		UserID: existingUser.InternalID,
		EspBase: models.EspBase{
			PublicID: uuid.New(),
			MacAddress: espData.MacAddress,
			DeviceStatus: "offline",
		},
	} 

	if err := s.er.AddEsp(newEsp); err != nil {
		return nil, err
	}

	return &newEsp.EspBase, nil
}

func (s *EspServiceImpl) GetEspDetail(espPubID string) (*models.EspBase, error) {
	espData, err := s.er.GetEspDetail(espPubID)

	if err != nil {
		return nil, err
	}

	return espData, nil
}

func (s *EspServiceImpl) GetEsps(userEmail string)([]models.EspBase, error){
	existingUser, err := s.ur.FindByEmail(userEmail)

	if err != nil {
		return nil, err
	}

	esps, err := s.er.GetEsps(existingUser.InternalID)

	if err != nil {
		return nil, err
	}

	return esps, nil

}