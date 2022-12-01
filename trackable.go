package trackable

import (
	"fmt"
)

// TODO: Interface should be the name of the interface (string, not a bool)

// TODO: Create a Session entity that allows sessions to be created where the
// TODO: ID pool is part of the session. Keep a package instance like 'http'
// TODO: package does because sessions will be redendant for most programs.

// TODO: Think about how to integrate file names and line numbers.
// TODO: - How, where, and when to collect them?
// TODO: - How to optimise print outs with them?
// TODO: - May have to redesign the Debug function?

// trackable represents an error that uses an integer as an identifier.
//
// If an ID is less than or equal to zero then the error is categorised as
// untracked according to errors.Is.
type trackable struct {
	id    int
	msg   string
	cause error
	iface bool
}

// Track returns a new trackable error, that is, one with a tracking ID.
//
// This function is designed to be called during package initialisation only.
// This means it should only be used to initialise package global variables,
// within init functions, or as part of a test.
func Track(msg string, args ...any) *trackable {
	return &trackable{
		id:  newId(),
		msg: fmt.Sprintf(msg, args...),
	}
}

// Interface is the same as Track except the trackable error is flagged as
// being at the boundary of a key interface.
func Interface(msg string, args ...any) *trackable {
	return &trackable{
		id:    newId(),
		msg:   fmt.Sprintf(msg, args...),
		iface: true,
	}
}

// Untracked returns a new error without a tracking ID.
func Untracked(msg string, args ...any) *trackable {
	return &trackable{
		msg: fmt.Sprintf(msg, args...),
	}
}

// Wrap returns a new error, without a tracking ID, that wraps a cause.
//
// It's an alternative to fmt.Errorf where the cause does not have to form part
// of the error message.
func Wrap(cause error, msg string, args ...any) *trackable {
	return &trackable{
		msg:   fmt.Sprintf(msg, args...),
		cause: cause,
	}
}

// WrapAtInterface is the same as Wrap but flags the but the trackable error is
// flagged as being at the boundary of a key interface.
func WrapAtInterface(cause error, msg string, args ...any) *trackable {
	return &trackable{
		msg:   fmt.Sprintf(msg, args...),
		cause: cause,
		iface: true,
	}
}

func (e trackable) Error() string {
	if e.cause == nil {
		return e.msg
	}

	return e.msg + ": " + e.cause.Error()
}

func (e trackable) String() string {
	return e.msg
}

func (e trackable) Unwrap() error {
	return e.cause
}

func (e trackable) Is(target error) bool {
	if e.id <= 0 {
		return false
	}

	if it, ok := target.(*trackable); ok {
		return e.id == it.id
	}

	return false
}

func (e trackable) Wrap(cause error) error {
	e.cause = cause
	return &e
}

func (e trackable) Because(msg string, args ...any) error {
	e.cause = Untracked(msg, args...)
	return &e
}

func (e trackable) BecauseOf(cause error, msg string, args ...any) error {
	e.cause = Wrap(cause, msg, args...)
	return &e
}

func (e trackable) Interface(msg string, args ...any) error {
	t := Untracked(msg, args...)
	t.iface = true
	e.cause = t
	return &e
}

func (e trackable) InterfaceOf(cause error, msg string, args ...any) error {
	t := Wrap(cause, msg, args...)
	t.iface = true
	e.cause = t
	return &e
}

func (e trackable) IsInterface() bool {
	return e.iface
}
