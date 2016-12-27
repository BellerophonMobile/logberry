package logberry

import (
	"sort"
	"reflect"
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
type EventDataInt64 int64
type EventDataUInt64 uint64
type EventDataFloat64 float64

func (x EventDataMap) WriteRecurse(out io.Writer) {

	fmt.Fprintf(out, "{")

	keys := make([]string, len(x))
	i := 0
	for k,_ := range(x) {
		keys[i] = k
		i++
	}

	sort.Strings(keys)
	
	for _,k := range(keys) {
		v := x[k]

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

func (x EventDataInt64) WriteRecurse(out io.Writer) {
	fmt.Fprintf(out, "%v", x)
}

func (x EventDataUInt64) WriteRecurse(out io.Writer) {
	fmt.Fprintf(out, "%v", x)
}

func (x EventDataFloat64) WriteRecurse(out io.Writer) {
	fmt.Fprintf(out, "%v", x)
}

func MakeEventData(data []interface{}) EventData {

	switch len(data) {
	case 0:
		return EventDataMap(nil)

	case 1:
		return Copy(data[0])
		
	default:

		data := make(EventDataSlice, len(data))
		
		for i, v := range(data) {
			data[i] = Copy(v)
		}

		return data
	}
	
	return EventDataMap(nil)

}

func Copy(data interface{}) EventData {

	val, null := rolldown(data)
	if null {
		return EventDataMap(nil)
	}

	switch val.Kind() {

	case reflect.Struct:
		return copystruct(val)

	case reflect.Map:
		return copymap(val)

	case reflect.Array: fallthrough
	case reflect.Slice:
		return copyslice(val)

	default:
		return copydata(val)

	}
	
	return EventDataString("###")
	
}

func rolldown(data interface{}) (reflect.Value, bool) {

	if data == nil {
		return reflect.Value{}, true
	}

	val := reflect.ValueOf(data)

	// Chain through any pointers or interfaces
	done := false
	for !done {
		switch val.Kind() {
		case reflect.Interface:
			fallthrough
		case reflect.Ptr:

			if val.IsNil() {
				return reflect.Value{}, true
			}

			val = val.Elem()

		default:
			done = true
		}
	}

	return val, false

}


func copystruct(val reflect.Value) EventDataMap {

	res := make(EventDataMap)

	var vtype = val.Type()
	var haspublic bool

	for i := 0; i < val.NumField(); i++ {
		var f = val.Field(i)
		if f.IsValid() && f.CanInterface() && !strings.Contains(vtype.Field(i).Tag.Get("logberry"), "quiet") {
			res[vtype.Field(i).Name] = Copy(f.Interface())
			haspublic = true
		}
	}

	// Special case: If the value is an error but has no accessible
	// fields, call its Error() function to get a text representation.
	if !haspublic && val.CanAddr() {
		v2 := val.Addr().Interface()
		if err, ok := (v2).(error); ok {
			res["Error"] = EventDataString(err.Error())
		}
	}
	
	return res

}

func copymap(val reflect.Value) EventDataMap {

	res := make(EventDataMap)
	
	var vals = val.MapKeys()
	for _, k := range vals {
		v := val.MapIndex(k)
		if k.CanInterface() && v.CanInterface() {
			res[fmt.Sprint(k.Interface())] = Copy(v.Interface())
		}
	}

	return res

}

func copyslice(val reflect.Value) EventDataSlice {

	arr := make(EventDataSlice, val.Len())

	for i := 0; i < val.Len(); i++ {
		arr[i] = Copy(val.Index(i))
	}

	return arr

}

func copydata(val reflect.Value) EventData {

	switch val.Kind() {

	case reflect.Int: fallthrough
	case reflect.Int8: fallthrough
	case reflect.Int16: fallthrough
	case reflect.Int32: fallthrough
	case reflect.Int64:
		return EventDataInt64(val.Int())		

	case reflect.Uint: fallthrough
	case reflect.Uint8: fallthrough
	case reflect.Uint16: fallthrough
	case reflect.Uint32: fallthrough
	case reflect.Uint64:
		return EventDataUInt64(val.Uint())

	case reflect.Float32: fallthrough
	case reflect.Float64:
		return EventDataFloat64(val.Uint())

	default:
		return EventDataString(val.String())

	}

	return EventDataString("##")

}
