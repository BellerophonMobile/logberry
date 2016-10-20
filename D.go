package logberry

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strings"
)

// D maps capture data to be logged as key/value pairs associated with
// Tasks or particular events.
type D map[string]interface{}

// CopyFrom populates a D from the given data object.  The following
// rules are applied:
//
//   data is nil                Nothing happens.
//
//   data is a struct           Each exposed field in data is set as a
//                              key/value in x, copying the value.
//
//   data is a map              Each key/value in data is set as a
//                              key/value in x.  fmt.Sprint is used to
//                              create a string key from the original
//                              map key.  The value stored is the
//                              original value, not a copy.
//
//   otherwise                  The value data is placed within a field
//                              keyed as "value" in x.  If the field
//                              already exists and is a single value,
//                              it is replaced with a list of both
//                              values.  Otherwise, if the field is
//                              already a list, data is appended to it.
//
// If data is a pointer or interface, the construction descends to the
// object it references.
//
// The modified host x is itself returned.
func (x D) CopyFrom(data interface{}) D {

	if data == nil {
		return x
	}

	var val = reflect.ValueOf(data)

	// Chain through any pointers or interfaces
	done := false
	for !done {
		switch val.Kind() {
		case reflect.Interface:
			fallthrough
		case reflect.Ptr:

			if val.IsNil() {
				return x
			}

			val = val.Elem()

		default:
			done = true
		}
	}

	// Apply the rules listed above
	switch val.Kind() {

	case reflect.Struct:

		var vtype = val.Type()
		var haspublic bool
		
		for i := 0; i < val.NumField(); i++ {
			var f = val.Field(i)
			if f.IsValid() && f.CanInterface() && !strings.Contains(vtype.Field(i).Tag.Get("logberry"), "quiet") {
				x[vtype.Field(i).Name] = f.Interface()
				haspublic = true
			}
		}

		// Special case: If the value is an error but has no accessible
		// fields, call its Error() function to get a text representation.
		if err, ok := data.(error); ok {
			if !haspublic {
				x["Error"] = err.Error()
			}
		}

	case reflect.Map:
		var vals = val.MapKeys()
		for _, k := range vals {
			v := val.MapIndex(k)
			if k.CanInterface() && v.CanInterface() {
				x[fmt.Sprint(k.Interface())] = v.Interface()
			}
		}

	default:
		var prev, find = x["value"]
		if find {
			switch prev.(type) {
			case []interface{}:
				x["value"] = append(prev.([]interface{}), data)

			default:
				x["value"] = []interface{}{prev, data}
			}
		} else {
			x["value"] = data
		}

		// end switch type
	}

	return x
	
}


// DAggregate returns a new D object populated from the given array
// using CopyFrom().  It is up to the caller to threadsafe and
// otherwise correctly share the D.  This may matter because the
// caller is using the same D object in another goroutine, or if a
// Root or OutputDriver buffers the event and logs it asynchronously.
func DAggregate(data []interface{}) D {

	if data == nil || len(data) == 0 {
		return D{}
	}
	
	var accum = D{}

	for _,d := range(data) {
		accum.CopyFrom(d)
	}

	return accum

	// end DAggregate
}

// String returns a text representation of the host D.  This is
// presented as a sequence of key=value pairs for arguably
// human-readable presentation.  To produce a JSON serialization,
// simply marshal on the D as usual.  It is equivalent to casting
// output from Text().
func (x D) String() string {
	return string(x.Text())
}

// Text returns a byte slice textual representation of the host D.
// This is presented as a sequence of key=value pairs for arguably
// human-readable presentation.  To produce a JSON serialization,
// simply marshal on the D as usual.
func (x D) Text() []byte {
	var buffer bytes.Buffer
	x.WriteTo(&buffer)
	return buffer.Bytes()
}

