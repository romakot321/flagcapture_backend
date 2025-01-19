package controllers

import (
  "strings"
  "math/rand"
  "log"
  "encoding/json"
  "errors"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
  "github.com/romakot321/game-backend/internal/api/models"
  "github.com/romakot321/game-backend/internal/api/services"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
    b := make([]byte, n)
    for i := range b {
        b[i] = letterBytes[rand.Intn(len(letterBytes))]
    }
    return string(b)
}

type ConnectionController interface {
  HandleConnection(wsconnection *websocket.Conn)
  Register(router fiber.Router)
}

type connectionController struct {
  connections map[string]*connection
  userService services.UserService
  roomService services.RoomService

  broadcast chan map[string]interface{}
}

func (s *connectionController) HandleConnection(wsconnection *websocket.Conn) {
  id := RandStringBytes(16)
  conn := &connection{ID: id, Wsconnection: wsconnection, Service: s}
  s.connections[id] = conn
  conn.Handle()
}

func (s *connectionController) Register(router fiber.Router) {
  router.Use(func(c *fiber.Ctx) error {
    if websocket.IsWebSocketUpgrade(c) {
      c.Locals("allowed", true)
      return c.Next()
    }
    return fiber.ErrUpgradeRequired
  })
  router.Get("/ws", websocket.New(func (c *websocket.Conn) {
    s.HandleConnection(c)
  }))
}

func (s *connectionController) runHub() {
  for {
    select {
    case message := <-s.broadcast:
      roomName := message["room"].(string)
      for _, conn := range s.connections {
        if conn.Room == nil || conn.Room.Name != roomName {
          continue
        }
        conn.Write(message)
      }
    }
  }
}

func NewConnectionController(userService services.UserService, roomService services.RoomService) ConnectionController {
  connections := make(map[string]*connection, 0)
  broadcast := make(chan map[string]interface{})
  s := &connectionController{
    connections: connections,
    userService: userService,
    broadcast: broadcast,
    roomService: roomService,
  }
  go s.runHub()
  return s
}

type connection struct {
  Service *connectionController
  Wsconnection *websocket.Conn
  ID string
  User *models.UserModel
  Room *models.Room
}

func (c *connection) read() (map[string]interface{}, error) {
  var data map[string]interface{}
  messageType, raw, err := c.Wsconnection.ReadMessage()
  if err != nil {
    return data, err
  }
  if messageType == websocket.TextMessage {
    json.Unmarshal(raw, &data)
    return data, nil
  }
  return data, errors.New("Unexpected message type while read")
}

func (c *connection) Write(data map[string]interface{}) {
  message, _ := json.Marshal(data)
  log.Print("Connection: ", c.ID, " Write: ", data)
  if err := c.Wsconnection.WriteMessage(websocket.TextMessage, message); err != nil {
    log.Fatal("Fail while write: ", err)
  }
}

func (c *connection) broadcast(data map[string]interface{}) {
  log.Print("Connection: ", c.ID, " Broadcast: ", data)
  data["room"] = c.Room.Name
  c.Service.broadcast <- data
}

func (c *connection) onOpen() {
  users := c.Service.userService.GetList()
  msg := make(map[string]interface{})
  data := make(map[string]interface{})
  data["users"] = make([]map[string]interface{}, 0)

  for _, user := range users {
    u := make(map[string]interface{})
    u["name"] = user.Name
    data["users"] = append(
      data["users"].([]map[string]interface{}),
      u,
    )
  }
  msg["event"] = models.MessageEventUsers
  msg["data"] = data
  c.Write(msg)
}

func (c *connection) Handle() {
  defer func() {
    delete(c.Service.connections, c.ID)
    c.Wsconnection.Close()
  }()

  c.onOpen()

  for {
    data, err := c.read()
    if err != nil {
      if strings.Contains(err.Error(), "close 1006") || strings.Contains(err.Error(), "close 1001") || strings.Contains(err.Error(), "connection reset by peer") {
        log.Print("Connection ", c.ID, " disconnected")
        break
      }
      log.Fatal("Unknown read error: ", err)
    }

    messageEvent, ok := data["event"].(float64)
    if !ok {
      continue
    }

    switch int(messageEvent) {
    case models.MessageEventAuthenticate:
      c.handleAuthenticate(models.MakeMessageAuthenticate(data))
    case models.MessageEventUserInput:
      c.handleUserInput(models.MakeMessageUserInput(data))
    case models.MessageEventUserPosition:
      c.handleUserPosition(models.MakeMessageUserPosition(data))
    case models.MessageEventNewEntity:
      c.handleNewEntity(models.MakeMessageNewEntity(data))
    }
  }
}

func (c *connection) handleAuthenticate(msg models.MessageAuthenticate) {
  user := c.Service.userService.Authenticate(msg.Data)
  room := c.Service.roomService.Authenticate(msg.Data)
  room.AddUser(user)
  c.User = user
  c.User.Status = models.UserStatusInRoom
  c.Room = room
  log.Print("Connection: " + c.ID + " Auth ", user, " ", room.Name)
  resp := make(map[string]interface{})
  resp["event"] = models.MessageEventAuthenticate
  resp["data"] = models.ToMap(user)
  c.broadcast(models.ToMap(resp))
}

func (c *connection) handleUserInput(msg models.MessageUserInput) {
  log.Print("Connection: " + c.ID + " User input. Key: ", msg.Data.Key, " IsRelease: ", msg.Data.IsRelease)
  msg.Data.Username = c.User.Name
  c.broadcast(models.ToMap(msg))
}

func (c *connection) handleUserPosition(msg models.MessageUserPosition) {
  log.Print("Connection: " + c.ID + " User position. x: ", msg.Data.X, " y: ", msg.Data.Y)
  msg.Data.Username = c.User.Name
  c.broadcast(models.ToMap(msg))
}

func (c *connection) handleNewEntity(msg models.MessageNewEntity) {
  log.Print("Connection: " + c.ID + " New Entity. Type: " + msg.Data.Type + " x: ", msg.Data.X, " y: ", msg.Data.Y)
  msg.Data.Username = c.User.Name
  c.broadcast(models.ToMap(msg))
}
