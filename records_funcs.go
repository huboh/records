package records

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
		return "", ErrUnSupportedKind
	}
}

func setValue(d reflect.Value, v string) error {
	switch d.Kind() {
	case reflect.Int:
		i, err := strconv.ParseInt(v, 10, 64)

		if err != nil {
			return err
		}

		d.SetInt(i)

		return nil

	case reflect.Bool:
		b, err := strconv.ParseBool(v)

		if err != nil {
			return err
		}

		d.SetBool(b)

		return nil

	case reflect.String:
		d.SetString(v)

		return nil

	default:
		return ErrUnSupportedKind
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

func getRecordKeys(s reflect.Type) (recordKeys []string) {
	row := make([]string, 0, s.NumField())

	forEachStructField(s, func(sf reflect.StructField, i int) {
		if csvKey, csvKeyExists := sf.Tag.Lookup(csvKeyName); csvKeyExists {
			row = append(row, csvKey)
		}
	})

	return row
}

func marshalRecord(sv reflect.Value) (csvRecord []string, err error) {
	record := make([]string, 0, sv.NumField())

	forEachStructField(sv.Type(), func(sf reflect.StructField, i int) {
		sfv := sv.Field(i)
		csvKey, csvKeyExists := sf.Tag.Lookup(csvKeyName)

		if csvKeyExists {
			v, e := getValue(sfv)

			if e != nil {
				err = KindErr{
					Message:    fmt.Sprintf("field '%v' is of unsupported kind: %v", csvKey, sfv.Kind()),
					WrappedErr: e,
				}
			}

			record = append(record, v)
		}
	})

	return record, err
}

func unmarshalRecord(record []string, csvKeyMap map[string]int, sv reflect.Value) (err error) {
	forEachStructField(sv.Type(), func(sf reflect.StructField, i int) {
		f := sv.Field(i)
		csvKey := sf.Tag.Get(csvKeyName)

		if fieldPosition, fieldExists := csvKeyMap[csvKey]; fieldExists && f.CanSet() {
			if e := setValue(f, record[fieldPosition]); e != nil {
				err = KindErr{
					Message:    fmt.Sprintf("field '%v' is of unsupported kind: %v", csvKey, f.Kind()),
					WrappedErr: e,
				}
			}
		}
	})

	return err
}
