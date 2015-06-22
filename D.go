package logberry

import (
	"fmt"
	"reflect"
	"bytes"
	"io"
)


// These strings are common data keys, to be optionally used in
// literal D objects or calls to AddData on Tasks.
const (
	CALCULATION string = "Calculation"
	FILE string = "File"
	RESOURCE string = "Resource"
	SERVICE string = "Service"
	USER string = "User"
	URL string = "URL"
	BYTES string = "Bytes"
	ID string = "ID"
	ENDPOINT string = "Endpoint"
)


// D maps capture data to be logged as key/value pairs associated with
// Tasks or particular events.
type D map[string]interface{}

// Set assigns field k of the data object to value v, overriding any
// existing value.  Its return is the host D itself.
func (x D) Set(k string, v interface{}) D {
	x[k] = v
	return x
}

// DBuild populates a D object from another arbitrary object.  If that
// object is a D, it is returned.  Otherwise, DBuild is the same as
// instantiating a blank D and invoking CopyFrom on it and the data.
func DBuild(data interface{}) D {

	d, ok := data.(D)
	if ok {
		return d
	}

	d = D{}
	d.CopyFrom(data)
	return d

	// end DBuild
}

// CopyFromD takes each key/value from the given data and sets it
// within the host.  Values are not copied.  The modified host D is
// itself returned.
func (x D) CopyFromD(data D) D {

	for k, v := range(data) {
		x[k] = v
	}

	return x

	// end CopyFromD
}

// CopyFrom populates the host x from the given data object.  This
// applies the following rules:
//
//   data is nil                Nothing happens.
//
//   data is a D                Calls x.CopyFromD(data)
//
//   data is a struct           Each exposed field in data is set as a
//                              key/value in x, maintaining the value.
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

	switch val.Kind() {

	case reflect.Struct:
		var vtype = val.Type()
		var haspublic bool
		for i := 0; i < val.NumField(); i++ {
			var f = val.Field(i)
			if f.CanInterface() {
				x[vtype.Field(i).Name] = f.Interface()
				haspublic = true
			}
		}

		if !haspublic {
			if err,ok := data.(error); ok {
				x["Error"] = err.Error()
			}
		}
		
	case reflect.Map:
		var vals = val.MapKeys()
		for _, k := range vals {
			x[fmt.Sprint(k.Interface())] = val.MapIndex(k).Interface()
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

	// end CopyFrom
}


// AggregateFrom populates x from the given data array.  It does this
// by calling x.CopyFrom() on each element in the array.  It returns
// the modified host x.
func (x D) AggregateFrom(data []interface{}) D {

	for _,e := range(data) {
		x.CopyFrom(e)
	}

	return x

	// and AggregateFrom
}

// DAggregate returns a new D object populated from the given array.
// The following rules are applied.  DBuild is run on the first
// element, and then CopyFrom is run from the resultant D on each
// other element of the array.  Note that an effect of this is that if
// the first element of the array is a D, it is reused, and will be
// modified if there are other elements in the array.
//
// It is up to the caller to threadsafe and otherwise correctly share
// the object.  This may happen because the caller is using the same D
// object in another goroutine, or if a Root or OutputDriver buffers
// the event and logs it asynchronously.
func DAggregate(data []interface{}) D {

	// This is done this way rather than creating a blank and calling
	// AggregateFrom() on it in order to not create D objects
	// unnecessarily, i.e., if the user passes a single one in, it will
	// be used directly.

	if data == nil || len(data) == 0 {
		return D{}
	}

	var accum = DBuild(data[0])

	for i := 1; i < len(data); i++ {
		accum.CopyFrom(data[i])
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

	switch val.Kind() {

		/*
	case reflect.Interface: fallthrough
	case reflect.Ptr:

		if val.IsNil() {
			_,e := fmt.Fprint(buffer, "nil")
			return e
		}

		return textrecurse(buffer, wrap, val.Elem().Interface())
*/

	case reflect.Map:
		if wrap {
			_,e := fmt.Fprint(buffer, "{")
			if e != nil { return e}
		}
		
		var vals = val.MapKeys()
		for _, k := range vals {
			vval := val.MapIndex(k)
			if vval.IsValid() && vval.CanInterface() {
				_,e := fmt.Fprintf(buffer, " %s=", k.Interface())
				if e != nil { return e}
			
				e = textrecurse(buffer, true, vval.Interface())
				if e != nil { return e}
			}
		}
		
		if wrap {
			_,e := fmt.Fprint(buffer, " }")
			if e != nil { return e}
		}

		
	case reflect.Struct:
		if wrap {
			_,e := fmt.Fprint(buffer, "{")
			if e != nil { return e}
		}

		var vtype = val.Type()
		var haspublic bool
		for i := 0; i < val.NumField(); i++ {
			var f = val.Field(i)
			if f.IsValid() && f.CanInterface() && vtype.Field(i).Tag.Get("quiet") == "" {
				_,e := fmt.Fprintf(buffer, " %s=", vtype.Field(i).Name)
				if e != nil { return e}
				
				e = textrecurse(buffer, true, f.Interface())
				haspublic = true
				if e != nil { return e}
			}
		}

		if !haspublic {
				if err,ok := data.(error); ok {
					_,e := fmt.Fprintf(buffer, " Message=%q", err.Error())
				  if e != nil { return e}
			}
		}
		
		if wrap {
			_, e := fmt.Fprint(buffer, " }")
			if e != nil { return e}
		}


	case reflect.Array: fallthrough
	case reflect.Slice:
		_,e := fmt.Fprint(buffer, "[")
		if e != nil { return e}

		if val.Len() > 0 {
			e := textrecurse(buffer, true, val.Index(0).Interface())
			if e != nil { return e}
		}

		for i := 1; i < val.Len(); i++ {
			_,e := fmt.Fprint(buffer, ", ")
			if e != nil { return e}
			
			e = textrecurse(buffer, true, val.Index(i).Interface())
			if e != nil { return e}
		}
		
		_, e = fmt.Fprint(buffer, " ]")
		if e != nil { return e}


	case reflect.String:
		_,e := fmt.Fprintf(buffer, "%q", val.String())
		if e != nil { return e}

	default:
		if val.IsValid() && val.CanInterface() {
			_,e := fmt.Fprintf(buffer, "%v", val.Interface())
			if e != nil { return e}
		} else {
			_,e := fmt.Fprintf(buffer, "nil")
			if e != nil { return e}
		}
		
		// end switch type
	}

	return nil

}
