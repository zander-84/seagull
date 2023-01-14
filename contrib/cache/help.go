package cache

import "reflect"

func getValue(in any) reflect.Value {
	if reflect.ValueOf(in).Type().Kind() != reflect.Ptr {
		return reflect.ValueOf(in)
	}
	return _getValue(reflect.ValueOf(in).Elem())

}

func _getValue(in reflect.Value) reflect.Value {
	if in.Type().Kind() != reflect.Ptr {
		return in
	}
	return _getValue(in.Elem())
}
