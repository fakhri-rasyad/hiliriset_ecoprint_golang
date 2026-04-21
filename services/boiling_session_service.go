package services

import (
	"hiliriset_ecoprint_golang/models"
	"hiliriset_ecoprint_golang/repositories"

	"github.com/google/uuid"
)

type BoSeService interface {
    CreateSession(userEmail string, req *models.BoilingSessionCreation) (*models.BoilingSessionResponse, error)
    GetSessionByPublicID(publicID uuid.UUID) (*models.BoilingSessionResponse, error)
    GetSessions(userEmail string) ([]models.BoilingSessionResponse, error)
    UpdateSessionStatus(publicID uuid.UUID, status string) error
    FinishSession(publicID uuid.UUID) error
}

type BoSeServiceImpl struct {
    bsRepo   repositories.BoSeRepository
    userRepo repositories.UserRepository
    krRepo   repositories.KomporRepository
    espRepo  repositories.EspRepository
}

func NewBoSeService(
    bsRepo repositories.BoSeRepository,
    userRepo repositories.UserRepository,
    krRepo repositories.KomporRepository,
    espRepo repositories.EspRepository,
) BoSeService {
    return &BoSeServiceImpl{
        bsRepo:   bsRepo,
        userRepo: userRepo,
        krRepo:   krRepo,
        espRepo:  espRepo,
    }
}

func (s *BoSeServiceImpl) CreateSession(userEmail string, req *models.BoilingSessionCreation) (*models.BoilingSessionResponse, error) {
    user, err := s.userRepo.FindByEmail(userEmail)
    if err != nil {
        return nil, err
    }

    // Resolve public IDs to internal IDs
    kompor, err := s.krRepo.GetKomporByPublicID(req.KomporPublicID)
    if err != nil {
        return nil, err
    }

    esp, err := s.espRepo.GetEspByPublicID(req.EspPublicID)
    if err != nil {
        return nil, err
    }

    base, err := s.bsRepo.CreateSession(user.InternalID, kompor.InternalID, esp.InternalID, req.FabricType)
    if err != nil {
        return nil, err
    }

    response := base.ToResponse()
    return &response, nil
}

func (s *BoSeServiceImpl) GetSessionByPublicID(publicID uuid.UUID) (*models.BoilingSessionResponse, error) {
    base, err := s.bsRepo.GetSessionByPublicID(publicID)
    if err != nil {
        return nil, err
    }

    response := base.ToResponse()
    return &response, nil
}

func (s *BoSeServiceImpl) GetSessions(userEmail string) ([]models.BoilingSessionResponse, error) {
    user, err := s.userRepo.FindByEmail(userEmail)
    if err != nil {
        return nil, err
    }

    bases, err := s.bsRepo.GetSessions(user.InternalID)
    if err != nil {
        return nil, err
    }

    responses := make([]models.BoilingSessionResponse, len(bases))
    for i, base := range bases {
        responses[i] = base.ToResponse()
    }

    return responses, nil
}

func (s *BoSeServiceImpl) UpdateSessionStatus(publicID uuid.UUID, status string) error {
    return s.bsRepo.UpdateSessionStatus(publicID, status)
}

func (s *BoSeServiceImpl) FinishSession(publicID uuid.UUID) error {
    return s.bsRepo.FinishSession(publicID)
}