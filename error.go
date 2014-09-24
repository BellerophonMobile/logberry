package logberry

type logerror struct {
	msg string
	src []error
}


func NewError(msg string) *logerror {
	return &logerror{msg, []error{}}
}

func WrapError(msg string, err error) *logerror {
	return &logerror{msg, []error{err}}
}


func (e *logerror) AddError(err error) {
	e.src = append(e.src, err)
}


func (e *logerror) Error() string {

	if len(e.src) == 1 {
		return e.msg + ": " + e.src[0].Error()
	}

	if len(e.src) != 0 {
		s := e.msg + "---multiple errors: " + e.src[0].Error()
		for i := 1; i < len(e.src); i++ {
			s += "; " + e.src[i].Error()
		}
		return s;
	}

	return e.msg;
}
