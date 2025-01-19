package api

import (
  "log"

	"github.com/gofiber/fiber/v2"
  "github.com/gofiber/fiber/v2/middleware/cors"
  "github.com/romakot321/game-backend/internal/api/repositories"
  "github.com/romakot321/game-backend/internal/api/services"
  "github.com/romakot321/game-backend/internal/api/controllers"
)

func Run() {
  userRepository := repositories.NewUserRepository()
  userService := services.NewUserService(userRepository)
  roomService := services.NewRoomService()
  connectionController := controllers.NewConnectionController(userService, roomService)
  roomController := controllers.NewRoomController(roomService)

  app := fiber.New()
  router := fiber.New()

  app.Use(cors.New(cors.Config{
    AllowOrigins: "*",
    AllowHeaders: "*",
    AllowMethods: "*",
    AllowCredentials: false,
  }))

  app.Mount("/api", router)
  router.Route("/game", connectionController.Register)
  router.Route("/room", roomController.Register)

  log.Fatal(app.Listen("0.0.0.0:9000"))
}
