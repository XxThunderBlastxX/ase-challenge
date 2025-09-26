package response

// Legacy response structures (kept for backward compatibility)
type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Data any `json:"data"`
}

// Enhanced response structures with better error handling
type BaseResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

type ErrorResponseWithCode struct {
	BaseResponse
	Code    string `json:"code"`
	Details any    `json:"details,omitempty"`
}

type SuccessResponseWithCode struct {
	BaseResponse
	Data any `json:"data,omitempty"`
}

// Validation error details structure
type ValidationErrorDetails struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Value   any    `json:"value,omitempty"`
}

type ValidationErrors struct {
	Errors []ValidationErrorDetails `json:"errors"`
}

// Pagination response structure
type PaginatedResponse struct {
	SuccessResponseWithCode
	Pagination PaginationMeta `json:"pagination,omitempty"`
}

type PaginationMeta struct {
	CurrentPage  int   `json:"current_page"`
	PerPage      int   `json:"per_page"`
	TotalPages   int   `json:"total_pages"`
	TotalRecords int64 `json:"total_records"`
	HasNext      bool  `json:"has_next"`
	HasPrev      bool  `json:"has_prev"`
}

// Helper functions to create standardized responses
func NewErrorResponse(code, message string, details ...any) ErrorResponseWithCode {
	response := ErrorResponseWithCode{
		BaseResponse: BaseResponse{
			Success: false,
			Message: message,
		},
		Code: code,
	}

	if len(details) > 0 {
		response.Details = details[0]
	}

	return response
}

func NewSuccessResponse(message string, data any) SuccessResponseWithCode {
	return SuccessResponseWithCode{
		BaseResponse: BaseResponse{
			Success: true,
			Message: message,
		},
		Data: data,
	}
}

func NewValidationErrorResponse(errors []ValidationErrorDetails) ErrorResponseWithCode {
	return ErrorResponseWithCode{
		BaseResponse: BaseResponse{
			Success: false,
			Message: "Validation failed",
		},
		Code: "VALIDATION_ERROR",
		Details: ValidationErrors{
			Errors: errors,
		},
	}
}

func NewPaginatedResponse(message string, data any, pagination PaginationMeta) PaginatedResponse {
	return PaginatedResponse{
		SuccessResponseWithCode: SuccessResponseWithCode{
			BaseResponse: BaseResponse{
				Success: true,
				Message: message,
			},
			Data: data,
		},
		Pagination: pagination,
	}
}
