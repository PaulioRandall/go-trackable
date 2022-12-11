package trackerr

import (
	"errors"
	"strings"
)

// ErrorFormatter formats an error for stack trace printing.
//
// Each error string will be printed on a line of its own so implementations
// should not prefix or suffix a linefeed unless they want gappy print outs.
type ErrorFormatter func(errMsg string, e error, isFirst bool) string

// ErrorStack calls ErrorStackf with simple default formatting.
//
// Checkpoints are prefixed and suffixed with `——` while ordinary errors are
// prefixed with `⤷ `:
//
//		——Workflow error——
//		⤷ Failed to read data
//		⤷ Error handling CSV file
//		——File could not be opened "splay/example/data/acid-rain.csv"——
//		⤷ open splay/example/data/acid-rain.csv
//		⤷ no such file or directory
func ErrorStack(e error) string {
	return ErrorStackf(e, func(msg string, e error, isFirst bool) string {
		sb := strings.Builder{}

		if IsCheckpoint(e) {
			sb.WriteString("——")
			sb.WriteString(msg)
			sb.WriteString("——")
			return sb.String()
		}

		if isFirst {
			sb.WriteString("  ")
		} else {
			sb.WriteString("⤷ ")
		}

		sb.WriteString(msg)
		return sb.String()
	})
}

// ErrorStackf returns a human readable stack trace for the error.
//
// The format function f may be nil for no formatting.
func ErrorStackf(e error, f ErrorFormatter) string {
	sb := strings.Builder{}

	for i, cause := range AsStack(e) {
		errMsg := ErrorWithoutCause(cause)

		if f != nil {
			errMsg = f(errMsg, cause, i == 0)
		}

		sb.WriteString(errMsg)
		sb.WriteRune('\n')
	}

	return sb.String()
}

// AsStack recursively unwraps the error returning a slice of errors.
//
// The passed error e will be first and root cause last.
func AsStack(e error) []error {
	var stack []error

	for e != nil {
		stack = append(stack, e)
		e = errors.Unwrap(e)
	}

	return stack
}

// ErrorWithoutCause removes the cause from error messages that use the format
// '%s: %w'. Where s is the error message and w is the cause's message.
func ErrorWithoutCause(e error) string {
	s := e.Error()

	if _, ok := e.(*untrackedError); ok {
		return s
	}

	if _, ok := e.(*trackedError); ok {
		return s
	}

	cause := errors.Unwrap(e)

	if cause == nil {
		return s
	}

	s = strings.TrimSuffix(s, cause.Error())
	s = strings.TrimSpace(s)
	return strings.TrimSuffix(s, ":")
}
