package logberry

import (
	"reflect"
	"fmt"
)

type D map[string]interface{}


//----------------------------------------------------------------------
//----------------------------------------------------------------------
func (x *D) Set(k string, v interface{}) *D {
	(*x)[k] = v
	return x
}


//----------------------------------------------------------------------
//----------------------------------------------------------------------
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

	switch val.Kind() {

	case reflect.Interface: fallthrough
	case reflect.Ptr:
		val = val.Elem()
		fallthrough

	case reflect.Struct:
		var vtype = val.Type()
		for i := 0; i < val.NumField(); i++ {
			var f = val.Field(i)
			d[vtype.Field(i).Name] = f.Interface()
		}

	case reflect.Map:
		var vals = val.MapKeys()
		for _, k := range(vals) {
			d[fmt.Sprint(k.Interface())] = val.MapIndex(k).Interface()
		}

	default:
		d["value"] = data

		// end switch type
	}

	return &d
	// end DBuild
}

//------------------------------------------------------
func (x *D) CopyFromD(data *D) *D {

	if data == nil {
		return x
	}

	for k,v := range((*data)) {
		(*x)[k] = v
	}

	return x

	// end CopyFromD
}

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

	case reflect.Interface: fallthrough
	case reflect.Ptr:
		val = val.Elem()
		fallthrough

	case reflect.Struct:
		var vtype = val.Type()
		for i := 0; i < val.NumField(); i++ {
			var f = val.Field(i)
			(*x)[vtype.Field(i).Name] = f.Interface()
		}

	default:
		var prev,find = (*x)["value"]
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

//------------------------------------------------------
func (x *D) AggregateFrom(dataarr []interface{}) *D {

	for i := 0; i < len(dataarr); i++ {
		x.CopyFrom(dataarr[i])
	}

	return x

	// and AggregateFrom
}

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

	// Alternatively, if wanted D objects to be safe, just do:
	//	return (&D{}).AggregateFrom(dataarr)

	// end DAggregate
}
