package dto

// Standard certified API error codes
const (
	ErrCodeBadRequest         = "bad_request"
	ErrCodeNotFound           = "not_found"
	ErrCodeMethodNotAllowed   = "method_not_allowed"
	ErrCodeRateLimited        = "rate_limited"
	ErrCodeInternalError      = "internal_error"
	ErrCodeServiceUnavailable = "service_unavailable"
)

// ErrorDetail contains the standard code and message for an API error.
type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ErrorEnvelope wraps an API error response.
type ErrorEnvelope struct {
	Error ErrorDetail `json:"error"`
}

// NewErrorEnvelope constructs an ErrorEnvelope.
func NewErrorEnvelope(code, message string) ErrorEnvelope {
	return ErrorEnvelope{
		Error: ErrorDetail{
			Code:    code,
			Message: message,
		},
	}
}
