package logberry

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"sort"
	"strings"
)

// A DBuilder is a type that can return logberry data when logged.
type DBuilder interface {
	D() D
}

type EventData interface {
	WriteTo(io.Writer)
}

type EventDataMap map[string]EventData
type EventDataSlice []EventData
type EventDataString string
type EventDataInt64 int64
type EventDataUInt64 uint64
type EventDataFloat64 float64
type EventDataBool bool

func (x EventDataMap) String() string {
	buff := new(bytes.Buffer)
	x.WriteTo(buff)
	return buff.String()
}

func (x EventDataMap) WriteTo(out io.Writer) {

	fmt.Fprintf(out, "{")

	keys := make([]string, len(x))
	i := 0
	for k := range x {
		keys[i] = k
		i++
	}

	sort.Strings(keys)

	for _, k := range keys {
		v := x[k]

		if strings.ContainsAny(k, "\"= {}[]") {
			fmt.Fprintf(out, " %q=", k)
		} else {
			fmt.Fprintf(out, " %v=", k)
		}

		v.WriteTo(out)

	}

	fmt.Fprintf(out, " }")

}

func (x EventDataSlice) WriteTo(out io.Writer) {

	fmt.Fprintf(out, "[")

	for k, v := range x {

		if k > 0 {
			fmt.Fprintf(out, ", ")
		}

		v.WriteTo(out)

	}

	fmt.Fprintf(out, "]")

}

func (x EventDataString) WriteTo(out io.Writer) {
	fmt.Fprintf(out, "%q", x)
}

func (x EventDataInt64) WriteTo(out io.Writer) {
	fmt.Fprintf(out, "%v", x)
}

func (x EventDataUInt64) WriteTo(out io.Writer) {
	fmt.Fprintf(out, "%v", x)
}

func (x EventDataFloat64) WriteTo(out io.Writer) {
	fmt.Fprintf(out, "%v", x)
}

func (x EventDataBool) WriteTo(out io.Writer) {
	fmt.Fprintf(out, "%v", x)
}

/*
func MakeEventData(data []interface{}) EventData {

	switch len(data) {
	case 0:
		return EventDataMap(nil)

	case 1:
		return makeeventdata(data[0])

	default:

		data := make(EventDataSlice, len(data))

		for i, v := range(data) {
			data[i] = copy(v)
		}

		return data
	}

	return EventDataMap(nil)

}
*/

func Copy(data interface{}) EventData {
	e, _ := copy(data)
	return e
}

func copy(data interface{}) (EventData, bool) {

	if der, ok := data.(DBuilder); ok {
		data = der.D()
	}

	val, null := rolldown(data)
	if null {
		return EventDataMap(nil), true
	}

	zero := true

	switch val.Kind() {

	case reflect.Struct:
		r := EventDataMap{}.aggregatestruct(val)
		if len(r) != 0 {
			zero = false
		}
		return r, zero

	case reflect.Map:
		r := EventDataMap{}.aggregatemap(val)
		if len(r) != 0 {
			zero = false
		}
		return r, zero

	default:
		return copydata(val)

	}

}

func Aggregate(data []interface{}) EventDataMap {

	x := EventDataMap{}
	for _, v := range data {
		x.Aggregate(v)
	}
	return x

}

func (x EventDataMap) Aggregate(data interface{}) EventDataMap {

	val, null := rolldown(data)
	if null {
		return x
	}

	switch val.Kind() {

	case reflect.Struct:
		x.aggregatestruct(val)

	case reflect.Map:
		x.aggregatemap(val)

	default:
		newval, zero := copydata(val)

		if zero {
			break
		}

		prev, find := x["value"]

		if find {
			switch p := prev.(type) {
			case EventDataSlice:
				x["value"] = append(p, newval)

			default:
				x["value"] = EventDataSlice{p, newval}
			}
		} else {
			x["value"] = newval
		}

	}

	return x

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

func (x EventDataMap) aggregatestruct(val reflect.Value) EventDataMap {

	var vtype = val.Type()

	fieldCount := 0

	for i := 0; i < val.NumField(); i++ {
		var f = val.Field(i)

		tags := vtype.Field(i).Tag.Get("logberry")
		isQuiet := strings.Contains(tags, "quiet")
		isAlways := strings.Contains(tags, "always")
		isHidden := strings.Contains(tags, "hidden")

		if f.IsValid() && f.CanInterface() && !isQuiet {
			fi := f.Interface()
			c, zero := copy(fi)
			if !zero || isAlways {
				if isHidden {
					x[vtype.Field(i).Name] = EventDataString("<!hidden!>")
				} else {
					x[vtype.Field(i).Name] = c
				}
				fieldCount++
			}
		}
	}

	// Special case: If the value is an error but has no accessible
	// fields, call its Error() function to get a text representation.
	if fieldCount == 0 {
		v2 := val
		if v2.CanAddr() {
			v2 = v2.Addr()
		}
		if err, ok := (v2.Interface()).(error); ok {
			x["Error()"] = EventDataString(err.Error())
		}
	}

	return x

}

func (x EventDataMap) aggregatemap(val reflect.Value) EventDataMap {

	var vals = val.MapKeys()
	for _, k := range vals {
		v := val.MapIndex(k)
		if k.CanInterface() && v.CanInterface() {
			x[fmt.Sprint(k.Interface())], _ = copy(v.Interface())
		}
	}

	return x

}

func copydata(val reflect.Value) (EventData, bool) {

	zero := true

	switch val.Kind() {
	case reflect.Array:
		fallthrough
	case reflect.Slice:
		arr := make(EventDataSlice, val.Len())
		for i := 0; i < val.Len(); i++ {
			arr[i], _ = copy(val.Index(i).Interface())
		}

		if len(arr) > 0 {
			zero = false
		}

		return arr, zero

	case reflect.Int:
		fallthrough
	case reflect.Int8:
		fallthrough
	case reflect.Int16:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int64:
		i := val.Int()
		if i != 0 {
			zero = false
		}
		return EventDataInt64(i), zero

	case reflect.Uint:
		fallthrough
	case reflect.Uint8:
		fallthrough
	case reflect.Uint16:
		fallthrough
	case reflect.Uint32:
		fallthrough
	case reflect.Uint64:
		u := val.Uint()
		if u != 0 {
			zero = false
		}
		return EventDataUInt64(u), zero

	case reflect.Float32:
		fallthrough
	case reflect.Float64:
		f := val.Float()
		if f != 0.0 {
			zero = false
		}
		return EventDataFloat64(f), zero

	case reflect.Bool:
		f := val.Bool()
		return EventDataBool(f), false

	default:
		// Special case: If the value is an error, call its Error() function
		// to get a text representation.
		if val.CanInterface() {
			v2 := val.Interface()
			if err, ok := (v2).(error); ok {
				return EventDataString(err.Error()), false
			}
		}

		s := val.String()
		if s != "" {
			zero = false
		}
		return EventDataString(s), zero

	}

}
