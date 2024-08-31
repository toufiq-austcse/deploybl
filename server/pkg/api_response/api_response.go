package api_response

import "strings"

// swagger:parameters Response
type Response struct {
	Code    int         `json:"code"` // This is Name
	Message string      `json:"message"`
	Errors  []string    `json:"errors"`
	Data    interface{} `json:"data"`
}
type EmptyObj struct {
}

func BuildResponse(code int, message string, data interface{}) Response {
	var errors = make([]string, 0)
	return Response{
		Code:    code,
		Message: message,
		Errors:  errors,
		Data:    data,
	}
}

func BuildErrorResponse(code int, message string, err string, data interface{}) Response {
	splitError := strings.Split(err, "\n")
	return Response{
		Code:    code,
		Message: message,
		Errors:  splitError,
		Data:    data,
	}
}
