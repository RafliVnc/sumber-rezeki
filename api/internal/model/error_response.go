package model

type ErrorResponse struct {
	Code    int            `json:"code"`
	Message string         `json:"message"`
	Details []ErrorDetails `json:"details,omitempty"`
}

func NewErrorResponse(code int, err string, details []ErrorDetails) *ErrorResponse {
	return &ErrorResponse{
		Code:    code,
		Message: err,
		Details: details,
	}
}

func (e *ErrorResponse) Error() string {
	return e.Message
}

type ErrorDetails struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}
