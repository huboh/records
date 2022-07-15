package marshaler

import (
	"fmt"
	"reflect"
	"testing"
)

type testcase struct {
	expectation bool
	dataType    string
	data        any
}

type predicateFunc func(v any) bool

type runPredicateTestProp struct {
	predicateTestType string
	predicateFunc     predicateFunc
	testCases         []testcase
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

		testCases: []testcase{
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

		testCases: []testcase{
			{dataType: "int", data: int(1), expectation: false},
			{dataType: "string", data: string("1"), expectation: false},
			{dataType: "*[]string", data: &[]string{}, expectation: false},
			{dataType: "*struct{}", data: &struct{}{}, expectation: false},
			{dataType: "[]struct{}", data: []struct{}{}, expectation: false},
			{dataType: "map[string]int", data: map[string]int{}, expectation: false},
			{dataType: "struct{}", data: struct{}{}, expectation: true},
			{dataType: "TestCase{}", data: testcase{}, expectation: true},
		},
	})
}

func Test_isPointer(t *testing.T) {
	runPredicateTest(runPredicateTestProp{
		t:                 t,
		predicateFunc:     isPointer,
		predicateTestType: "pointer",

		testCases: []testcase{
			{dataType: "int", data: int(1), expectation: false},
			{dataType: "string", data: string("1"), expectation: false},
			{dataType: "TestCase{}", data: testcase{}, expectation: false},
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

		testCases: []testcase{
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

		testCases: []testcase{
			{dataType: "int", data: (int)(1), expectation: false},
			{dataType: "struct{}", data: struct{}{}, expectation: false},
			{dataType: "[2]string", data: [2]string{}, expectation: false},
			{dataType: "map[string]int", data: map[string]int{}, expectation: false},
			{dataType: "*[]string", data: &[]string{}, expectation: false},
			{dataType: "*[3]TestCase{}", data: &[2]testcase{}, expectation: false},
			{dataType: "*[]struct{}", data: &[]struct{}{}, expectation: true},
			{dataType: "*[]TestCase{}", data: &[]testcase{}, expectation: true},
		},
	})
}

//

func Test_getValue(t *testing.T) {
	type GetValueTestCase struct {
		expectation string
		name        string
		data        any
	}

	testCases := []GetValueTestCase{
		{name: "integer test", data: 120, expectation: "120"},
		{name: "boolean test", data: false, expectation: "false"},
		{name: "string test", data: "hello", expectation: "hello"},
		{name: "unsupported type test", data: []string{}, expectation: ""},
	}

	for _, tc := range testCases {
		var (
			testName            = tc.name
			expectation         = tc.expectation
			testfuncData        = reflect.ValueOf(tc.data)
			testResult, testErr = getValue(testfuncData)
		)

		t.Run(testName, func(t *testing.T) {
			if testResult != expectation {
				t.Error(testErr)
			}
		})
	}
}

func Test_setValue(t *testing.T) {
	var (
		sData     = ""
		sValue    = "test"
		sPtrValue = reflect.ValueOf(&sData)
		sTestErr  = setValue(sPtrValue.Elem(), sValue)
	)

	if sData != sValue || sTestErr != nil {
		t.Error(sTestErr)
	}

	//

	var (
		b         = false
		bValue    = "true"
		bPtrValue = reflect.ValueOf(&b)
		bTestErr  = setValue(bPtrValue.Elem(), bValue)
	)

	if b != true || bTestErr != nil {
		t.Error(bTestErr)
	}

	//

	var (
		i         = 100
		iValue    = 120
		iPtrValue = reflect.ValueOf(&i)
		iTestErr  = setValue(iPtrValue.Elem(), fmt.Sprint(iValue))
	)

	if i != iValue || iTestErr != nil {
		t.Error(iTestErr)
	}

	//

	var (
		st         = struct{}{}
		stPtrValue = reflect.ValueOf(&st)
		stTestErr  = setValue(stPtrValue.Elem(), "not gonna work")
	)

	if stTestErr == nil {
		t.Error("expected error, got nil")
	}
}