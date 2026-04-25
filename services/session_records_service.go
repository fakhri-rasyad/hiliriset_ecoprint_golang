package services

import (
	"hiliriset_ecoprint_golang/models"
	"hiliriset_ecoprint_golang/repositories"
	"sync"

	"github.com/google/uuid"
)

type SessionRecordService interface {
    CreateRecord(sessionPubID uuid.UUID, data models.SessionRecordInput) error
    GetRecordByPubID(recordPubID uuid.UUID) (*models.SessionRecordOutput, error)
    GetRecords(sessionPubID uuid.UUID) ([]models.SessionRecordOutput, error)
    InvalidateSession(sessionPubID uuid.UUID)
}

type SessionRecordServiceImpl struct {
    srRepo       repositories.SeReRepository
    bsRepo       repositories.BoSeRepository
    sessionCache sync.Map
}

func NewSessionRecordService(
    srRepo repositories.SeReRepository,
    bsRepo repositories.BoSeRepository,
) SessionRecordService {
    return &SessionRecordServiceImpl{
        srRepo: srRepo,
        bsRepo: bsRepo,
    }
}

func (s *SessionRecordServiceImpl) resolveSessionID(sessionPubID uuid.UUID) (int64, error) {
    if cached, ok := s.sessionCache.Load(sessionPubID); ok {
        return cached.(int64), nil
    }

    existingSession, err := s.bsRepo.GetSessionByPublicID(sessionPubID)
    if err != nil {
        return 0, err
    }

    s.sessionCache.Store(sessionPubID, existingSession.InternalID)
    return existingSession.InternalID, nil
}

func (s *SessionRecordServiceImpl) CreateRecord(sessionPubID uuid.UUID, data models.SessionRecordInput) error {
    sessionID, err := s.resolveSessionID(sessionPubID)
    if err != nil {
        return err
    }

    record := &models.SessionRecordGorm{
        PublicID:  uuid.New(),
        SessionID: sessionID,
        AirTemp:   data.AirTemp,
        WaterTemp: data.WaterTemp,
        Humidity:  data.Humidity,
    }

    return s.srRepo.CreateRecord(record)
}

func (s *SessionRecordServiceImpl) GetRecordByPubID(recordPubID uuid.UUID) (*models.SessionRecordOutput, error) {
    base, err := s.srRepo.GetRecordByPubID(recordPubID)
    if err != nil {
        return nil, err
    }

    output := base.ToOutput()
    return &output, nil
}

func (s *SessionRecordServiceImpl) GetRecords(sessionPubID uuid.UUID) ([]models.SessionRecordOutput, error) {
    sessionID, err := s.resolveSessionID(sessionPubID)
    if err != nil {
        return nil, err
    }

    bases, err := s.srRepo.GetRecords(sessionID)
    if err != nil {
        return nil, err
    }

    outputs := make([]models.SessionRecordOutput, len(bases))
    for i, base := range bases {
        outputs[i] = base.ToOutput()
    }

    return outputs, nil
}

func (s *SessionRecordServiceImpl) InvalidateSession(sessionPubID uuid.UUID) {
    s.sessionCache.Delete(sessionPubID)
}
