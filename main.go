package main

import (
	"hiliriset_ecoprint_golang/config"
	"hiliriset_ecoprint_golang/controllers"
	mqttpackage "hiliriset_ecoprint_golang/mqtt_package"
	"hiliriset_ecoprint_golang/repositories"
	"hiliriset_ecoprint_golang/routes"
	"hiliriset_ecoprint_golang/services"
	websocketutils "hiliriset_ecoprint_golang/websocket_utils"

	"log"

	"github.com/gofiber/fiber/v3"

	_ "hiliriset_ecoprint_golang/docs"

	swaggo "github.com/gofiber/contrib/v3/swaggo"
)

// @title           Hiliriset Ecoprint API
// @version         1.0
// @description     API untuk sistem monitoring perebusan kain ecoprint
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:3000
// @BasePath  /

// @securityDefinitions.apiKey  ApiKeyAuth
// @in                          header
// @name                        Authorization
// @description                 Type "Bearer" followed by a space and JWT token.

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
    // 1. Config and DB
    config.LoadEnv()
    config.ConnectToDB()

    // 2. Fiber app
    app := fiber.New()
    app.Get("/swagger/*", swaggo.HandlerDefault)
    app.Get("/docs/*", swaggo.New(swaggo.Config{
        URL:               "http://example.com/doc.json",
        DeepLinking:       false,
        DocExpansion:      "none",
        OAuth2RedirectUrl: "http://localhost:3000/swagger/oauth2-redirect.html",
    }))

    // 3. Repositories
    userRepository  := repositories.NewUserRepository()
    komporRepository := repositories.NewKomporRepository()
    espRepository   := repositories.NewEspRepository()
    boSeRepository  := repositories.NewBoSeRepository()
    seReRepository  := repositories.NewSeReRepository()

    // 4. WebSocket hub — needed by MQTT handler
    wsHub := websocketutils.NewHub()
    go wsHub.Run()

    // 5. Services that don't depend on MQTT
    seReService := services.NewSessionRecordService(seReRepository, boSeRepository)

    // 6. MQTT — handler first, then client, then inject publisher back
    mqttHandler := mqttpackage.NewMQTTHandler(boSeRepository, seReService, wsHub, espRepository, komporRepository)
    mqttClient  := mqttpackage.NewMQTTClient(mqttHandler)
    mqttHandler.SetPublisher(mqttClient)
    defer mqttClient.Disconnect()

    // 7. Services that depend on MQTT
    boSeService  := services.NewBoSeService(boSeRepository, userRepository, komporRepository, espRepository, mqttClient,mqttHandler)
    userService  := services.NewUserService(userRepository)
    komporService := services.NewKomporService(userRepository, komporRepository)
    espService   := services.NewEspService(espRepository, userRepository)

    // 8. Controllers
    userController  := controllers.NewUserController(userService)
    komporController := controllers.NewKomporController(komporService)
    espController   := controllers.NewEspController(espService)
    boSeController  := controllers.NewBoSeController(boSeService, seReService)
    seReController  := controllers.NewSessionRecordController(seReService)

    wsController := websocketutils.NewWSController(wsHub)
    // 9. Routes
    routes.Setup(app, userController, komporController, espController, boSeController, seReController, wsController)

    port := config.APPConfig.APPPort
    log.Print("App running on port: ", port)
    log.Fatal(app.Listen(":" + port))
}
