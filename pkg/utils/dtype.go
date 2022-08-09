package utils

import (
	"go/types"
	"reflect"
	"strconv"
)

type EmptyValue types.Nil

type V struct {
	value interface{}
}

func NewV(value interface{}) V {
	return V{
		value: value,
	}
}

func (v *V) IsSet() bool {
	ref := reflect.ValueOf(v.value)

	return ref.IsValid()
}

func (v *V) IsEmpty() bool {
	if !v.IsSet() {
		return true
	}

	ref := reflect.ValueOf(v.value)

	if ref.Kind() == reflect.Array || ref.Kind() == reflect.Slice {
		return ref.Len() == 0
	}

	return ref.Type() == reflect.ValueOf(EmptyValue{}).Type()
}

func (v *V) IsString() bool {
	if !v.IsSet() {
		return false
	}

	ref := reflect.ValueOf(v.value)

	return ref.Kind() == reflect.String
}

func (v *V) IsInt() bool {
	if !v.IsSet() {
		return false
	}

	ref := reflect.ValueOf(v.value)

	switch ref.Kind() {
	case reflect.Int:
		return true
	case reflect.String:
		if _, err := strconv.ParseInt(v.value.(string), 10, 64); err == nil {
			return true
		}

		return false
	default:
		return false
	}
}

func (v *V) IsBool() bool {
	if !v.IsSet() {
		return false
	}

	ref := reflect.ValueOf(v.value)

	switch ref.Kind() {
	case reflect.Bool:
		return true
	default:
		return false
	}
}

func (v *V) String() string {
	if !v.IsSet() {
		return ""
	}

	ref := reflect.ValueOf(v.value)

	switch ref.Kind() {
	case reflect.Int:
		return strconv.Itoa(int(ref.Int()))
	case reflect.String:
		return v.value.(string)
	default:
		return ""
	}
}

func (v *V) Int() int {
	if !v.IsSet() {
		return 0
	}

	ref := reflect.ValueOf(v.value)

	switch ref.Kind() {
	case reflect.Int:
		return int(ref.Int())
	case reflect.String:
		if i, err := strconv.ParseInt(v.value.(string), 10, 64); err == nil {
			return int(i)
		}

		return 0
	default:
		return 0
	}
}

func (v *V) Int64() int64 {
	return int64(v.Int())
}

func (v *V) Bool() bool {
	if !v.IsSet() {
		return false
	}

	ref := reflect.ValueOf(v.value)

	switch ref.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if ref.Int() == 0 {
			return false
		}

		return true
	case reflect.String:
		if ref.String() == "" || ref.String() == "false" || ref.String() == "0" {
			return false
		}

		return true
	default:
		return false
	}
}
