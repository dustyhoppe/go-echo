package errors

type MethodNotAllowedError struct{}

func (e *MethodNotAllowedError) Error() string {
	return "method not allowed"
}

func NewMethodNowAllowedError() *MethodNotAllowedError {
	return &MethodNotAllowedError{}
}
