package api_response

import "strings"

// swagger:parameters Response
type Response struct {
	Code    int         `json:"code"` // This is Name
	Message string      `json:"message"`
	Errors  []string    `json:"errors"`
	Data    interface{} `json:"data"`
}
type PaginationResponse struct {
	Code       int         `json:"code"` // This is Name
	Message    string      `json:"message"`
	Errors     []string    `json:"errors"`
	Data       interface{} `json:"data"`
	Pagination *Pagination `json:"pagination"`
}
type Pagination struct {
	Total       int64 `json:"total"`
	CurrentPage int64 `json:"current_page"`
	LastPage    int64 `json:"last_page"`
	PerPage     int64 `json:"per_page"`
}
type EmptyObj struct{}

func BuildResponse(code int, message string, data interface{}) Response {
	errors := make([]string, 0)
	return Response{
		Code:    code,
		Message: message,
		Errors:  errors,
		Data:    data,
	}
}

func BuildPaginationResponse(
	code int,
	message string,
	data interface{},
	pagination *Pagination,
) PaginationResponse {
	errors := make([]string, 0)
	return PaginationResponse{
		Code:       code,
		Message:    message,
		Errors:     errors,
		Data:       data,
		Pagination: pagination,
	}
}

func BuildErrorResponse(code int, message string, err string, data interface{}) Response {
	splitError := []string{}
	if err != "" {
		splitError = strings.Split(err, "\n")
	}

	return Response{
		Code:    code,
		Message: message,
		Errors:  splitError,
		Data:    data,
	}
}
