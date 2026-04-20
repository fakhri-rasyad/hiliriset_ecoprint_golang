// ========================================================
// user_service.go
// ========================================================
package services

import (
	"errors"
	"hiliriset_ecoprint_golang/models"
	"hiliriset_ecoprint_golang/repositories"
	"hiliriset_ecoprint_golang/utils"

	"github.com/google/uuid"
)

type UserService interface {
    CreateUser(req *models.UserGorm) error
    LoginUser(email, password string) (*models.UserBase, error)
}

type UserServiceImpl struct {
    ur repositories.UserRepository
}

func NewUserService(ur repositories.UserRepository) UserService {
    return &UserServiceImpl{ur: ur}
}

func (s *UserServiceImpl) CreateUser(req *models.UserGorm) error {
    existing, _ := s.ur.FindByEmail(req.Email)
    if existing != nil && existing.InternalID != 0 {
        return errors.New("email telah digunakan, silahkan memakai email baru")
    }

    hashedPassword, err := utils.HashPassword(req.Password)
    if err != nil {
        return err
    }

    req.Password = hashedPassword
    req.Role = "user"

	req.PublicID = uuid.New()
    _, err = s.ur.CreateUser(req)
    return err
}

func (s *UserServiceImpl) LoginUser(email, password string) (*models.UserBase, error) {
    existingUser, err := s.ur.FindByEmail(email)
    if err != nil {
        return nil, errors.New("email belum terdaftar")
    }

    if !utils.CheckPasswordHash(password, existingUser.Password) {
        return nil, errors.New("password yang dimasukkan salah")
    }

    return existingUser, nil
}