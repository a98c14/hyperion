package errors

import "fmt"

// Wraps the given error with function name
func Wrap(fn string, err error) error {
	return fmt.Errorf("%s\n-> %w", fn, err)
}

// Wraps the given error with function name and custom message
func WrapMsg(fn string, msg string, err error) error {
	return fmt.Errorf("%s | %s \n-> %w", fn, msg, err)
}
