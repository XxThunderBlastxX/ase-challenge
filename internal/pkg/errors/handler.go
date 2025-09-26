package errors

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/xxthunderblastxx/ase-challenge/internal/pkg/response"
)

// ErrorHandler creates a Fiber error handler middleware
func ErrorHandler() fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		// Log the error for debugging
		log.Printf("Error occurred: %v", err)

		// Check if it's our custom AppError
		if appErr, ok := err.(*AppError); ok {
			return c.Status(appErr.StatusCode).JSON(response.ErrorResponseWithCode{
				BaseResponse: response.BaseResponse{
					Success: false,
					Message: appErr.Message,
				},
				Code:    string(appErr.Code),
				Details: appErr.Details,
			})
		}

		// Check for Fiber errors
		if fiberErr, ok := err.(*fiber.Error); ok {
			return c.Status(fiberErr.Code).JSON(response.ErrorResponseWithCode{
				BaseResponse: response.BaseResponse{
					Success: false,
					Message: fiberErr.Message,
				},
				Code: "FIBER_ERROR",
			})
		}

		// Handle unknown errors
		return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorResponseWithCode{
			BaseResponse: response.BaseResponse{
				Success: false,
				Message: "An unexpected error occurred",
			},
			Code: string(UnknownError),
		})
	}
}

// HandleError is a helper function to return errors from handlers
func HandleError(c *fiber.Ctx, err error) error {
	if appErr, ok := err.(*AppError); ok {
		return c.Status(appErr.StatusCode).JSON(response.ErrorResponseWithCode{
			BaseResponse: response.BaseResponse{
				Success: false,
				Message: appErr.Message,
			},
			Code:    string(appErr.Code),
			Details: appErr.Details,
		})
	}

	// Convert standard errors to AppErrors
	appErr := FromError(err)
	return c.Status(appErr.StatusCode).JSON(response.ErrorResponseWithCode{
		BaseResponse: response.BaseResponse{
			Success: false,
			Message: appErr.Message,
		},
		Code:    string(appErr.Code),
		Details: appErr.Details,
	})
}

// HandleSuccess is a helper function to return success responses
func HandleSuccess(c *fiber.Ctx, data any, statusCode ...int) error {
	code := fiber.StatusOK
	if len(statusCode) > 0 {
		code = statusCode[0]
	}

	return c.Status(code).JSON(response.SuccessResponseWithCode{
		BaseResponse: response.BaseResponse{
			Success: true,
			Message: "Operation completed successfully",
		},
		Data: data,
	})
}

// HandleCreatedSuccess is a helper function for creation success responses
func HandleCreatedSuccess(c *fiber.Ctx, data any) error {
	return c.Status(fiber.StatusCreated).JSON(response.SuccessResponseWithCode{
		BaseResponse: response.BaseResponse{
			Success: true,
			Message: "Resource created successfully",
		},
		Data: data,
	})
}

// HandleNoContent is a helper function for no content responses
func HandleNoContent(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNoContent)
}
