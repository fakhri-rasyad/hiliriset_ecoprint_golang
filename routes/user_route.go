package routes

import (
	"hiliriset_ecoprint_golang/config"
	"hiliriset_ecoprint_golang/controllers"
	"hiliriset_ecoprint_golang/utils"
	"log"

	jwtware "github.com/gofiber/contrib/v3/jwt"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/extractors"
	"github.com/joho/godotenv"
)

func Setup(
	app *fiber.App,
	userController controllers.UserController,
	komporController controllers.KomporController,
	espController controllers.EspController,
) {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Gagal membuka file .env")
	}

	app.Post("v1/auth/register", userController.RegisterUser)
	app.Post("v1/auth/login", userController.LoginUser)

	

	api := app.Group("/api/v1")
	api.Use(
		jwtware.New(
			jwtware.Config{
				SigningKey: jwtware.SigningKey{Key: []byte(config.APPConfig.JWTSecret)},
				Extractor: extractors.FromAuthHeader("Bearer"),
				ErrorHandler: func (c fiber.Ctx, err error) error {
					return utils.UnauthorizedReponse(c, "User unathorized", err)
				},
	}))

	api.Get("/kompors", komporController.GetKompors)
	api.Post("/kompors", komporController.AddKompor)

	api.Get("/esps", espController.GetEsps)
	api.Get("/esps/:esp_pub_id", espController.GetEspDetail)
	api.Post("/esps", espController.CreateEsp)
	
}