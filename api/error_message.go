package api_errors

type ApiError struct {
	ErrorMesssage string `json:"error_message"`
	Code          int    `json:"code"`
}

func NewApiErrror(message string, code int) *ApiError {
	return &ApiError{
		ErrorMesssage: message,
		Code:          code,
	}
}

func (s *ApiError) Error() string {
	return s.ErrorMesssage
}
