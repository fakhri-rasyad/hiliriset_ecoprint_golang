package repositories

import (
	"hiliriset_ecoprint_golang/config"
	"hiliriset_ecoprint_golang/models"
)

type UserRepository interface {
    CreateUser(req *models.UserGorm) (*models.UserBase, error)
    FindByEmail(email string) (*models.UserBase, error)
    FindByPublicID(publicID string) (*models.UserBase, error)
}

type UserRepositoryImpl struct{}

func NewUserRepository() UserRepository {
    return &UserRepositoryImpl{}
}

func (r *UserRepositoryImpl) CreateUser(req *models.UserGorm) (*models.UserBase, error) {
    if err := config.DB.Create(req).Error; err != nil {
        return nil, err
    }

    base := req.ToBase()
    return &base, nil
}

func (r *UserRepositoryImpl) FindByEmail(email string) (*models.UserBase, error) {
    var gormModel models.UserGorm

    if err := config.DB.Where("email = ?", email).First(&gormModel).Error; err != nil {
        return nil, err
    }

    base := gormModel.ToBase()
    return &base, nil
}

func (r *UserRepositoryImpl) FindByPublicID(publicID string) (*models.UserBase, error) {
    var gormModel models.UserGorm

    if err := config.DB.Where("public_id = ?", publicID).First(&gormModel).Error; err != nil {
        return nil, err
    }

    base := gormModel.ToBase()
    return &base, nil
}