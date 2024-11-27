package dto

// swagger:model error
type ResponseError struct {
	Code    int    `json:"code"`
	Err     string `json:"error"`
	Message string `json:"message"`
}

func NewRespError(code int, err, message string) *ResponseError {
	return &ResponseError{
		Code:    code,
		Err:     err,
		Message: message,
	}
}
