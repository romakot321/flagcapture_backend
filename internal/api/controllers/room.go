package controllers

import (
  "github.com/gofiber/fiber/v2"
  "github.com/romakot321/game-backend/internal/api/services"
  "github.com/romakot321/game-backend/internal/api/models"
)

type RoomController interface {
  Register(router fiber.Router)
}

type roomController struct {
  roomService services.RoomService
}

func (c *roomController) Register(router fiber.Router) {
  router.Get("/", c.list)
  router.Post("/", c.create)
}

func (c *roomController) list(ctx *fiber.Ctx) error {
  rooms := c.roomService.List()
  return ctx.JSON(rooms)
}

func (c *roomController) create(ctx *fiber.Ctx) error {
  var schema models.RoomCreateSchema;

  if err := ctx.BodyParser(&schema); err != nil {
    return ctx.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
      "error": "Unprocessable entity",
      "details": err.Error(),
    })
  }

  room := c.roomService.Create(schema.Name);

  return ctx.JSON(room);
}

func NewRoomController(roomService services.RoomService) RoomController {
  return &roomController{roomService: roomService}
}
