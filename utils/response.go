package utils

import "github.com/gofiber/fiber/v3"

type Response struct {
	Status     string
	StatusCode int
	Message    string
	Data       interface{}
	Error      string
}

func BadRequest(ctx fiber.Ctx, message string, err error) error {
	return ctx.Status(fiber.StatusBadRequest).JSON(Response{
		Status: "400 Bad Request",
		StatusCode: fiber.StatusBadRequest,
		Message: message,
		Error: err.Error(),
	})
}

func Unauthorized(ctx fiber.Ctx, message string, err error) error {
	return ctx.Status(fiber.StatusUnauthorized).JSON(Response{
		Status: "401 Unathorized",
		StatusCode: fiber.StatusUnauthorized,
		Message: message,
		Error: err.Error(),
	})
}

func NotFound(ctx fiber.Ctx, message string, err error) error {
	return ctx.Status(fiber.StatusNotFound).JSON(Response{
		Status: "404 Not Found",
		StatusCode: fiber.StatusNotFound,
		Message: message,
		Error: err.Error(),
	})
}

func InternalError(ctx fiber.Ctx, message string, err error) error{
	return ctx.Status(fiber.StatusInternalServerError).JSON(Response{
		Status: "500 Internal Server Error",
		StatusCode: fiber.StatusInternalServerError,
		Message: message,
		Error: err.Error(),
	})
}

func CreationSuccess(ctx fiber.Ctx, message string, data interface{}) error {
	return ctx.Status(fiber.StatusCreated).JSON(Response{
		Status: "201 Status Created",
		StatusCode: fiber.StatusCreated,
		Message: message,
		Data: data,
	})
}

func SuccessResponse(ctx fiber.Ctx, message string, data interface{}) error {
	return ctx.Status(fiber.StatusOK).JSON(Response{
		Status: "200 Status OK",
		StatusCode: fiber.StatusOK,
		Message:  message,
		Data: data,
	})
}

func UnauthorizedReponse(ctx fiber.Ctx, message string, err error) error {
	return ctx.Status(fiber.StatusUnauthorized).JSON(Response{
		Status: "401 Unathorized",
		StatusCode: fiber.StatusUnauthorized,
		Message: message,
		Error: err.Error(),
	})
}