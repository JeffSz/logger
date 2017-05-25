package logger

import "reflect"

func Contains(collect interface{}, value interface{}) bool{
	baseType := reflect.ValueOf(collect)
	switch reflect.TypeOf(collect).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < baseType.Len(); i++ {
			if baseType.Index(i).Interface() == value{
				return true
			}
		}
		return false
	case reflect.Map:
		return baseType.MapIndex(reflect.ValueOf(value)).IsValid()
	default:
		return false
	}
}