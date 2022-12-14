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
// prefixed with `⤷ `. if the first error is not a checkpoint then the prefix
// is omitted:
//
//		Workflow error
//		⤷ Failed to read data
//		⤷ Error handling CSV file
//		——File could not be opened "play/example/data/acid-rain.csv"——
//		⤷ open splay/example/data/acid-rain.csv
//		⤷ no such file or directory
func ErrorStack(e error) string {
	return ErrorStackf(e, func(errMsg string, e error, isFirst bool) string {
		sb := strings.Builder{}

		if IsCheckpoint(e) {
			sb.WriteString("——")
			sb.WriteString(errMsg)
			sb.WriteString("——")
			return sb.String()
		}

		if !isFirst {
			sb.WriteString("⤷ ")
		}

		sb.WriteString(errMsg)
		return sb.String()
	})
}

// ErrorStackf returns a human readable stack trace for the error. The format
// function f may be nil for no formatting.
//
//		alice := trackerr.Untracked("Alice's message")
//		bob := trackerr.Checkpoint(alice, "Bob's message")
//		charlie := trackerr.Wrap(bob, "Charlie's message")
//		dan := trackerr.Wrap(charlie, "Dan's message")
//
//		s := trackerr.ErrorStackf(e, func(errMsg string, err error, isFirst bool) string {
//			if isFirst {
//				return "ERROR: " + errMsg
//			}
//			if trackerr.IsCheckpoint(err) {
//				return "*** " + errMsg + " ***"
//			}
//			return "Caused by: " + errMsg
//		}
//
//		// ERROR: Dan's message
//		// Caused by: Charlie's message
//		// *** Bob's message ***
//		// Caused by: Alice's message
func ErrorStackf(e error, f ErrorFormatter) string {
	sb := strings.Builder{}

	for i, cause := range AsStack(e) {
		errMsg := cause.Error()

		if f != nil {
			errMsg = f(errMsg, cause, i == 0)
		}

		sb.WriteString(errMsg)
		sb.WriteRune('\n')
	}

	return sb.String()
}

// AsStack recursively unwraps the error returning a slice of errors. The
// passed error e will be first and root cause last.
//
//		charlie := trackerr.Untracked("Charlie's message")
//		bob := trackerr.Wrap(charlie, "Bob's message")
//		alice := trackerr.Wrap(bob, "Alice's message")
//
//		result := AsStack(alice)
//
//		// result: [
//		// 	alice,
//		// 	bob,
//		// 	charlie,
//		// ]
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

	if _, ok := e.(*UntrackedError); ok {
		return s
	}

	if _, ok := e.(*TrackedError); ok {
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
