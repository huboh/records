# Records

## A light weight & fast Go data marshaler

Records is a light weight Go package that marshals & unmarshals CSV records`(slice of slice of strings)` to/from CSV entries`(slice of structs)`.

Go provides the `csv.NewReader` & `csv.NewWriter` functions to read a CSV file into a *slice of slice of strings* & to write a *slice of slice of strings* out to a CSV file, without a way to map that data to the fields in a struct. this package add that missing functionality.

## Installation

Once Go is installed, run the following command to get `Records`.

```bash
go get github.com/huboh/records
```

## Usage

### Unmarshaling CSV records

import Records

```go
package main

import (
  "github.com/huboh/records"
)
```

get csv records from a data source, most likely from a file

```go
r, err := os.Open("") // import "os"
csvReader := csv.NewReader(r) // import "encoding/csv"
csvRecords, err := csvReader.ReadAll() // read csv records
```

unmarshal csv records. note that unexported fields or fields without csv struct tags are ignored.

```go
// create a struct that represent your csv data.
type CsvFileEntries struct {
  Age        int    `csv:"age"`
  Name       string `csv:"name"`
  IsEmployee bool   `csv:"isEmployee"`
}

entries := []CsvFileEntries{} // initialize variable to hold csv records
err := records.Unmarshal(csvRecords, &entries)

if err != nil && errors.Is(err, records.ErrUnSupportedKind) {
  // encountered fields with unsupported types
  // supported types are: all int, uint & float types, bool, string.
}

// success..  entries => []CsvFileEntries{{...}}
```

### Marshaling CSV entries

import Records

```go
package main

import (
  "github.com/huboh/records"
)
```

transform csv data to csv records. note that fields without csv struct tags are ignored.

```go
// create a struct that represent your csv data.
type CsvFileEntries struct {
  Age        int    `csv:"age"`
  Name       string `csv:"name"`
  IsEmployee bool   `csv:"isEmployee"`
}

entries := []CsvFileEntries{...}
csvRecords, err := records.Marshal(entries) // transform csv data to csv records

if err != nil && error.Is(err, records.ErrUnSupportedKind) {
  // encountered fields with unsupported types
  // supported types are: all int, uint & float types, bool, string
}

// success..  csvRecords => [][]string{{...}}
```

after successfully marshalling your csv data to records, writing the csv records to a file is as easy as:

```go
w, e := os.Create("")
csvWriter := csv.NewWriter(w)
csvWriterErr := csvWriter.WriteAll(csvRecords) // writes csv records to file
```

## Contributions

Contributions are welcome to this project to further improve it to suit the public need. I hope you enjoy the simplicity of this package.

## License

This package is provided under MIT license.
