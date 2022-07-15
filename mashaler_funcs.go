package marshaler

import (
	"fmt"
	"reflect"
	"strconv"
)

const (
	csvStructTag        = string("csv")
	columnNamesPosition = int(0)
)

func isSlice(v any) bool {
	return reflect.TypeOf(v).Kind() == reflect.Slice
}

func isStruct(v any) bool {
	return reflect.TypeOf(v).Kind() == reflect.Struct
}

func isPointer(v any) bool {
	return reflect.TypeOf(v).Kind() == reflect.Pointer
}

func isSliceOfStruct(v any) bool {
	return isSlice(v) && (reflect.TypeOf(v).Elem().Kind() == reflect.Struct)
}

func isPointerToSliceOfStructs(v any) bool {
	return isPointer(v) && isSliceOfStruct(reflect.ValueOf(v).Elem().Interface())
}

//

func getValue(f reflect.Value) (string, error) {
	switch f.Kind() {
	case reflect.Int:
		return strconv.FormatInt(f.Int(), 10), nil

	case reflect.Bool:
		return strconv.FormatBool(f.Bool()), nil

	case reflect.String:
		return f.String(), nil

	default:
		return "", fmt.Errorf("cannot handle value of kind %v", f.Kind())
	}
}

func setValue(v reflect.Value, d string) error {
	switch v.Kind() {
	case reflect.Int:
		i, err := strconv.ParseInt(d, 10, 64)

		if err != nil {
			return err
		}

		v.SetInt(i)

		return nil

	case reflect.Bool:
		b, err := strconv.ParseBool(d)

		if err != nil {
			return err
		}

		v.SetBool(b)

		return nil

	case reflect.String:
		v.SetString(d)

		return nil

	default:
		return fmt.Errorf("cannot handle value of kind %v", v.Kind())
	}
}
