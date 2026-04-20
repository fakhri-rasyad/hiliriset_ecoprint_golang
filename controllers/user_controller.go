package controllers

import (
	"hiliriset_ecoprint_golang/models"
	"hiliriset_ecoprint_golang/services"
	"hiliriset_ecoprint_golang/utils"

	"github.com/gofiber/fiber/v3"
)

type UserController interface {
    RegisterUser(ctx fiber.Ctx) error
    LoginUser(ctx fiber.Ctx) error
}

type UserControllerImpl struct {
    userService services.UserService
}

func NewUserController(s services.UserService) UserController {
    return &UserControllerImpl{userService: s}
}

// RegisterUser godoc
// @Summary      Register akun user
// @Description  Endpoint untuk mendaftarkan akun user ke sistem
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        user body models.UserGorm true "Raw User Json Data"
// @Success      201  {object}  models.UserDataResponse
// @Failure      400  {object}  utils.Response
// @Router       /v1/auth/register [post]
func (c *UserControllerImpl) RegisterUser(ctx fiber.Ctx) error {
    var req models.UserGorm

    if err := ctx.Bind().Body(&req); err != nil {
        return utils.BadRequest(ctx, "Error saat parsing data", err)
    }

    if err := c.userService.CreateUser(&req); err != nil {
        return utils.BadRequest(ctx, "Gagal meregistrasikan user", err)
    }

    return utils.CreationSuccess(ctx, "Akun user telah dibuat", models.UserDataResponse{
        Username: req.Username,
        Email:    req.Email,
        Role:     req.Role,
    })
}

// LoginUser godoc
// @Summary      Login user
// @Description  Endpoint untuk authentikasi dan verifikasi user
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        user body models.UserGorm true "Raw user json data"
// @Success      200  {object}  models.UserLoginResponse
// @Failure      400  {object}  utils.Response
// @Router       /v1/auth/login [post]
func (c *UserControllerImpl) LoginUser(ctx fiber.Ctx) error {
    var req models.UserGorm

    if err := ctx.Bind().Body(&req); err != nil {
        return utils.BadRequest(ctx, "Error saat parsing data", err)
    }

    userBase, err := c.userService.LoginUser(req.Email, req.Password)
    if err != nil {
        return utils.BadRequest(ctx, "Error saat login", err)
    }

    newToken, err := utils.GenerateToken(userBase.Username, userBase.Email)
    if err != nil {
        return utils.InternalError(ctx, "Error saat membuat token", err)
    }

    return utils.SuccessResponse(ctx, "Login sukses", models.UserLoginResponse{
        Username:    userBase.Username,
        BearerToken: newToken,
    })
}