package services

import (
	"hiliriset_ecoprint_golang/models"
	"hiliriset_ecoprint_golang/repositories"

	"github.com/google/uuid"
)

type KomporService interface{
	AddKompors(kompor *models.KomporRequestBody, userEmail string) error
	GetKompor(userEmail string) ([]models.Kompors, error)
}

type KomporServiceImpl struct {
	UserRepository repositories.UserRepository
	KomporRepository repositories.KomporRepository
} 

func NewKomporService(ur repositories.UserRepository,kr repositories.KomporRepository) KomporService {
	return &KomporServiceImpl{UserRepository: ur, KomporRepository: kr}
}

func (s *KomporServiceImpl) AddKompors(kompor *models.KomporRequestBody, userEmail string) error {
	existingUser, err := s.UserRepository.FindByEmail(userEmail)

	if err != nil {
		return err
	}

	newKompor := models.Kompors{
		PublicID: uuid.New(),
		KomporName: kompor.KomporName,
		UserId: int(existingUser.InternalID),
	}

	return s.KomporRepository.AddKompor(&newKompor)
}
func (s *KomporServiceImpl) GetKompor(userEmail string) ([]models.Kompors, error){
	exisistingUser, err := s.UserRepository.FindByEmail(userEmail)
	if err != nil {
		return nil, err
	}
	return s.KomporRepository.GetKompor(int(exisistingUser.InternalID))
}

