package services

import (
	"hiliriset_ecoprint_golang/models"
	"hiliriset_ecoprint_golang/repositories"

	"github.com/google/uuid"
)

type EspService interface {
    AddEsp(userEmail string, req *models.CreateEspRequest) (*models.EspBase, error)
    GetEsps(userEmail string) ([]models.EspBase, error)
    GetEspByPublicID(publicID uuid.UUID) (*models.EspBase, error)
    DeleteEsp(publicID uuid.UUID) error
}

type EspServiceImpl struct {
    er repositories.EspRepository
    ur repositories.UserRepository
}

func NewEspService(er repositories.EspRepository, ur repositories.UserRepository) EspService {
    return &EspServiceImpl{er, ur}
}

func (s *EspServiceImpl) AddEsp(userEmail string, req *models.CreateEspRequest) (*models.EspBase, error) {
    existingUser, err := s.ur.FindByEmail(userEmail)
    if err != nil {
        return nil, err
    }

    return s.er.AddEsp(req, existingUser.InternalID)
}

func (s *EspServiceImpl) GetEsps(userEmail string) ([]models.EspBase, error) {
    existingUser, err := s.ur.FindByEmail(userEmail)
    if err != nil {
        return nil, err
    }

    return s.er.GetEsps(existingUser.InternalID)
}

func (s *EspServiceImpl) GetEspByPublicID(publicID uuid.UUID) (*models.EspBase, error) {
    return s.er.GetEspByPublicID(publicID)
}

func (s *EspServiceImpl) DeleteEsp(publicID uuid.UUID) error {
    _, err := s.er.GetEspByPublicID(publicID)
    if err != nil {
        return err
    }

    return s.er.DeleteEsp(publicID)
}