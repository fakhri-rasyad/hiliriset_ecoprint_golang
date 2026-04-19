package controllers

import (
	"hiliriset_ecoprint_golang/models"
	"hiliriset_ecoprint_golang/services"
	"hiliriset_ecoprint_golang/utils"

	"github.com/gofiber/fiber/v3"
)

type EspController interface{
    CreateEsp(ctx fiber.Ctx) error
    GetEspDetail(ctx fiber.Ctx) error
    GetEsps(ctx fiber.Ctx) error
}

type EspControllerImpl struct{
	s services.EspService
}

func NewEspController(s services.EspService) EspController{
	return &EspControllerImpl{s: s}
}

//CreateEsp godoc
//@Summary			CreateEsp
//@Description		Endpoint untuk menambahkan esp baru
//@Tags				Esps
//@Security 		ApiKeyAuth	
//@Param        	esp body models.CreateEspRequest true "Raw json body"
//@Accept			json
//@Produce			json
//@Success			201 {object} models.EspBase
//@Failure			400 {object} utils.Response
//@Failure			500 {object} utils.Response
//@Router			/api/v1/esps [post]
func (c *EspControllerImpl) CreateEsp(ctx fiber.Ctx) error {
	userEmail, err := utils.GetEmailClaim(ctx)
	if err != nil {
		return utils.BadRequest(ctx, "Gagal menambahkan esp: User tidak ditemukan", err)
	}
	var espData models.CreateEspRequest

	if err := ctx.Bind().Body(&espData); err != nil {
		return utils.BadRequest(ctx, "Gagal menambahkan esp: Data tidak valid", err)
	}
	
	addedEsp, err := c.s.AddNewEsps(userEmail, &espData)

	if err != nil {
		return utils.InternalError(ctx, "Gagal menambahkan esps: internal error", err)
	}

	return utils.CreationSuccess(ctx, "Sukses menambahkan esp", addedEsp)
}

//GetEspDetail godoc
//@Summary 		GetEspDetail
//@Description 	Endpoint untuk mengambil detail esps menggunakan uuid esp
//@Tags 		Esps
//@Produce		json
//@Param		esp_pub_id path string true "Public Id ESP"
//@Success		200 {object} models.EspBase		
//@Failure		404 {object} utils.Response
//@Failure		500 {object} utils.Response
//@Router		/api/v1/esps/{esp_pub_id} [get]
//@Security		ApiKeyAuth
func (c *EspControllerImpl) GetEspDetail(ctx fiber.Ctx) error {
	espPubID := ctx.Params("esp_pub_id")
	espData, err := c.s.GetEspDetail(espPubID)

	if err != nil {
		return utils.NotFound(ctx, "Gagal mengambil detail esp", err)
	}

	return utils.SuccessResponse(ctx, "Sukses mengambil data esp", espData)
}

//GetEsps godoc
//@Summary 		GetEsps
//@Description 	Endpoint untuk mengambil daftar esps pengguna menggunakan email pengguna
//@Tags 		Esps
//@Produce		json
//@Success		200 {object} []models.EspBase		
//@Failure		400 {object} utils.Response
//@Failure		500 {object} utils.Response
//@Router		/api/v1/esps [get]
//@Security		ApiKeyAuth
func (c *EspControllerImpl) GetEsps(ctx fiber.Ctx) error {
	userEmail, err := utils.GetEmailClaim(ctx)
	if err != nil {
		return utils.BadRequest(ctx, "Gagal mengambil esp: User tidak ditemukan", err)
	}

	esps, err := c.s.GetEsps(userEmail)

	if err != nil {
		return utils.BadRequest(ctx, "Gagal mengambil daftar esps", err)
	}

	return utils.SuccessResponse(ctx, "Sukses mengembalikan daftar esp", esps)
}