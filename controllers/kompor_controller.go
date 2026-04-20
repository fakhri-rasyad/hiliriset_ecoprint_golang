package controllers

import (
	"hiliriset_ecoprint_golang/models"
	"hiliriset_ecoprint_golang/services"
	"hiliriset_ecoprint_golang/utils"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type KomporController interface {
    AddKompor(ctx fiber.Ctx) error
    GetKompors(ctx fiber.Ctx) error
    GetKomporByPublicID(ctx fiber.Ctx) error
    DeleteKompor(ctx fiber.Ctx) error
}

type KomporControllerImpl struct {
    s services.KomporService
}

func NewKomporController(s services.KomporService) KomporController {
    return &KomporControllerImpl{s: s}
}

// AddKompor godoc
// @Summary      AddKompor
// @Description  Menambahkan kompor pengguna
// @Tags         Kompors
// @Accept       json
// @Produce      json
// @Param        kompor body models.KomporRequest true "Raw json kompor data"
// @Success      201  {object}  models.KomporBase
// @Failure      400  {object}  utils.Response
// @Security     ApiKeyAuth
// @Router       /api/v1/kompors [post]
func (c *KomporControllerImpl) AddKompor(ctx fiber.Ctx) error {
    userEmail, err := utils.GetEmailClaim(ctx)
    if err != nil {
        return utils.Unauthorized(ctx, "User tidak ditemukan", err)
    }

    var req models.KomporRequest
    if err := ctx.Bind().Body(&req); err != nil {
        return utils.BadRequest(ctx, "Gagal parsing data kompor", err)
    }

    newKompor, err := c.s.AddKompor(&req, userEmail)
    if err != nil {
        return utils.BadRequest(ctx, "Gagal menambahkan kompor", err)
    }

    return utils.CreationSuccess(ctx, "Sukses menambahkan kompor", newKompor)
}

// GetKompors godoc
// @Summary      GetKompors
// @Description  Mengembalikan daftar kompor pengguna
// @Tags         Kompors
// @Produce      json
// @Success      200  {object}  []models.KomporBase
// @Failure      400  {object}  utils.Response
// @Security     ApiKeyAuth
// @Router       /api/v1/kompors [get]
func (c *KomporControllerImpl) GetKompors(ctx fiber.Ctx) error {
    userEmail, err := utils.GetEmailClaim(ctx)
    if err != nil {
        return utils.Unauthorized(ctx, "User tidak ditemukan", err)
    }

    kompors, err := c.s.GetKompors(userEmail)
    if err != nil {
        return utils.BadRequest(ctx, "Gagal mengambil daftar kompor", err)
    }

    return utils.SuccessResponse(ctx, "Sukses mengambil daftar kompor", kompors)
}

// GetKomporByPublicID godoc
// @Summary      GetKomporByPublicID
// @Description  Mengambil detail kompor berdasarkan public id
// @Tags         Kompors
// @Produce      json
// @Param        public_id path string true "Public ID Kompor"
// @Success      200  {object}  models.KomporBase
// @Failure      400  {object}  utils.Response
// @Failure      404  {object}  utils.Response
// @Security     ApiKeyAuth
// @Router       /api/v1/kompors/{public_id} [get]
func (c *KomporControllerImpl) GetKomporByPublicID(ctx fiber.Ctx) error {
    publicID, err := uuid.Parse(ctx.Params("public_id"))
    if err != nil {
        return utils.BadRequest(ctx, "Format public_id tidak valid", err)
    }

    kompor, err := c.s.GetKomporByPublicID(publicID)
    if err != nil {
        return utils.NotFound(ctx, "Kompor tidak ditemukan", err)
    }

    return utils.SuccessResponse(ctx, "Sukses mengambil detail kompor", kompor)
}

// DeleteKompor godoc
// @Summary      DeleteKompor
// @Description  Menghapus kompor berdasarkan public id
// @Tags         Kompors
// @Produce      json
// @Param        public_id path string true "Public ID Kompor"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Failure      404  {object}  utils.Response
// @Security     ApiKeyAuth
// @Router       /api/v1/kompors/{public_id} [delete]
func (c *KomporControllerImpl) DeleteKompor(ctx fiber.Ctx) error {
    publicID, err := uuid.Parse(ctx.Params("public_id"))
    if err != nil {
        return utils.BadRequest(ctx, "Format public_id tidak valid", err)
    }

    if err := c.s.DeleteKompor(publicID); err != nil {
        return utils.NotFound(ctx, "Gagal menghapus kompor", err)
    }

    return utils.SuccessResponse(ctx, "Sukses menghapus kompor", nil)
}