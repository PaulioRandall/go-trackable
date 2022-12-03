package trackable

import (
	"fmt"
)

// TODO: Create separate struct for interface errors

var global = NewSession()

type session struct {
	id int
}

// NewSession creates a new session, that is, a new separate pool of unique
// trackable error ids.
func NewSession() *session {
	return &session{}
}

// newId is the function used to generate trackable error IDs. Only IDs greater
// than zero are considered trackable.
func (s *session) newId() int {
	s.id++
	return s.id
}

// Track returns a new trackable error, that is, one with a tracking ID.
//
// This function is designed to be called during package initialisation only.
// This means it should only be used to initialise package global variables,
// within init functions, or as part of a test.
func (s *session) Track(msg string, args ...any) *trackable {
	return &trackable{
		id:  s.newId(),
		msg: fmt.Sprintf(msg, args...),
	}
}

// Interface is a trackable error given a name to indicate it being at the
// boundary of a key interface.
//
// Interface errors don't need a message. To add one use the Because receiving
// function on the returned trackable error.
func (s *session) Interface(name string) *trackable {
	return &trackable{
		id:    s.newId(),
		iface: name,
	}
}

// Track calls the Track function on the internal global singleton Session.
//
// This function is designed to be called during package initialisation only.
// This means it should only be used to initialise package global variables,
// within init functions, or as part of a test.
func Track(msg string, args ...any) *trackable {
	return global.Track(msg, args...)
}

// Interface calls the Interface function on the internal global singleton
// Session.
//
// Interface errors don't need a message. To add one use the Because receiving
// function on the returned trackable error.
//
// This function is designed to be called during package initialisation only.
// This means it should only be used to initialise package global variables,
// within init functions, or as part of a test.
func Interface(name string) *trackable {
	return global.Interface(name)
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

// WrapAtInterface returns a new error, without a trracking ID, that wraps a
// cause and provides the name of the interface.
//
// These errors don't need a message. To add one use the Because receiving
// function on the returned trackable error.
func WrapAtInterface(cause error, name string) *trackable {
	return &trackable{
		cause: cause,
		iface: name,
	}
}
