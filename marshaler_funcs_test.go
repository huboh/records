package marshaler

import (
	"fmt"
	"testing"
)

type TestCase struct {
	expectation bool
	dataType    string
	data        interface{}
}

type PredicateFunc func(v any) bool

type runPredicateTestProp struct {
	predicateTestType string
	predicateFunc     PredicateFunc
	testCases         []TestCase
	t                 *testing.T
}

func runPredicateTest(props runPredicateTestProp) {
	for _, testCase := range props.testCases {
		var (
			testDataType    = testCase.dataType
			testResult      = props.predicateFunc(testCase.data)
			testExpectation = testCase.expectation
			testName        = fmt.Sprintf("returned %t for %s", testExpectation, testDataType)
		)

		props.t.Run(testName, func(t *testing.T) {
			if testResult != testExpectation {
				t.Errorf("unexpected result: type %s is not a %v", testDataType, props.predicateTestType)
			}
		})
	}
}

func Test_isSlice(t *testing.T) {
	runPredicateTest(runPredicateTestProp{
		t:                 t,
		predicateFunc:     isSlice,
		predicateTestType: "slice",

		testCases: []TestCase{
			{dataType: "int", data: int(1), expectation: false},
			{dataType: "string", data: string("1"), expectation: false},
			{dataType: "struct{}", data: struct{}{}, expectation: false},
			{dataType: "[2]string", data: [2]string{}, expectation: false},
			{dataType: "*[]string", data: &[]string{}, expectation: false},
			{dataType: "map[string]int", data: map[string]int{}, expectation: false},
			{dataType: "[]int", data: []int{}, expectation: true},
			{dataType: "[]string", data: []string{}, expectation: true},
			{dataType: "[]struct{}", data: []struct{}{}, expectation: true},
		},
	})
}

func Test_isStruct(t *testing.T) {
	runPredicateTest(runPredicateTestProp{
		t:                 t,
		predicateFunc:     isStruct,
		predicateTestType: "struct",

		testCases: []TestCase{
			{dataType: "int", data: int(1), expectation: false},
			{dataType: "string", data: string("1"), expectation: false},
			{dataType: "*[]string", data: &[]string{}, expectation: false},
			{dataType: "*struct{}", data: &struct{}{}, expectation: false},
			{dataType: "[]struct{}", data: []struct{}{}, expectation: false},
			{dataType: "map[string]int", data: map[string]int{}, expectation: false},
			{dataType: "struct{}", data: struct{}{}, expectation: true},
			{dataType: "TestCase{}", data: TestCase{}, expectation: true},
		},
	})
}

func Test_isPointer(t *testing.T) {
	runPredicateTest(runPredicateTestProp{
		t:                 t,
		predicateFunc:     isPointer,
		predicateTestType: "pointer",

		testCases: []TestCase{
			{dataType: "int", data: int(1), expectation: false},
			{dataType: "string", data: string("1"), expectation: false},
			{dataType: "TestCase{}", data: TestCase{}, expectation: false},
			{dataType: "struct{}", data: struct{}{}, expectation: false},
			{dataType: "[]struct{}", data: []struct{}{}, expectation: false},
			{dataType: "map[string]int", data: map[string]int{}, expectation: false},
			{dataType: "*struct{}", data: &struct{}{}, expectation: true},
			{dataType: "*[]string", data: &[]string{}, expectation: true},
		},
	})
}

func Test_isSliceOfStruct(t *testing.T) {
	runPredicateTest(runPredicateTestProp{
		t:                 t,
		predicateFunc:     isSliceOfStruct,
		predicateTestType: "slice of struct",

		testCases: []TestCase{
			{dataType: "int", data: (int)(1), expectation: false},
			{dataType: "struct{}", data: struct{}{}, expectation: false},
			{dataType: "[]struct{}", data: []struct{}{}, expectation: true},
			{dataType: "[2]string", data: [2]string{}, expectation: false},
			{dataType: "*[]string", data: &[]string{}, expectation: false},
			{dataType: "map[string]int", data: map[string]int{}, expectation: false},
		},
	})
}

func Test_isPointerToSliceOfStructs(t *testing.T) {
	runPredicateTest(runPredicateTestProp{
		t:                 t,
		predicateFunc:     isPointerToSliceOfStructs,
		predicateTestType: "pointer to slice of struct",

		testCases: []TestCase{
			{dataType: "int", data: (int)(1), expectation: false},
			{dataType: "struct{}", data: struct{}{}, expectation: false},
			{dataType: "[2]string", data: [2]string{}, expectation: false},
			{dataType: "map[string]int", data: map[string]int{}, expectation: false},
			{dataType: "*[]string", data: &[]string{}, expectation: false},
			{dataType: "*[3]TestCase{}", data: &[2]TestCase{}, expectation: false},
			{dataType: "*[]struct{}", data: &[]struct{}{}, expectation: true},
			{dataType: "*[]TestCase{}", data: &[]TestCase{}, expectation: true},
		},
	})
}
