package response

type ErrorResponse struct {
	Message string `json: "message"`
}

func NewErrorResponse(message string) *ErrorResponse {
	return &ErrorResponse{
		Message: message,
	}
}