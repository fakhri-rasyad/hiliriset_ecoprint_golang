package controllers

import (
	"hiliriset_ecoprint_golang/models"
	"hiliriset_ecoprint_golang/services"
	"hiliriset_ecoprint_golang/utils"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type BoSeController interface {
    CreateSession(ctx fiber.Ctx) error
    GetSessions(ctx fiber.Ctx) error
    GetSessionByPublicID(ctx fiber.Ctx) error
    UpdateSessionStatus(ctx fiber.Ctx) error
    FinishSession(ctx fiber.Ctx) error
}

type BoSeControllerImpl struct {
    s services.BoSeService
}

func NewBoSeController(s services.BoSeService) BoSeController {
    return &BoSeControllerImpl{s: s}
}

// CreateSession godoc
// @Summary      CreateSession
// @Description  Membuat sesi perebusan baru
// @Tags         BoisingSessions
// @Accept       json
// @Produce      json
// @Param        session body models.BoilingSessionCreation true "Data sesi perebusan"
// @Success      201  {object}  models.BoilingSessionResponse
// @Failure      400  {object}  utils.Response
// @Failure      401  {object}  utils.Response
// @Security     ApiKeyAuth
// @Router       /api/v1/sessions [post]
func (c *BoSeControllerImpl) CreateSession(ctx fiber.Ctx) error {
    userEmail, err := utils.GetEmailClaim(ctx)
    if err != nil {
        return utils.Unauthorized(ctx, "User tidak ditemukan", err)
    }

    var req models.BoilingSessionCreation
    if err := ctx.Bind().Body(&req); err != nil {
        return utils.BadRequest(ctx, "Gagal parsing data sesi", err)
    }

    newSession, err := c.s.CreateSession(userEmail, &req)
    if err != nil {
        return utils.BadRequest(ctx, "Gagal membuat sesi perebusan", err)
    }

    return utils.CreationSuccess(ctx, "Sukses membuat sesi perebusan", newSession)
}

// GetSessions godoc
// @Summary      GetSessions
// @Description  Mengembalikan daftar sesi perebusan pengguna
// @Tags         BoisingSessions
// @Produce      json
// @Success      200  {object}  []models.BoilingSessionResponse
// @Failure      400  {object}  utils.Response
// @Failure      401  {object}  utils.Response
// @Security     ApiKeyAuth
// @Router       /api/v1/sessions [get]
func (c *BoSeControllerImpl) GetSessions(ctx fiber.Ctx) error {
    userEmail, err := utils.GetEmailClaim(ctx)
    if err != nil {
        return utils.Unauthorized(ctx, "User tidak ditemukan", err)
    }

    sessions, err := c.s.GetSessions(userEmail)
    if err != nil {
        return utils.BadRequest(ctx, "Gagal mengambil daftar sesi", err)
    }

    return utils.SuccessResponse(ctx, "Sukses mengambil daftar sesi", sessions)
}

// GetSessionByPublicID godoc
// @Summary      GetSessionByPublicID
// @Description  Mengambil detail sesi perebusan berdasarkan public id
// @Tags         BoisingSessions
// @Produce      json
// @Param        public_id path string true "Public ID Sesi"
// @Success      200  {object}  models.BoilingSessionResponse
// @Failure      400  {object}  utils.Response
// @Failure      404  {object}  utils.Response
// @Security     ApiKeyAuth
// @Router       /api/v1/sessions/{public_id} [get]
func (c *BoSeControllerImpl) GetSessionByPublicID(ctx fiber.Ctx) error {
    publicID, err := uuid.Parse(ctx.Params("public_id"))
    if err != nil {
        return utils.BadRequest(ctx, "Format public_id tidak valid", err)
    }

    session, err := c.s.GetSessionByPublicID(publicID)
    if err != nil {
        return utils.NotFound(ctx, "Sesi tidak ditemukan", err)
    }

    return utils.SuccessResponse(ctx, "Sukses mengambil detail sesi", session)
}

// UpdateSessionStatus godoc
// @Summary      UpdateSessionStatus
// @Description  Memperbarui status sesi perebusan
// @Tags         BoisingSessions
// @Accept       json
// @Produce      json
// @Param        public_id path string true "Public ID Sesi"
// @Param        status body models.BoilingSessionStatusUpdate true "Status baru sesi"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Failure      404  {object}  utils.Response
// @Security     ApiKeyAuth
// @Router       /api/v1/sessions/{public_id}/status [patch]
func (c *BoSeControllerImpl) UpdateSessionStatus(ctx fiber.Ctx) error {
    publicID, err := uuid.Parse(ctx.Params("public_id"))
    if err != nil {
        return utils.BadRequest(ctx, "Format public_id tidak valid", err)
    }

    var req models.BoilingSessionStatusUpdate
    if err := ctx.Bind().Body(&req); err != nil {
        return utils.BadRequest(ctx, "Gagal parsing status", err)
    }

    if err := c.s.UpdateSessionStatus(publicID, req.Status); err != nil {
        return utils.NotFound(ctx, "Gagal memperbarui status sesi", err)
    }

    return utils.SuccessResponse(ctx, "Sukses memperbarui status sesi", nil)
}

// FinishSession godoc
// @Summary      FinishSession
// @Description  Menyelesaikan sesi perebusan
// @Tags         BoisingSessions
// @Produce      json
// @Param        public_id path string true "Public ID Sesi"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Failure      404  {object}  utils.Response
// @Security     ApiKeyAuth
// @Router       /api/v1/sessions/{public_id}/finish [patch]
func (c *BoSeControllerImpl) FinishSession(ctx fiber.Ctx) error {
    publicID, err := uuid.Parse(ctx.Params("public_id"))
    if err != nil {
        return utils.BadRequest(ctx, "Format public_id tidak valid", err)
    }

    if err := c.s.FinishSession(publicID); err != nil {
        return utils.NotFound(ctx, "Gagal menyelesaikan sesi", err)
    }

    return utils.SuccessResponse(ctx, "Sukses menyelesaikan sesi", nil)
}