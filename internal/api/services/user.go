package services

import (
  "log"

  "github.com/romakot321/game-backend/internal/api/models"
  "github.com/romakot321/game-backend/internal/api/repositories"
)

type UserService interface {
  Authenticate(msg models.MessageAuthenticateData) *models.UserModel
  GetList() []*models.UserModel
}

type userService struct {
  userRepository repositories.UserRepository
}

func (s userService) Authenticate(msg models.MessageAuthenticateData) *models.UserModel {
  user := s.userRepository.Get(msg.Username)
  if user == nil {
    user = models.MakeUserModel(msg.Username)
    s.userRepository.Add(user)
    log.Print("Create user with name ", msg.Username)
  }
  return user
}

func (s userService) GetList() []*models.UserModel {
  return s.userRepository.GetList()
}

func NewUserService(userRepository repositories.UserRepository) UserService {
  return &userService{userRepository: userRepository}
}
