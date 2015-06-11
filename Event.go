package logberry

import (
	"time"
)

// String tags identifying common classes of events.
const (	
	BEGIN         string = "begin"
	END           string = "end"
	CONFIGURATION string = "configuration"
	READY         string = "ready"
	STOPPED       string = "stopped"
	INFO          string = "info"
	WARNING       string = "warning"
	ERROR         string = "error"
)


// Event captures an annotated occurrence or message, a log entry.
type Event struct {
	TaskID uint64
	ParentID uint64

	Component string
	
	Event string
	Message string
	Data D

	Timestamp time.Time
}
