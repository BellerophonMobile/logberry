package logberry

//----------------------------------------------------------------------
//----------------------------------------------------------------------

type StatementClass int

const (
	ERROR StatementClass = iota
	METADATA
	INFO
)

type Data map[string]string


//------------------------------------------------------
var outputdrivers = []OutputDriver{}


//----------------------------------------------------------------------
//----------------------------------------------------------------------
func AddOutput(driver OutputDriver) error {

	outputdrivers = append(outputdrivers, driver)
	return nil

}


//----------------------------------------------------------------------
//----------------------------------------------------------------------
func logprimitive(component string,
	                class StatementClass,
	                msg string,
                  data interface{}) error {

	var accumerror *logerror = nil

	for _,driver := range(outputdrivers) {
		err := driver.log(component, class, msg, data)
		if err != nil {
			if accumerror == nil {
				accumerror = NewError("Error outputting to log(s)")
			}
			accumerror.AddError(err)
		}
	}

	return accumerror
}


//----------------------------------------------------------------------
//----------------------------------------------------------------------
type ComponentLog struct {
	Component string
}

func NewComponentLog(component string) *ComponentLog {
	return &ComponentLog{
		Component: component,
	}
}

func (log *ComponentLog) Build(build BuildMetadata) error {
	return logprimitive(log.Component, METADATA, "Build", build)
}

func (log *ComponentLog) Info(msg string, data interface{}) error {
	return logprimitive(log.Component, INFO, msg, data)
}

func (log *ComponentLog) Error(msg string, err error) error {
	return logprimitive(log.Component, ERROR, msg,
		&Data{ "Error": err.Error() })
}