// WriteTo serializes the host D to the given io.Writer.  This is
// presented as a sequence of key=value pairs for arguably
// human-readable presentation.  To produce a JSON serialization,
// simply marshal on the D as usual.
func (x D) WriteTo(w io.Writer) error {
	return textrecurse(w, false, x)
}

func textrecurse(buffer io.Writer, wrap bool, data interface{}) error {

	var val = reflect.ValueOf(data)

	// Chain through any pointers or interfaces
	done := false
	for !done {
		switch val.Kind() {
		case reflect.Interface:
			fallthrough
		case reflect.Ptr:

			if val.IsNil() {
				return nil
			}

			val = val.Elem()

		default:
			done = true
		}
	}

	// Handle the possibilities
	switch val.Kind() {

	case reflect.Map:
		if wrap {
			_, e := fmt.Fprint(buffer, "{")
			if e != nil {
				return e
			}
		}

		var vals = val.MapKeys()
		for _, k := range vals {
			vval := val.MapIndex(k)
			if k.CanInterface() && vval.IsValid() && vval.CanInterface() {
				_, e := fmt.Fprintf(buffer, " %s=", k.Interface())
				if e != nil {
					return e
				}

				e = textrecurse(buffer, true, vval.Interface())
				if e != nil {
					return e
				}
			}
		}

		if wrap {
			_, e := fmt.Fprint(buffer, " }")
			if e != nil {
				return e
			}
		}

	case reflect.Struct:
		if wrap {
			_, e := fmt.Fprint(buffer, "{")
			if e != nil {
				return e
			}
		}

		var vtype = val.Type()
		var haspublic bool
		for i := 0; i < val.NumField(); i++ {
			var f = val.Field(i)
			if f.IsValid() && f.CanInterface() && !strings.Contains(vtype.Field(i).Tag.Get("logberry"), "quiet") {
				_, e := fmt.Fprintf(buffer, " %s=", vtype.Field(i).Name)
				if e != nil {
					return e
				}

				e = textrecurse(buffer, true, f.Interface())
				haspublic = true
				if e != nil {
					return e
				}
			}
		}

		if err, ok := data.(error); ok {
			if !haspublic {
				_, e := fmt.Fprintf(buffer, " Message=%q", err.Error())
				if e != nil {
					return e
				}
			} else if _, ok := err.(*Error); !ok {
				_, e := fmt.Fprintf(buffer, " Type=%T", err)
				if e != nil {
					return e
				}
			}
		}

		if wrap {
			_, e := fmt.Fprint(buffer, " }")
			if e != nil {
				return e
			}
		}

	case reflect.Array:
		fallthrough
	case reflect.Slice:
		_, e := fmt.Fprint(buffer, "[")
		if e != nil {
			return e
		}

		if val.Len() > 0 {
			v := val.Index(0)
			if v.CanInterface() {
				e := textrecurse(buffer, true, v.Interface())
				if e != nil {
					return e
				}
			} else {
				_, e := fmt.Fprint(buffer, "--")
				if e != nil {
					return e
				}
			}
		}

		for i := 1; i < val.Len(); i++ {
			_, e := fmt.Fprint(buffer, ", ")
			if e != nil {
				return e
			}

			v := val.Index(i)			
			if v.CanInterface() {
				e = textrecurse(buffer, true, val.Index(i).Interface())
				if e != nil {
					return e
				}
			} else {
				_, e := fmt.Fprint(buffer, "--")
				if e != nil {
					return e
				}
			}
		}

		_, e = fmt.Fprint(buffer, " ]")
		if e != nil {
			return e
		}

	case reflect.String:
		_, e := fmt.Fprintf(buffer, "%q", val.String())
		if e != nil {
			return e
		}

	default:
		if val.IsValid() && val.CanInterface() {
			_, e := fmt.Fprintf(buffer, "%v", val.Interface())
			if e != nil {
				return e
			}
		} else {
			_, e := fmt.Fprintf(buffer, "--")
			if e != nil {
				return e
			}
		}

		// end switch type
	}

	return nil

}
