package controllers

import (
	"hiliriset_ecoprint_golang/models"
	"hiliriset_ecoprint_golang/services"
	"hiliriset_ecoprint_golang/utils"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type EspController interface {
    CreateEsp(ctx fiber.Ctx) error
    GetEsps(ctx fiber.Ctx) error
    GetEspByPublicID(ctx fiber.Ctx) error
    DeleteEsp(ctx fiber.Ctx) error
}

type EspControllerImpl struct {
    s services.EspService
}

func NewEspController(s services.EspService) EspController {
    return &EspControllerImpl{s: s}
}

// CreateEsp godoc
// @Summary      CreateEsp
// @Description  Endpoint untuk menambahkan esp baru
// @Tags         Esps
// @Accept       json
// @Produce      json
// @Param        esp body models.CreateEspRequest true "Raw json body"
// @Success      201  {object}  models.EspBase
// @Failure      400  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Security     ApiKeyAuth
// @Router       /api/v1/esps [post]
func (c *EspControllerImpl) CreateEsp(ctx fiber.Ctx) error {
    userEmail, err := utils.GetEmailClaim(ctx)
    if err != nil {
        return utils.Unauthorized(ctx, "User tidak ditemukan", err)
    }

    var req models.CreateEspRequest
    if err := ctx.Bind().Body(&req); err != nil {
        return utils.BadRequest(ctx, "Data tidak valid", err)
    }

    newEsp, err := c.s.AddEsp(userEmail, &req)
    if err != nil {
        return utils.InternalError(ctx, "Gagal menambahkan esp", err)
    }

    return utils.CreationSuccess(ctx, "Sukses menambahkan esp", newEsp)
}

// GetEsps godoc
// @Summary      GetEsps
// @Description  Endpoint untuk mengambil daftar esps pengguna
// @Tags         Esps
// @Produce      json
// @Success      200  {object}  []models.EspBase
// @Failure      400  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Security     ApiKeyAuth
// @Router       /api/v1/esps [get]
func (c *EspControllerImpl) GetEsps(ctx fiber.Ctx) error {
    userEmail, err := utils.GetEmailClaim(ctx)
    if err != nil {
        return utils.Unauthorized(ctx, "User tidak ditemukan", err)
    }

    esps, err := c.s.GetEsps(userEmail)
    if err != nil {
        return utils.BadRequest(ctx, "Gagal mengambil daftar esps", err)
    }

    return utils.SuccessResponse(ctx, "Sukses mengembalikan daftar esp", esps)
}

// GetEspByPublicID godoc
// @Summary      GetEspByPublicID
// @Description  Endpoint untuk mengambil detail esp menggunakan public id
// @Tags         Esps
// @Produce      json
// @Param        public_id path string true "Public ID ESP"
// @Success      200  {object}  models.EspBase
// @Failure      400  {object}  utils.Response
// @Failure      404  {object}  utils.Response
// @Security     ApiKeyAuth
// @Router       /api/v1/esps/{public_id} [get]
func (c *EspControllerImpl) GetEspByPublicID(ctx fiber.Ctx) error {
    publicID, err := uuid.Parse(ctx.Params("public_id"))
    if err != nil {
        return utils.BadRequest(ctx, "Format public_id tidak valid", err)
    }

    esp, err := c.s.GetEspByPublicID(publicID)
    if err != nil {
        return utils.NotFound(ctx, "ESP tidak ditemukan", err)
    }

    return utils.SuccessResponse(ctx, "Sukses mengambil data esp", esp)
}

// DeleteEsp godoc
// @Summary      DeleteEsp
// @Description  Endpoint untuk menghapus esp menggunakan public id
// @Tags         Esps
// @Produce      json
// @Param        public_id path string true "Public ID ESP"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Failure      404  {object}  utils.Response
// @Security     ApiKeyAuth
// @Router       /api/v1/esps/{public_id} [delete]
func (c *EspControllerImpl) DeleteEsp(ctx fiber.Ctx) error {
    publicID, err := uuid.Parse(ctx.Params("public_id"))
    if err != nil {
        return utils.BadRequest(ctx, "Format public_id tidak valid", err)
    }

    if err := c.s.DeleteEsp(publicID); err != nil {
        return utils.NotFound(ctx, "Gagal menghapus esp", err)
    }

    return utils.SuccessResponse(ctx, "Sukses menghapus esp", nil)
}