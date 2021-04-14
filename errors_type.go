package berrors

type ErrorType string

const (
	Email    ErrorType = "email"
	Required ErrorType = "required"
	Login    ErrorType = "login"
)

var ErrorTypes = map[ErrorType]string{
	Email:    "The email given is not correct",
	Required: "The field is required",
	Login:    "The login and password does not match",
}
