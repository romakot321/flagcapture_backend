package models

const (
  MessageEventAuthenticate = iota
  MessageEventConnected = iota
  MessageEventUserInput = iota
  MessageEventNewEntity = iota
  MessageEventUsers = iota
  MessageEventUserPosition = iota
)

type message struct {
  Event int `json:"-"`
  Data interface{} `json:"-"`
}

type MessageAuthenticateData struct {
  Username string `json:"username"`
  Room string `json:"room"`
}

type MessageAuthenticate struct {
  message
  Event int `json:"event"`
  Data MessageAuthenticateData `json:"data"`
}

func MakeMessageAuthenticate(data map[string]interface{}) MessageAuthenticate {
  return MessageAuthenticate{
    Event: MessageEventAuthenticate,
    Data: MessageAuthenticateData{
      Username: data["data"].(map[string]interface{})["username"].(string),
      Room: data["data"].(map[string]interface{})["room"].(string),
    },
  }
}

type MessageConnected struct {
  message
  Event int `json:"event"`
  Data interface{} `json:"data"`
}

type MessageUserInputData struct {
  Key string `json:"key"`
  IsRelease bool `json:"is_release"`
  Username string `json:"username"`
  Extra interface{} `json:"extra"`
}

type MessageUserInput struct {
  message
  Event int `json:"event"`
  Data MessageUserInputData `json:"data"`
}

func MakeMessageUserInput(data map[string]interface{}) MessageUserInput {
  msgData := data["data"].(map[string]interface{});
  username, ok := msgData["username"].(string);
  if (!ok) { username = ""; }

  return MessageUserInput{
    Event: MessageEventUserInput,
    Data: MessageUserInputData{
      Key: msgData["key"].(string),
      IsRelease: msgData["is_release"].(bool),
      Username: username,
      Extra: msgData["extra"],
    },
  }
}

type MessageUserPositionData struct {
  Username string `json:"username"`
  X int `json:"x"`
  Y int `json:"y"`
}

type MessageUserPosition struct {
  message
  Event int `json:"event"`
  Data MessageUserPositionData `json:"data"`
}

func MakeMessageUserPosition(data map[string]interface{}) MessageUserPosition {
  msgData := data["data"].(map[string]interface{});
  username, ok := msgData["username"].(string);
  if (!ok) { username = ""; }

  return MessageUserPosition{
    Event: MessageEventUserPosition,
    Data: MessageUserPositionData{
      X: int(msgData["x"].(float64)),
      Y: int(msgData["y"].(float64)),
      Username: username,
    },
  }
}

type MessageNewEntityData struct {
  Username string `json:"username"`
  X int `json:"x"`
  Y int `json:"y"`
  Type string `json:"type"`
}

type MessageNewEntity struct {
  message
  Event int `json:"event"`
  Data MessageNewEntityData `json:"data"`
}

func MakeMessageNewEntity(data map[string]interface{}) MessageNewEntity {
  msgData := data["data"].(map[string]interface{});
  username, ok := msgData["username"].(string);
  if (!ok) { username = ""; }

  return MessageNewEntity{
    Event: MessageEventNewEntity,
    Data: MessageNewEntityData{
      X: int(msgData["x"].(float64)),
      Y: int(msgData["y"].(float64)),
      Type: msgData["type"].(string),
      Username: username,
    },
  }
}
