package models

type Room struct {
  Id int `json:"id"`
  Users []*UserModel `json:"users"`
  Name string `json:"name"`
}

func (r *Room) AddUser(model *UserModel) {
  for _, user := range r.Users {
    if user.Name == model.Name {
      return
    }
  }
  r.Users = append(r.Users, model)
}

type RoomCreateSchema struct {
  Name string `json:"name"`
}
