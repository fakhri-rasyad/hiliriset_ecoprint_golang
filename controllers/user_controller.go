package controllers

import (
	"hiliriset_ecoprint_golang/models"
	"hiliriset_ecoprint_golang/services"
	"hiliriset_ecoprint_golang/utils"

	"github.com/gofiber/fiber/v3"
)

type UserController interface{
	RegisterUser(ctx fiber.Ctx) error
	LoginUser(ctx fiber.Ctx) error
}

type UserControllerImpl struct {
	userService services.UserService
}

func NewUserController(s services.UserService) UserController {
	return &UserControllerImpl{userService: s}
}

// ShowAccount godoc
// @Summary      Register akun user
// @Description  Endpoint untuk mendaftarkan akun user ke sistem
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param		 user body  models.UserDataResponse true "Raw User Json Data"
// @Success      201  {object}  models.UserDataResponse
// @Failure      400  {object}  utils.Response
// @Router       /v1/auth/register [post]
func (c *UserControllerImpl) RegisterUser(ctx fiber.Ctx) error {
	user := new(models.User)

	if err := ctx.Bind().Body(&user); err != nil {
		return utils.BadRequest(ctx, "Error saat parsing data", err)
	}

	if err := c.userService.CreateUser(user) ; err != nil {
		return utils.BadRequest(ctx, "Gagal meregistrasikan user", err)
	}
	
	UserResponseData := models.UserDataResponse{
		Username: user.Username,
		Email: user.Email,
		Role: user.Role,
	}
	return utils.CreationSuccess(ctx, "Akun user telah dibuat", UserResponseData)
}

//ShowAccount godoc
// @Summary 		Login user
// @Description 	Endpoint untuk authentikasi dan verifikasi user
// @Tags			User
// @Accept			json
// @Produce			json
// @Param 			user body models.UserDataResponse true "Raw user json data"
// @Success 		200 {object} models.UserLoginResponse
// @Failure 		400 {object} utils.Response
// @Router 			/v1/auth/login [post]
func (c *UserControllerImpl) LoginUser(ctx fiber.Ctx) error{
	user := new(models.User)

	if err := ctx.Bind().Body(&user); err != nil {
		return utils.BadRequest(ctx, "Error saat parsing data", err)
	}

	if err := c.userService.LoginUser(user.Email, user.Password); err != nil {
		return utils.BadRequest(ctx, "Error saat login", err)
	}

	newToken, err := utils.GenerateToken(user.Username, user.Email)

	if err != nil {
		return utils.BadRequest(ctx, "Error saat login", err)
	}

	LoginResponse := models.UserLoginResponse{
		Username: user.Username,
		BearerToken: newToken,
	}

	return utils.SuccessResponse(ctx, "Login success", LoginResponse)
}