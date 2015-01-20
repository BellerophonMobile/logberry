package logberry

import (
	"fmt"
	"reflect"
)

// D captures data to be logged as key/value pairs.
type D map[string]interface{}

// Set sets field k of the data object to v, overriding any existing
// value.  It returns the modified x.
func (x *D) Set(k string, v interface{}) *D {
	(*x)[k] = v
	return x
}

// DBuild populates a D object from another arbitrary object.  This
// construction applies the following rules:
//
//   data is nil                An empty D object is returned.
//
//   data is a *D               Data is returned, cast to *D.
//
//   data is a D                A pointer to data is returned.
//
//   data is a struct           Each exposed field in data is set as a
//                              key/value in the new D.  The value is
//                              the original value, not another D.
//
//   data is a map              Each key/value in data is set as a
//                              key/value in the new D.  fmt.Sprint is
//                              used to create a string key from the
//                              original map key.  The value stored is
//                              the original value, not another D.
//
//   otherwise                  The value data is placed within a field
//                              keyed as "value" of a new D.
//
// If data is a pointer or interface, the construction descends to the
// object it references.
func DBuild(data interface{}) *D {

	if data == nil {
		return &D{}
	}

	switch data.(type) {
	case *D:
		return data.(*D)

	case D:
		var res = data.(D)
		return &res
	}

	var d = D{}

	var val = reflect.ValueOf(data)

	// Chain through any pointers or interfaces
	done := false
	for !done {
		switch val.Kind() {
		case reflect.Interface:
			fallthrough
		case reflect.Ptr:
			val = val.Elem()

		default:
			done = true
		}
	}

	switch val.Kind() {

	case reflect.Struct:
		var vtype = val.Type()
		for i := 0; i < val.NumField(); i++ {
			var f = val.Field(i)
			if f.CanInterface() {
				d[vtype.Field(i).Name] = f.Interface()
			} else {
				d[vtype.Field(i).Name] = "unavailable"
			}
		}

	case reflect.Map:
		var vals = val.MapKeys()
		for _, k := range vals {
			d[fmt.Sprint(k.Interface())] = val.MapIndex(k).Interface()
		}

	default:
		d["value"] = data

		// end switch type
	}

	return &d

	// end DBuild
}

// CopyFromD takes each key/value from data and sets it within x.
// Values are not copied.  The parameter data may be nil, in which
// case nothing happens.  The modified x is returned.
func (x *D) CopyFromD(data *D) *D {

	if data == nil {
		return x
	}

	for k, v := range *data {
		(*x)[k] = v
	}

	return x

	// end CopyFromD
}

// CopyFromD populates x with data from the given data object.  This
// applies the following rules:
//
//   data is nil                Nothing happens.
//
//   data is a *D               Calls x.CopyFromD(data)
//
//   data is a D                Calls x.CopyFromD(&data)
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
// The modified x is returned.
func (x *D) CopyFrom(data interface{}) *D {

	if data == nil {
		return x
	}

	switch data.(type) {
	case *D:
		x.CopyFromD(data.(*D))
		return x

	case D:
		var res = data.(D)
		x.CopyFromD(&res)
		return x
	}

	var val = reflect.ValueOf(data)

	switch val.Kind() {

	case reflect.Interface:
		fallthrough
	case reflect.Ptr:
		val = val.Elem()
		fallthrough

	case reflect.Struct:
		var vtype = val.Type()
		for i := 0; i < val.NumField(); i++ {
			var f = val.Field(i)
			(*x)[vtype.Field(i).Name] = f.Interface()
		}

	case reflect.Map:
		var vals = val.MapKeys()
		for _, k := range vals {
			(*x)[fmt.Sprint(k.Interface())] = val.MapIndex(k).Interface()
		}

	default:
		var prev, find = (*x)["value"]
		if find {
			switch prev.(type) {
			case []interface{}:
				(*x)["value"] = append(prev.([]interface{}), data)

			default:
				(*x)["value"] = []interface{}{prev, data}
			}
		} else {
			(*x)["value"] = data
		}

		// end switch type
	}

	return x

	// end CopyFrom
}

// AggregateFrom populates x from the given data array.  It does this
// by calling x.CopyFrom() on each element in the array.  It returns
// the modified x.
func (x *D) AggregateFrom(dataarr []interface{}) *D {

	for i := 0; i < len(dataarr); i++ {
		x.CopyFrom(dataarr[i])
	}

	return x

	// and AggregateFrom
}

// DAggregate returns a new D object populated from the given array.
// The following rules are applied:
//
//   array is nil or empty      A new empty D is returned.
//
//   array has a single element DBuild(dataarr[0]) is returned.
//
//   otherwise                  DBuild is run on the first element, and
//                              then CopyFrom is run from that D on each
//                              other element of the array.
//
// Note that an effect of these rules is that if the first element of
// the array is a D, it is reused, and will be modified if there are
// other elements in the array.
//
// It is up to the caller to threadsafe and otherwise correctly share
// the object.  This may happen because the caller is using the same D
// object in another goroutine, or if an OutputDriver buffers the
// event and logs it asynchronously, as ThreadSafeOutput does.  The
// latter does not apply to later uses of the D object, directly or
// implicitly as in a Task object, through the same ThreadSafeOutput.
//
// All of that discussion is generally not an issue for most uses of
// the Logberry API.
func DAggregate(dataarr []interface{}) *D {

	// This is done this way rather than creating a blank and calling
	// AggregateFrom() on it in order to not create D objects
	// unnecessarily, i.e., if the user passes a single one in, it will
	// be used directly.

	if dataarr == nil || len(dataarr) == 0 {
		return &D{}
	}

	if len(dataarr) == 1 {
		return DBuild(dataarr[0])
	}

	var accum = DBuild(dataarr[0])

	for i := 1; i < len(dataarr); i++ {
		accum.CopyFrom(dataarr[i])
	}

	return accum

	// Alternatively, if wanted D objects to be threadsafe, just do:
	//	return (&D{}).AggregateFrom(dataarr)

	// end DAggregate
}
