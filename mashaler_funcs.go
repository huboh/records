package marshaler

import (
	"fmt"
	"reflect"
	"strconv"
)

const (
	csvKeyName = string("csv")
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

// forEachStruct iterates through a `reflect.Value` representation of slice of structs & calls the given function for each struct
func forEachStruct(slice reflect.Value, f func(s reflect.Value, i int)) {
	if (slice.Kind() == reflect.Slice) && (slice.Type().Elem().Kind() == reflect.Struct) {
		for i := 0; i < slice.Len(); i++ {
			f(slice.Index(i), i)
		}
	}
}

// forEachStructField iterates through a `reflect.Type` representation of a  structs & calls the given function for each struct field
func forEachStructField(s reflect.Type, f func(f reflect.StructField, i int)) {
	if s.Kind() == reflect.Struct {
		for i := 0; i < s.NumField(); i++ {
			f(s.Field(i), i)
		}
	}
}

func marshalHeader(s reflect.Type) []string {
	row := make([]string, 0, s.NumField())

	forEachStructField(s, func(sf reflect.StructField, i int) {
		if csvKey, csvKeyExists := sf.Tag.Lookup(csvKeyName); csvKeyExists {
			row = append(row, csvKey)
		}
	})

	return row
}

func marshalStruct(sv reflect.Value) (csvRow []string, err error) {
	row := make([]string, 0, sv.NumField())

	forEachStructField(sv.Type(), func(sf reflect.StructField, i int) {
		sfv := sv.Field(i)
		_, csvKeyExists := sf.Tag.Lookup(csvKeyName)

		if csvKeyExists {
			if v, e := getValue(sfv); e != nil {
				err = e
			} else {
				row = append(row, v)
			}
		}
	})

	return row, err
}

func unmarshalStruct(row []string, csvKeyMap map[string]int, sv reflect.Value) (err error) {
	forEachStructField(sv.Type(), func(sf reflect.StructField, i int) {
		f := sv.Field(i)
		csvKey := sf.Tag.Get(csvKeyName)

		if fieldPosition, fieldExists := csvKeyMap[csvKey]; fieldExists {
			if e := setValue(f, row[fieldPosition]); e != nil {
				err = e
			}
		}
	})

	return err
}
