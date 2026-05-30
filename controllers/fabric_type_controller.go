package controllers

import (
	"hiliriset_ecoprint_golang/services"
	"hiliriset_ecoprint_golang/utils"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type FabricTypeController interface {
	GetFabricByUUID(ctx fiber.Ctx) error
  GetFabricAll(ctx fiber.Ctx) error
}

func NewFabricTypeController(s services.FabricTypeService) FabricTypeController {
	return &FabricTypeControllerImpl{s: s}
}

type FabricTypeControllerImpl struct {
	s services.FabricTypeService
}

// GetFabricByUUID godoc
// @Summary     GetFabricByUUID
// @Description Mengambil informasi kain berdasarkan id public
// @Tags        Fabrics
// @Produce     json
// @Param       fabric_uuid path string true "Public id kain"
// @Success     200 {object} models.FabricType
// @Failure     400 {object} utils.Response
// @Error       404 {object} utils.Response
// @Security    ApiKeyAuth
// @Router      /api/v1/fabric_types/{fabric_uuid} [get]
func (f *FabricTypeControllerImpl) GetFabricByUUID(ctx fiber.Ctx) error {
  fabric_uuid, err := uuid.Parse(ctx.Params("fabric_uuid"))
  if err != nil {
    return utils.InternalError(ctx, "Kesalahan saat parsing uuid", err)
  }
  fabric, err := f.s.GetFabricByUUID(fabric_uuid)
  if err != nil {
    return utils.NotFound(ctx, "Kain tidak ditemukan", err)
  }

  return utils.SuccessResponse(ctx, "Fabric berhasil diambil", fabric)

}

// GetFabricByAll godoc
// @Summary     GetFabricAll
// @Description Mengambil daftar kain
// @Tags        Fabrics
// @Produce     json
// @Success     200 {object} []models.FabricType
// @Failure     400 {object} utils.Response
// @Router      /api/v1/fabric_types [get]
func (f *FabricTypeControllerImpl) GetFabricAll(ctx fiber.Ctx) error {
  fabrics, err := f.s.GetAllFabric()
  if err != nil {
    return utils.BadRequest(ctx, "Kesalahan saat mengambil daftar kain", err)
  }

  return utils.SuccessResponse(ctx, "Daftar kain sukses dikembalikan", fabrics)
}

