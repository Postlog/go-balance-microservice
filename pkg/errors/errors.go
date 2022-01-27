package errors

import "fmt"

// ArgumentError represents an error related to incorrectly passed arguments
type ArgumentError struct {
	message string
}

func NewArgumentError(message string) ArgumentError {
	return ArgumentError{message}
}

func (e ArgumentError) Error() string {
	return e.message
}

// APIRequestError wrapper that returns generic error
func APIRequestError(err error) error {
	return fmt.Errorf("unexpected error during sending request: %s", err)
}

// APIBadStatusCode wrapper that returns generic error
func APIBadStatusCode(code int) error {
	return fmt.Errorf("API respond with bad status code (%d)", code)
}

// APIUnexpectedSchema wrapper that returns generic error
func APIUnexpectedSchema(err error) error {
	return fmt.Errorf("API respond with unexpected JSON schema: %s", err)
}
