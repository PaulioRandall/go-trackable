package trackable

// id is an ID counter, always increment before assigning a new ID.
var id int

// newId is the function used to generate trackable error IDs.
//
// Only IDs greater than zero are considered trackable.
var newId = func() int {
	id++
	return id
}
