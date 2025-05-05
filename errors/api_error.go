package errors

import "github.com/gin-gonic/gin"

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewErr(code int, message string) *ErrorResponse {
	return &ErrorResponse{
		Code:    code,
		Message: message,
	}
}

func APIError(c *gin.Context, err *ErrorResponse) {
	c.JSON(err.Code, err)
}

var (
	ErrBadRequestBody      = NewErr(400, "Bad Request body")
	ErrHeaderIsMissing     = NewErr(403, "Authorization header is missing")
	ErrInvalidHeaderFormat = NewErr(403, "Invalid authorization header format")
	ErrIncorrectToken      = NewErr(403, "Incorrect Token")
	ErrInternalServer      = NewErr(500, "An unexpected error occurred while processing the request")
)
