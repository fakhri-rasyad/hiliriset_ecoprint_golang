package controllers

import (
	"hiliriset_ecoprint_golang/models"
	"hiliriset_ecoprint_golang/services"
	"hiliriset_ecoprint_golang/utils"

	"github.com/gofiber/fiber/v3"
)

type KomporController interface {
	AddKompor(ctx fiber.Ctx) error 
	GetKompors(ctx fiber.Ctx) error
}

type KomporControllerImpl struct {
	s services.KomporService

}

func NewKomporController(s services.KomporService) KomporController{
	return &KomporControllerImpl{s: s}
}


//ShowAccount godoc
//@Summary 		AddKompor
//@Description 	Menambahkan kompor pengguna
//@Tags 		Kompors
//@Accept		json
//@Produce		json
//@Param		kompor body models.KomporRequestBody true "Raw json kompor data"
//@Success		201 {object} models.KomporRequestBody 
//@Failure		400 {object} utils.Response
//@Security		ApiKeyAuth
//@Router		/api/v1/kompors [post]
func (c *KomporControllerImpl) AddKompor(ctx fiber.Ctx) error {
	userEmail, err := utils.GetEmailClaim(ctx)
	var KomporModel models.KomporRequestBody
	
	if err != nil {
		return err
	}

	if err := ctx.Bind().Body(&KomporModel); err != nil {
		return utils.BadRequest(ctx, "Gagal parsing data kompor", err)
	}
	
	if err := c.s.AddKompors(&KomporModel, userEmail); err != nil {
		return utils.BadRequest(ctx, "Gagal menambahkan kompor", err)
	}

	return utils.CreationSuccess(ctx, "Sukses menambahkan kompor", KomporModel)
}

//ShowAccount godoc
//@Summary 		GetKompors
//@Description 	Mengembalikan daftar kompor pengguna
//@Tags 		Kompors
//@Produce		json
//@Success		200 {object} []models.KomporRequestBody 
//@Failure		400 {object} utils.Response
//@Security		ApiKeyAuth
//@Router		/api/v1/kompors [get]
func (c *KomporControllerImpl) GetKompors(ctx fiber.Ctx) error {
	userEmail, err := utils.GetEmailClaim(ctx)
	
	if err != nil {
		return err
	}

	komporsData, err := c.s.GetKompor(userEmail)

	if err != nil {
		return utils.BadRequest(ctx, "Gagal mengambil daftar kompor", err)
	}

	return utils.SuccessResponse(ctx, "Sukses mengambil daftar kompor", komporsData)
}