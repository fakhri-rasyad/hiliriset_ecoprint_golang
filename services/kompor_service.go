package services

import (
	"hiliriset_ecoprint_golang/models"
	"hiliriset_ecoprint_golang/repositories"

	"github.com/google/uuid"
)

type KomporService interface {
    AddKompor(req *models.KomporRequest, userEmail string) (*models.KomporBase, error)
    GetKompors(userEmail string) ([]models.KomporBase, error)
    GetKomporByPublicID(publicID uuid.UUID) (*models.KomporBase, error)
    DeleteKompor(publicID uuid.UUID) error
}

type KomporServiceImpl struct {
    ur repositories.UserRepository
    kr repositories.KomporRepository
}

func NewKomporService(ur repositories.UserRepository, kr repositories.KomporRepository) KomporService {
    return &KomporServiceImpl{ur, kr}
}

func (s *KomporServiceImpl) AddKompor(req *models.KomporRequest, userEmail string) (*models.KomporBase, error) {
    existingUser, err := s.ur.FindByEmail(userEmail)
    if err != nil {
        return nil, err
    }

    return s.kr.AddKompor(req, existingUser.InternalID)
}

func (s *KomporServiceImpl) GetKompors(userEmail string) ([]models.KomporBase, error) {
    existingUser, err := s.ur.FindByEmail(userEmail)
    if err != nil {
        return nil, err
    }

    return s.kr.GetKompors(existingUser.InternalID)
}

func (s *KomporServiceImpl) GetKomporByPublicID(publicID uuid.UUID) (*models.KomporBase, error) {
    return s.kr.GetKomporByPublicID(publicID)
}

func (s *KomporServiceImpl) DeleteKompor(publicID uuid.UUID) error {
    _, err := s.kr.GetKomporByPublicID(publicID)
    if err != nil {
        return err
    }

    return s.kr.DeleteKompor(publicID)
}