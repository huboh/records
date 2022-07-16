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

import Records

```go
package main

import (
  "github.com/huboh/records"
)
```

### Unmarshaling CSV records

```go
// error handling omited

import (
  "encoding/csv"

  "github.com/huboh/records"
)

// create a struct that represent your csv data.
// NB: unexported fields or fields without csv struct tags are ignored.
type CsvFileEntries struct {
  Age        int    `csv:"age"`
  Name       string `csv:"name"`
  IsEmployee bool   `csv:"isEmployee"`
}

func main() {
  var (
    r, err          = os.Open(csvFile)
    csvReader       = csv.NewReader(r)
    csvRecords, err = csvReader.ReadAll()

    entries = []CsvFileEntries{} // initialize variable to hold csv entries
  )

  entriesErr := records.Unmarshal(csvRecords, &entries)

  // entries => []CsvFileEntries{{20, "john", false}, {20, "mary", true}, {24, "saint", false}, {30, "helen", true}}
}
```

### Marshaling CSV entries

```go
// error handling omited

import (
  "encoding/csv"

  "github.com/huboh/records"
)

// create a struct that represent your csv data.
// NB: unexported fields or fields without csv struct tags are ignored.
type CsvFileEntries struct {
  Age        int    `csv:"age"`
  Name       string `csv:"name"`
  IsEmployee bool   `csv:"isEmployee"`
}

func main() {
  entries := []CsvFileEntries{
    {20, "john", false},
    {20, "mary", true},
    {24, "saint", false},
    {30, "helen", true},
  }

  csvRecords, entriesErr := records.Marshal(entries)

  // csvRecords => [][]string{{"age", "name", "isEmployee"},{"20", "john", "false"},{"20", "mary", "true"},{"24", "saint", "false"},{"30", "helen", "true"}}
}
```

## Contributions

Contributions are welcome to this project to further improve it to suit the public need. I hope you enjoy the simplicity of this package.

## License

This package is provided under MIT license.
