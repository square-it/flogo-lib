package activity

// Error is an activity error
type Error struct {
	errorStr  string
	errorCode string
	errorData interface{}
}

// NewError creates a error object
func NewError(errorText string) *Error {
	return &Error{errorStr: errorText}
}

// NewErrorWithData creates a error object with associated data
func NewErrorWithData(errorText string, errorData interface{}) *Error {
	return &Error{errorStr: errorText, errorData: errorData}
}

func NewErrorWithCode(errorText string, code string) *Error {
	return &Error{errorStr: errorText, errorCode: code}
}

func NewErrorWithDataAndCode(errorText string, errorData interface{}, code string) *Error {
	return &Error{errorStr: errorText, errorData: errorData, errorCode: code}
}

// Error implements error.Error()
func (e *Error) Error() string {
	return e.errorStr
}

// Data returns any associated error data
func (e *Error) Data() interface{} {
	return e.errorData
}

// Code returns any associated error code
func (e *Error) Code() string {
	return e.errorCode
}
