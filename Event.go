package logberry

import (
	"time"
)

// These strings are common classes of events, used in the Task event
// generation functions to identify types of occurrences.
const (
	BEGIN         string = "begin"
	END           string = "end"
	CONFIGURATION string = "configuration"
	READY         string = "ready"
	STOPPED       string = "stopped"
	INFO          string = "info"
	SUCCESS       string = "success"
	WARNING       string = "warning"
	ERROR         string = "error"
)

// Event captures an annotated occurrence or message, a log entry.
type Event struct {
	TaskID   uint64
	ParentID uint64

	Component string

	Event   string
	Message string
	Data    EventDataMap

	Timestamp time.Time
}
