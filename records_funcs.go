package records

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"github.com/huboh/records/internal/structType"
)

const (
	csvTagName = string("csv")
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
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(f.Int(), 10), nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(f.Uint(), 10), nil

	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(f.Float(), 'f', -1, 64), nil

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
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(v, 10, 64)

		if err != nil {
			return err
		}

		d.SetInt(i)

		return nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		ui, err := strconv.ParseUint(v, 10, 64)

		if err != nil {
			return err
		}

		d.SetUint(ui)

		return nil

	case reflect.Float32, reflect.Float64:
		ui, err := strconv.ParseFloat(v, 64)

		if err != nil {
			return err
		}

		d.SetFloat(ui)

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

// getEntryTags gets the csv struct tag from each entry field
func getEntryTags(s reflect.Type) (entryTags []string) {
	return structType.Reduce(s, entryTags, func(tags []string, structField reflect.StructField, i int) []string {
		if tag, tagExists := structField.Tag.Lookup(csvTagName); tagExists {
			tags = append(tags, tag)
		}

		return tags
	})
}

// marshalEntry transform a CSV entry to a CSV record, it ignores struct fields without a csv struct tag
func marshalEntry(sv reflect.Value) (record []string, err error) {
	entryReducer := func(csvRecord []string, structField reflect.StructField, index int) []string {
		field := sv.Field(index)
		fieldName := structField.Name
		_, tagExists := structField.Tag.Lookup(csvTagName)

		if tagExists {
			record, recordErr := getValue(field)
			csvRecord = append(csvRecord, record)

			if recordErr != nil {
				err = recordErr

				if errors.Is(recordErr, ErrUnSupportedKind) {
					err = KindErr{
						WrappedErr: recordErr,
						Message:    fmt.Sprintf("could not parse '%v', its kind '%v' is not supported", fieldName, field.Kind()),
					}
				}
			}
		}

		return csvRecord
	}

	return structType.Reduce(sv.Type(), record, entryReducer), err
}

// unmarshalRecord transforms a CSV record to a CSV entry, it ignores unexported fields & fields without csv struct tag
func unmarshalRecord(record []string, tagMap map[string]int, sv reflect.Value) (err error) {
	structType.ForEach(sv.Type(), func(sf reflect.StructField, i int) {
		field := sv.Field(i)
		fieldKind := field.Kind()
		fieldName := sf.Name

		if tagPos, tagExists := tagMap[sf.Tag.Get(csvTagName)]; tagExists && field.CanSet() {
			e := setValue(field, record[tagPos])

			if e != nil {
				err = e

				if errors.Is(e, ErrUnSupportedKind) {
					err = KindErr{e, fmt.Sprintf("could not set '%v', its kind '%v' is not supported", fieldName, fieldKind)}
				}
			}
		}
	})

	return err
}
