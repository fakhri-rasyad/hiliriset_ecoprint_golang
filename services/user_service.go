package services

import (
	"hiliriset_ecoprint_golang/models"
	"hiliriset_ecoprint_golang/repositories"
)

type UserService interface {
	CreateUser(*models.User) error
	FindByEmail(email string) (*models.User, error)
}

type UserServiceImpl struct{
	userRepo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &UserServiceImpl{userRepo: repo}
}

func (s *UserServiceImpl) CreateUser(user *models.User) error{
	return s.userRepo.Create(user)
}

func (s *UserServiceImpl) FindByEmail(email string) (*models.User, error){
	return s.userRepo.FindByEmail(email)
}