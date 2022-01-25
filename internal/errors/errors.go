package errors

import "fmt"

type ArgumentError struct {
	message string
}

func NewArgumentError(message string) ArgumentError {
	return ArgumentError{message}
}

func (e ArgumentError) Error() string {
	return e.message
}

func APIRequestError(err error) error {
	return fmt.Errorf("unexpected error during sending request: %s", err.Error())
}

func APIBadStatusCode(code int) error {
	return fmt.Errorf("API respond with bad status code (%d)", code)
}

func APIUnexpectedSchema(err error) error {
	return fmt.Errorf("API respond with unexpected JSON schema: %s", err.Error())
}
