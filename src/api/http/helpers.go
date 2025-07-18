package http

import (
	response "github.com/bncunha/erp-api/src/api/responses"
)

func HandleError(err error) *response.ErrorResponse {
	
	return &response.ErrorResponse{
		Message: err.Error(),
	}
}