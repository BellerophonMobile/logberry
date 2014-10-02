package logberry


type ComponentLog struct {
	Component string
	Data *D
}

//----------------------------------------------------------------------
//----------------------------------------------------------------------
func NewComponent(component string, data interface{}) *ComponentLog {
	var d = DBuild(data)
	LogPrimitive(component, INSTANTIATE, "Instantiate", d)
	return &ComponentLog{
		Component: component,
		Data: d,
	}
}

// This is intended to be used by libraries and distinct components,
// so that they can easily report their build information/repository
// status independently of the main program.
func (log *ComponentLog) Build(build BuildMetadata) {
	LogPrimitive(log.Component, CONFIGURATION, "Build", DBuild(build))
}

func (log *ComponentLog) Finalize() {
	LogPrimitive(log.Component, FINALIZE, "Finalize", nil)
}


//------------------------------------------------------
func (log *ComponentLog) Info(msg string, data ...interface{}) {
	LogPrimitive(log.Component, INFO, msg, DAggregate(data))
}

func (log *ComponentLog) Warning(msg string, data ...interface{}) {
	LogPrimitive(log.Component, WARNING, msg, DAggregate(data))
}

func (log *ComponentLog) Error(msg string, err error, data ...interface{}) error {
	// Note that this can't/shouldn't just throw err into the data blob
	// because the standard errors package error doesn't expose
	// anything, even the message.  So you basically have to reduce to a
	// string via Error().
	LogPrimitive(log.Component, ERROR, msg,
		DAggregate(data).Set("Error", err.Error()))
	return WrapError(err, msg)
}

// Failure is the same as Error but doesn't take an error object.
func (log *ComponentLog) Failure(msg string, data ...interface{}) error {
	LogPrimitive(log.Component, ERROR, msg, DAggregate(data))
	return NewError(msg)
}


//------------------------------------------------------
func (log *ComponentLog) Resource(msg string, resource interface{}) {
	LogPrimitive(log.Component, RESOURCE, msg,
		&D{"Resource": resource})
}

func (log *ComponentLog) Service(msg string, service interface{}) {
	LogPrimitive(log.Component, SERVICE, msg,
		&D{"Service": service})
}
