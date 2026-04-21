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
    boSeController controllers.BoSeController,
) {
    err := godotenv.Load(".env")
    if err != nil {
        log.Fatal("Gagal membuka file .env")
    }

    auth := app.Group("/v1/auth")
    auth.Post("/register", userController.RegisterUser)
    auth.Post("/login", userController.LoginUser)

    api := app.Group("/api/v1")
    api.Use(jwtware.New(jwtware.Config{
        SigningKey: jwtware.SigningKey{Key: []byte(config.APPConfig.JWTSecret)},
        Extractor: extractors.FromAuthHeader("Bearer"),
        ErrorHandler: func(c fiber.Ctx, err error) error {
            return utils.UnauthorizedReponse(c, "User unauthorized", err)
        },
    }))

    api.Get("/kompors", komporController.GetKompors)
    api.Post("/kompors", komporController.AddKompor)
    api.Get("/kompors/:public_id", komporController.GetKomporByPublicID)
    api.Delete("/kompors/:public_id", komporController.DeleteKompor)

    api.Get("/esps", espController.GetEsps)
    api.Post("/esps", espController.CreateEsp)
    api.Get("/esps/:public_id", espController.GetEspByPublicID)
    api.Delete("/esps/:public_id", espController.DeleteEsp)

    api.Post("/sessions", boSeController.CreateSession)
    api.Get("/sessions", boSeController.GetSessions)
    api.Get("/sessions/:public_id", boSeController.GetSessionByPublicID)
    api.Patch("/sessions/:public_id/status", boSeController.UpdateSessionStatus)
    api.Patch("/sessions/:public_id/finish", boSeController.FinishSession)
}