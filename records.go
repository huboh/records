// Package records marshals & unmarshals CSV records(slice of slice of strings) to/from CSV entries(slice of structs).
package records

import (
	"errors"
	"reflect"

	"github.com/huboh/records/internal/slice"
)

// Marshal maps all of the CSV entries(slice of structs) v to CSV records(slice of slice of strings).
// the first record in the output is the csv header with the records title, gotten from the csv struct tag.
//
// NB:
//
// - struct fields without csv tags are ignored.
//
//  // marshal struct to csv records
//  entries := []CsvFileEntries{...}
//  csvRecords, err := records.Marshal(entries)
//
//  if err != nil && error.Is(err, records.ErrUnSupportedKind) {
//    // encountered fields with unsupported types
//    // supported types are: all int, uint & float types, bool, string)
//  }
//  // success
//
//  // write csv records to file
//  w, e := os.Create("")
//  csvWriter := csv.NewWriter(w)
//  csvWriterErr := csvWriter.WriteAll(csvRecords)
func Marshal(v any) (records [][]string, err error) {
	if !isSliceOfStruct(v) {
		return nil, errors.New("v must be a slice of structs")
	}

	var (
		slice      = reflect.ValueOf(v)
		recordType = slice.Type().Elem()
		csvRecords = make([][]string, 0, slice.Len())
	)

	csvRecords = append(csvRecords, getEntryTags(recordType))

	forEachStruct(slice, func(s reflect.Value, i int) {
		record, recordErr := marshalEntry(s)

		if recordErr != nil {
			err = recordErr
		}

		csvRecords = append(csvRecords, record)
	})

	return csvRecords, err
}

// Unmarshal transforms CSV records(slice of slice of strings) to CSV entries(slice of structs).
// the first csv record is assumed to be the csv header names, it builds each CSV entry by mapping it's csv struct tags to the column names.
//
// NB: unexported struct fields & struct fields without csv tags are ignored.
//
//  // reading csv from file
//  r, err := os.Open("")
//  csvReader := csv.NewReader(r)
//  csvRecords, err := csvReader.ReadAll()
//
//  // unmarshal csv records to structs
//  csvEntries := []CsvFileEntries{} // initialize variable to hold csv entries
//  err := records.Unmarshal(csvRecords, &csvEntries)
//
//  // handle possible errors
//  if err != nil && error.Is(err, records.ErrUnSupportedKind) {
// 	  // encountered fields with unsupported types
// 	  // supported types are: all int, uint & float types, bool, string
//  }
//
//  // success
func Unmarshal(records [][]string, v any) (err error) {
	if !isPointerToSliceOfStructs(v) {
		return errors.New("v must be a pointer to a slice of structs")
	}

	var (
		sliceVal     = reflect.ValueOf(v).Elem()
		recordType   = reflect.ValueOf(v).Elem().Type().Elem()
		recordKeys   = records[0]
		recordKeyMap = make(map[string]int, len(recordKeys))
	)

	for csvHeaderkeyIndex, csvHeaderKey := range recordKeys {
		recordKeyMap[csvHeaderKey] = csvHeaderkeyIndex
	}

	slice.ForEach(records[1:], func(record []string, index int) {
		value := reflect.New(recordType).Elem()
		unmarshalErr := unmarshalRecord(record, recordKeyMap, value)

		if unmarshalErr != nil {
			err = unmarshalErr
		}

		sliceVal.Set(reflect.Append(sliceVal, value))
	})

	return err
}
