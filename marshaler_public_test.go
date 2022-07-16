package records_test

import (
	"encoding/csv"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/huboh/records"
)

const (
	csvTestFile = "./testdata/test.csv"
)

type csvTestFileEntries struct {
	Age        int    `csv:"age"`
	Name       string `csv:"name"`
	IsEmployee bool   `csv:"isEmployee"`
}

func checkErr(e error, f func(e error)) {
	if e != nil {
		f(e)
	}
}

func TestMarshal(t *testing.T) {
	entries := []csvTestFileEntries{
		{20, "john", false},
		{20, "mary", true},
		{24, "saint", false},
		{30, "helen", true},
	}

	csvRecordsExpectation := [][]string{
		{"age", "name", "isEmployee"},
		{"20", "john", "false"},
		{"20", "mary", "true"},
		{"24", "saint", "false"},
		{"30", "helen", "true"},
	}

	csvRecords, entriesErr := records.Marshal(entries)

	checkErr(entriesErr, func(e error) {
		t.Error("unmarshal error", e)
	})

	if diff := cmp.Diff(csvRecords, csvRecordsExpectation); diff != "" {
		t.Error(diff)
	}
}

func TestUnmarshal(t *testing.T) {
	r, err := os.Open(csvTestFile)

	checkErr(err, func(e error) {
		t.Fatal("could not open test file", e)
	})

	csvReader := csv.NewReader(r)
	csvRecords, err := csvReader.ReadAll()

	checkErr(err, func(e error) {
		t.Fatal("could not read test file", e)
	})

	entries := []csvTestFileEntries{}
	entriesErr := records.Unmarshal(csvRecords, &entries)
	extriesExpectation := []csvTestFileEntries{
		{20, "john", false}, {20, "mary", true}, {24, "saint", false}, {30, "helen", true},
	}

	checkErr(entriesErr, func(e error) {
		t.Error("unmarshal error", e)
	})

	if diff := cmp.Diff(entries, extriesExpectation); diff != "" {
		t.Error(diff)
	}
}
