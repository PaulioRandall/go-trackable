package trackerr

// TrackedError represents a trackable node in an error stack trace.
type TrackedError struct {
	UntrackedError
	id int
}

// Is returns true if the passed error is equivalent to the receiving
// error.
//
// This is a shallow comparison so causes are not checked. It is designed
// to work with the Is function exposed by the standard errors package.
func (e TrackedError) Is(other error) bool {
	if e2, ok := other.(*TrackedError); ok {
		return e.id == e2.id
	}
	return false
}
