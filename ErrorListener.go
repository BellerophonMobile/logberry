package logberry

// ErrorListeners are notified of internal logging errors.  Examples
// include an inability to write an entry to disk, or contact a
// logging server.  That notification could be utilized to prompt the
// administrator in some way or take other action.
type ErrorListener interface {
	Error(err error)
}
