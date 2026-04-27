package controllers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"

	"hiliriset_ecoprint_golang/services"
	"hiliriset_ecoprint_golang/utils"
)

type SessionRecordController interface {
    GetRecords(ctx fiber.Ctx) error
    GetRecordByPubID(ctx fiber.Ctx) error
}

type SessionRecordControllerImpl struct {
    s services.SessionRecordService
  }

func NewSessionRecordController(s services.SessionRecordService) SessionRecordController {
    return &SessionRecordControllerImpl{s: s}
}

// GetRecords godoc
// @Summary      GetRecords
// @Description  Mengambil semua records berdasarkan session
// @Tags         SessionRecords
// @Produce      json
// @Param        session_id path string true "Public ID session"
// @Success      200  {object}  []models.SessionRecordOutput
// @Failure      400  {object}  utils.Response
// @Failure      404  {object}  utils.Response
// @Security     ApiKeyAuth
// @Router       /api/v1/sessions/{session_id}/records [get]
func (c *SessionRecordControllerImpl) GetRecords(ctx fiber.Ctx) error {
    sessionPubID, err := uuid.Parse(ctx.Params("session_id"))
    if err != nil {
        return utils.BadRequest(ctx, "Format session_id tidak valid", err)
    }

    records, err := c.s.GetRecords(sessionPubID)
    if err != nil {
        return utils.NotFound(ctx, "Records tidak ditemukan", err)
    }

    return utils.SuccessResponse(ctx, "Sukses mengambil records", records)
}

// GetRecordByPubID godoc
// @Summary      GetRecordByPubID
// @Description  Mengambil detail satu record berdasarkan public id
// @Tags         SessionRecords
// @Produce      json
// @Param        record_id path string true "Public ID record"
// @Success      200  {object}  models.SessionRecordOutput
// @Failure      400  {object}  utils.Response
// @Failure      404  {object}  utils.Response
// @Security     ApiKeyAuth
// @Router       /api/v1/records/{record_id} [get]
func (c *SessionRecordControllerImpl) GetRecordByPubID(ctx fiber.Ctx) error {
    recordPubID, err := uuid.Parse(ctx.Params("record_id"))
    if err != nil {
        return utils.BadRequest(ctx, "Format record_id tidak valid", err)
    }

    record, err := c.s.GetRecordByPubID(recordPubID) // fixed: was GetRecordDetails
    if err != nil {
        return utils.NotFound(ctx, "Record tidak ditemukan", err)
    }

    return utils.SuccessResponse(ctx, "Sukses mengambil record", record)
}
