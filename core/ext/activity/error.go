package activity

// Error is an activity error
type Error struct {
	errorStr string
}

// NewError creates a error object
func NewError(errorText string) *Error {
	return &Error{errorStr: errorText}
}

// Error implements error.Error()
func (e *Error) Error() string {
	return e.errorStr
}
