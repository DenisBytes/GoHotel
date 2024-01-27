package api

import "github.com/gofiber/fiber/v2"

func ErrorHandler (c *fiber.Ctx, err error) error{
	if apiError, ok := err.(Error); ok {
		return c.Status(apiError.Code).JSON(apiError)
	}
	apiError := NewError(fiber.StatusInternalServerError, err.Error())
	return c.Status(apiError.Code).JSON(apiError.Err)
}

type Error struct {
	Code int    `json:"code"`
	Err  string `json:"error"`
}

// this bellow imlplements the error interface. so in our hndlers we can return this instead of normal errors
func (e Error) Error() string {
	return e.Err
}

func NewError(code int, err string) Error {
	return Error{
		Code: code,
		Err:  err,
	}
}

func ErrInvalidID() Error {
	return Error{
		Code: fiber.StatusBadRequest,
		Err:  "invalid id given",
	}
}

func ErrUnauthorized() Error {
	return Error{
		Code: fiber.StatusUnauthorized,
		Err:  "unathorized request",
	}
}

func ErrBadRequest() Error {
	return Error{
		Code: fiber.StatusBadRequest,
		Err:  "invalid json request",
	}
}

func ErrResourceNotFound(res string) Error {
	return Error{
		Code: fiber.StatusNotFound,
		Err:  res + " resource not found",
	}
}