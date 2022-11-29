package trackable

// ErrorStack accepts an error and produces a printable stack trace.
//
// It combines an error's message with it's cause's message at the point its
// Error function is called.
//
// A custom algorithm can be assigned during program initialisation. By default
// DescMiniTrace is used.
var ErrorStack = DescMiniTrace

// DescMiniTrace (Descending Minimalist Trace) is a minimalist ErrorStack
// algorithm.
//
// It simply appends the cause to the message with a line break, tab, and
// a Unicode arrow to represent the ascent. This means the root cause will be
// the last error printed which is usually what you want. Ain't no one got time
// for scrolling in terminals.
func DescMiniTrace(msg string, cause error) string {
	if cause == nil {
		return msg
	}
	return msg + "\n⤷ " + cause.Error()
}

// AsceMiniTrace (Ascending Minimalist Trace) is the reverse of DescMiniTrace.
//
// The message is appended to the cause instead of the other way around.
func AsceMiniTrace(msg string, cause error) string {
	if cause == nil {
		return "↱ " + msg
	}
	return cause.Error() + "\n↱ " + msg
}
