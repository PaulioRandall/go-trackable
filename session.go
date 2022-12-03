package trackable

import (
	"fmt"
)

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

// Interface is the same as Track except the trackable error is given an
// interface name as to indicate it being at the boundary of a key interface.
func (s *session) Interface(name string, msg string, args ...any) *trackable {
	return &trackable{
		id:    s.newId(),
		msg:   fmt.Sprintf(msg, args...),
		iface: name,
	}
}

// Track returns a new trackable error, that is, one with a tracking ID.
//
// This function is designed to be called during package initialisation only.
// This means it should only be used to initialise package global variables,
// within init functions, or as part of a test.
func Track(msg string, args ...any) *trackable {
	return global.Track(msg, args...)
}

// Interface is the same as Track except the trackable error is given an
// interface name as to indicate it being at the boundary of a key interface.
func Interface(name string, msg string, args ...any) *trackable {
	return global.Interface(name, msg, args...)
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
