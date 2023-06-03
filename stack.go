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
//		Workflow error
//		⤷ Failed to read data
//		⤷ Error handling CSV file
//		⤷ open splay/example/data/acid-rain.csv
//		⤷ no such file or directory
func ErrorStack(e error) string {
	return ErrorStackf(e, func(errMsg string, e error, isFirst bool) string {
		sb := strings.Builder{}

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

	for i, cause := range SliceStack(e) {
		errMsg := cause.Error()

		if f != nil {
			errMsg = f(errMsg, cause, i == 0)
		}

		sb.WriteString(errMsg)
		sb.WriteRune('\n')
	}

	return sb.String()
}

// SliceStack recursively unwraps the error returning a slice of errors. The
// passed error e will be first and root cause last.
//
//		charlie := trackerr.Untracked("Charlie's message")
//		bob := trackerr.Wrap(charlie, "Bob's message")
//		alice := trackerr.Wrap(bob, "Alice's message")
//
//		result := SliceStack(alice)
//
//		// result: [
//		// 	alice,
//		// 	bob,
//		// 	charlie,
//		// ]
func SliceStack(e error) []error {
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

// Stack accepts a an array of ErrorWrappers and converts it into a stack trace
// by recursively calling CasuedBy.
//
// The first item is the root cause and the last item the head.
//
//		head := trackerr.New("head message")
//		mid := trackerr.New("mid level message")
//		root := trackerr.New("root cause message")
//
//		e := Stack(root, mid, head)
//
//		// head message
//		// ⤷ mid level message
//		// ⤷ root cause message
func Stack(e error, errs ...ErrorThatWraps) error {
	if e == nil {
		return nil
	}

	for _, err := range errs {
		e = err.CausedBy(e).(ErrorThatWraps)
	}

	return e
}

// Squash calls trackerr.ErrorStack with the error e then uses the
// result as the message for a new error; which is returned.
func Squash(e error) error {
	s := ErrorStack(e)
	return Untracked(s)
}

// Squashf is the same as squash but allows an ErrorFormatter to be used to
// format the error stack string.
func Squashf(e error, f ErrorFormatter) error {
	s := ErrorStackf(e, f)
	return Untracked(s)
}
