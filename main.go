package main

import (
	"hiliriset_ecoprint_golang/config"
	"hiliriset_ecoprint_golang/controllers"
	"hiliriset_ecoprint_golang/repositories"
	"hiliriset_ecoprint_golang/routes"
	"hiliriset_ecoprint_golang/services"
	"log"

	"github.com/gofiber/fiber/v3"

	_ "hiliriset_ecoprint_golang/docs"

	swaggo "github.com/gofiber/contrib/v3/swaggo"
)

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:3000
// @BasePath  /

// @securityDefinitions.apiKey  ApiKeyAuth
// @in							header
// @name						Authorization
// @description					Type "Bearer" followed by a space and JWT token.

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	config.LoadEnv()
	config.ConnectToDB()

	app := fiber.New()
	app.Get("/swagger/*", swaggo.HandlerDefault)
    app.Get("/docs/*", swaggo.New(swaggo.Config{
        URL:               "http://example.com/doc.json",
        DeepLinking:       false,
        DocExpansion:      "none",
        OAuth2RedirectUrl: "http://localhost:3000/swagger/oauth2-redirect.html",
    }))

	userRepository := repositories.NewUserRepository()
	userService := services.NewUserService(userRepository)
	userController := controllers.NewUserController(userService)

	komporRepository := repositories.NewKomporRepository()
	komporService := services.NewKomporService(userRepository, komporRepository)
	komporController := controllers.NewKomporController(komporService)

	espRepository := repositories.NewEspRepository()
	espService := services.NewEspService(espRepository, userRepository)
	espController := controllers.NewEspController(espService)


	routes.Setup(app, userController, komporController, espController)

	port := config.APPConfig.APPPort

	log.Print("App running on port: ", port)

	log.Fatal(app.Listen(":" + port))
}