package logberry

// ErrorListeners are registered to Roots and notified of internal
// logging errors.  Examples include an inability to write to disk, or
// contact a logging server.  That notification could be utilized to
// prompt the administrator in some way or take other action.  It is
// an error with unspecified behavior to add an ErrorListener to more
// than one Root simultaneously.
type ErrorListener interface {
	Error(err error)
}
