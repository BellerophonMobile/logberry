package logberry

import (
	"fmt"
	"io"
	"strings"
)

type EventData interface {
	WriteRecurse(io.Writer)
}

type EventDataMap map[string]EventData
type EventDataSlice []EventData
type EventDataString string
type EventDataInt32 int32
type EventDataFloat32 float32

func (x EventDataMap) WriteRecurse(out io.Writer) {

	fmt.Fprintf(out, "{")
	
	for k,v := range(x) {
		if strings.ContainsAny(k, "\"= {}") {
			fmt.Fprintf(out, " %q=", k)
		} else {
			fmt.Fprintf(out, " %v=", k)
		}

		v.WriteRecurse(out)

	}

	fmt.Fprintf(out, " }")

}

func (x EventDataSlice) WriteRecurse(out io.Writer) {

	fmt.Fprintf(out, "[ ")
	
	for k,v := range(x) {

		if k > 0 {
			fmt.Fprintf(out, ", ")
		}
		
		v.WriteRecurse(out)

	}

	fmt.Fprintf(out, " ]")

}

func (x EventDataString) WriteRecurse(out io.Writer) {
	fmt.Fprintf(out, "%q", x)
}

func (x EventDataInt32) WriteRecurse(out io.Writer) {
	fmt.Fprintf(out, "%v", x)
}

func (x EventDataFloat32) WriteRecurse(out io.Writer) {
	fmt.Fprintf(out, "%v", x)
}
