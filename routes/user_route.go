package routes

import (
	"hiliriset_ecoprint_golang/controllers"
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
)

func Setup(
	app *fiber.App,
	userController controllers.UserController,
) {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Gagal membuka file .env")
	}

	app.Post("v1/auth/register", userController.RegisterUser)
}