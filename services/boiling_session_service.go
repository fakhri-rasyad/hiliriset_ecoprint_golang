package services

import (
	"encoding/json"
	"errors"
	"hiliriset_ecoprint_golang/models"
	"hiliriset_ecoprint_golang/repositories"
	"log"

	"github.com/google/uuid"
)

type BoSeService interface {
	CreateSession(userEmail string, req *models.BoilingSessionCreation) (*models.BoilingSessionResponse, error)
	GetSessionByPublicID(publicID uuid.UUID) (*models.BoilingSessionResponse, error)
	GetSessions(userEmail string) ([]models.BoilingSessionResponse, error)
	UpdateSessionStatus(publicID uuid.UUID, status string) error
	FinishSession(publicID uuid.UUID) error
}

type MQTTPublisher interface {
	Publish(topic string, payload string) error
}

type SessionEspRegistrar interface {
    RegisterEspSession(espPublicID uuid.UUID)
}

type BoSeServiceImpl struct {
	bsRepo   repositories.BoSeRepository
	userRepo repositories.UserRepository
	krRepo   repositories.KomporRepository
	espRepo  repositories.EspRepository
  fabricTypeRepo repositories.FabricTypeRepository
	mqttPub  MQTTPublisher
  espRegistrar SessionEspRegistrar

}

func NewBoSeService(
	bsRepo repositories.BoSeRepository,
	userRepo repositories.UserRepository,
	krRepo repositories.KomporRepository,
	espRepo repositories.EspRepository,
  fabricTypeRepo repositories.FabricTypeRepository,
	mqttPub MQTTPublisher,
  espRegistrar SessionEspRegistrar,
) BoSeService {
	return &BoSeServiceImpl{
		bsRepo:   bsRepo,
		userRepo: userRepo,
		krRepo:   krRepo,
		espRepo:  espRepo,
    fabricTypeRepo: fabricTypeRepo,
		mqttPub:  mqttPub,
    espRegistrar: espRegistrar,
	}
}

func (s *BoSeServiceImpl) CreateSession(userEmail string, req *models.BoilingSessionCreation) (*models.BoilingSessionResponse, error) {
	user, err := s.userRepo.FindByEmail(userEmail)
	if err != nil {
		return nil, err
	}

	kompor, err := s.krRepo.GetKomporByPublicID(req.KomporPublicID)
	if err != nil {
		return nil, err
	}

	esp, err := s.espRepo.GetEspByPublicID(req.EspPublicID)
	if err != nil {
		return nil, err
	}


	if kompor.IsActive {
		return nil, errors.New("kompor sedang digunakan di sesi lain")
	}
	if esp.IsActive {
		return nil, errors.New("esp sedang digunakan di sesi lain")
	}

  fabric, err := s.fabricTypeRepo.GetByName(req.FabricType)
  if err != nil {
    return nil, errors.New("fabric type tidak ditemukan: " + req.FabricType)
  }

	base, err := s.bsRepo.CreateSession(
		user.InternalID,
		kompor.InternalID,
		esp.InternalID,
		esp.PublicID,
		fabric.InternalID,
	)
	if err != nil {
		return nil, err
	}

	if err := s.krRepo.SetActive(req.KomporPublicID, true); err != nil {
		log.Printf("failed to set kompor active: %v", err)
	}
	if err := s.espRepo.SetActive(req.EspPublicID, true); err != nil {
		log.Printf("failed to set esp active: %v", err)
	}

	startCommand := &models.MQTTCommandPayload{Command: "Start"}
	payload, err := json.Marshal(startCommand)
	if err != nil {
		log.Printf("failed to marshal start command: %v", err)
	} else {
		topic := "esp/" + esp.PublicID.String() + "/command"
		if err := s.mqttPub.Publish(topic, string(payload)); err != nil {
			log.Printf("failed to publish start to esp %s: %v", esp.PublicID, err)
		}
	}

  s.espRegistrar.RegisterEspSession(esp.PublicID)

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
