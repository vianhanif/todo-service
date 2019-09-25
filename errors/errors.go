package errors

// BaseError is the basic error with minimum construct
//
// should be the base of all errors emitted from within application, use it as promoted fields
type BaseError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// NewBaseError returns new instance of BaseError
func NewBaseError(code, name string) BaseError {
	return BaseError{
		Code: code, Message: name,
	}
}

func (e BaseError) Error() string {
	return e.Message
}

// CommonError return basic/common error
type CommonError struct {
	BaseError
}

// NewCommonError returns new CommonError
func NewCommonError(msg string) CommonError {
	return CommonError{
		BaseError: NewBaseError("CommonError", msg),
	}
}

// ValidationError is error used for all validation related errors
type ValidationError struct {
	BaseError
	Fields []ValidationErrorField `json:"fields,omitempty"`
}

// ValidationErrorField is a type to contain information on field error
type ValidationErrorField struct {
	Name    string `json:"name"`    // name of the field
	Message string `json:"message"` // error message related to the field
}

// NewValidationError returns new ValidationError
func NewValidationError(msg string) ValidationError {
	return ValidationError{
		BaseError: NewBaseError("ValidationError", msg),
		Fields:    []ValidationErrorField{},
	}
}

// ClearFieldErrors clears all field errors
func (e *ValidationError) ClearFieldErrors() {
	e.Fields = e.Fields[:0]
}

// FieldError sets field error
func (e *ValidationError) FieldError(name, message string) {
	for i, f := range e.Fields {
		if f.Name == name {
			e.Fields[i].Message = message
			return
		}
	}
	e.Fields = append(e.Fields, ValidationErrorField{
		Name:    name,
		Message: message,
	})
}

// GetFieldError returns field error
func (e *ValidationError) GetFieldError(name string) *ValidationErrorField {
	for _, f := range e.Fields {
		if f.Name == name {
			return &f
		}
	}
	return nil
}

// AuthError indicates error on authorization
type AuthError struct {
	BaseError
}

// NewAuthError returns new AuthError
func NewAuthError(msg string) AuthError {
	return AuthError{
		BaseError: NewBaseError("UnauthorizedError", msg),
	}
}

// ServiceError indicates error on service level
type ServiceError struct {
	BaseError
}

// NewServiceError returns new ServiceError
func NewServiceError(msg string) ServiceError {
	return ServiceError{
		BaseError: NewBaseError("ServiceError", msg),
	}
}

// NotFoundError indicates error when data not found
type NotFoundError struct {
	BaseError
}

// NewNotFoundError returns new NotFoundError
func NewNotFoundError(msg string) NotFoundError {
	return NotFoundError{
		BaseError: NewBaseError("NotFoundError", msg),
	}
}
