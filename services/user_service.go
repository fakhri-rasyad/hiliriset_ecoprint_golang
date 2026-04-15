package services

import (
	"errors"
	"hiliriset_ecoprint_golang/models"
	"hiliriset_ecoprint_golang/repositories"
	"hiliriset_ecoprint_golang/utils"

	"github.com/google/uuid"
)

type UserService interface {
	CreateUser(*models.User) error
}

type UserServiceImpl struct{
	userRepo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &UserServiceImpl{userRepo: repo}
}

func (s *UserServiceImpl) CreateUser(user *models.User) error{
	existingUser, _:= s.userRepo.FindByEmail(user.Email)
	if existingUser.InternalID != 0 {
		return errors.New("Email telah digunakan, silahkan memakai email baru")
	}
	hashedPassword, _ := utils.HashPassword(user.Password)

	user.Password = hashedPassword
	user.Role = "user"
	user.PublicID = uuid.New()

	return s.userRepo.Create(user)
}