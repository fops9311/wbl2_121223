package sortutil

import (
	"fmt"
	"sort"
	"testing"
)

type TestSortOpts struct {
	UniqueSortOpt          bool
	ReverseSortOpt         bool
	CaseSensetiveSortOpt   bool
	ColumnSortOpt          int
	ColumnSeparatorSortOpt string
	NumberSortOpt          bool
	FileNameSortOpt        string
	DebugSortOpt           bool
}

func TestSortFile(t *testing.T) {
	var testCases []SortFileTestCase = []SortFileTestCase{
		{
			opts: TestSortOpts{
				UniqueSortOpt:          true,
				ReverseSortOpt:         false,
				CaseSensetiveSortOpt:   false,
				ColumnSortOpt:          1,
				ColumnSeparatorSortOpt: " ",
				NumberSortOpt:          false,
				FileNameSortOpt:        "test_data.txt",
				DebugSortOpt:           true,
			},
			want: []string{"RedHat 1", "data 2", "laptop 2", "computer 3", "mouse 4", "debian 4", "LAPTOP 6", "laptop 8"},
		},
	}
	for _, test := range testCases {
		ApplyTestOpts(test.opts)

		test.got = SortFile()
		AssertString(fmt.Sprintf("%v", test.got), fmt.Sprintf("%v", test.want), t)
	}
}

type SortFileTestCase struct {
	opts TestSortOpts
	got  []string
	want []string
}

func TestSortBy(t *testing.T) {
	var testCases []SortByTestCase = []SortByTestCase{
		{
			opts: TestSortOpts{
				NumberSortOpt:  true,
				ReverseSortOpt: false,
			},
			inS:  []string{"2", "1"},
			want: []string{"1", "2"},
		},
		{
			opts: TestSortOpts{
				NumberSortOpt:  true,
				ReverseSortOpt: true,
			},
			inS:  []string{"1", "2"},
			want: []string{"2", "1"},
		},
	}
	for _, test := range testCases {
		ApplyTestOpts(test.opts)
		sort.Sort(test.inS)
		test.got = (test.inS)
		AssertString(fmt.Sprintf("%v", test.got), fmt.Sprintf("%v", test.want), t)
	}
}

type SortByTestCase struct {
	opts TestSortOpts
	inS  sortBy
	got  []string
	want []string
}

func TestApplyFilterOpts(t *testing.T) {
	var testCases []ApplyFilterOptsTestCase = []ApplyFilterOptsTestCase{
		{
			opts: TestSortOpts{
				UniqueSortOpt: false,
			},
			inS:  []string{"1", "2"},
			want: []string{"1", "2"},
		},
		{
			opts: TestSortOpts{
				UniqueSortOpt: true,
			},
			inS:  []string{"1", "1"},
			want: []string{"1"},
		},
	}
	for _, test := range testCases {
		ApplyTestOpts(test.opts)
		test.got = applyFilterOpts(test.inS)
		AssertString(fmt.Sprintf("%v", test.got), fmt.Sprintf("%v", test.want), t)
	}
}

type ApplyFilterOptsTestCase struct {
	opts TestSortOpts
	inS  []string
	got  []string
	want []string
}

func TestApplyCompareOpts(t *testing.T) {
	var testCases []ApplyCompareOptsTestCase = []ApplyCompareOptsTestCase{
		{
			opts: TestSortOpts{
				NumberSortOpt: false,
			},
			inStr1: "1",
			inStr2: "2",
			want:   true,
		},
		{
			opts: TestSortOpts{
				NumberSortOpt: true,
			},
			inStr1: "1",
			inStr2: "2",
			want:   true,
		},
		{
			opts: TestSortOpts{
				NumberSortOpt: false,
			},
			inStr1: "2",
			inStr2: "1",
			want:   false,
		},
		{
			opts: TestSortOpts{
				NumberSortOpt: true,
			},
			inStr1: "2",
			inStr2: "1",
			want:   false,
		},
	}
	for _, test := range testCases {
		ApplyTestOpts(test.opts)
		test.got = applyCompareOpts(test.inStr1, test.inStr2)
		AssertBool(test.got, test.want, t)
	}
}

type ApplyCompareOptsTestCase struct {
	opts   TestSortOpts
	inStr1 string
	inStr2 string
	got    bool
	want   bool
}

func TestCompareWithDirection(t *testing.T) {
	var testCases []CompareWithDirectionTestCase = []CompareWithDirectionTestCase{
		{
			opts: TestSortOpts{
				ReverseSortOpt: false,
			},
			in1:    1,
			in2:    2,
			inStr1: "1",
			inStr2: "2",
			want:   true,
		},
		{
			opts: TestSortOpts{
				ReverseSortOpt: true,
			},
			in1:    1,
			in2:    2,
			inStr1: "1",
			inStr2: "2",
			want:   false,
		},
	}

	for _, test := range testCases {
		ApplyTestOpts(test.opts)
		test.got = compareWithDirection(test.in1, test.in2)
		AssertBool(test.got, test.want, t)
		test.got = compareWithDirection(test.inStr1, test.inStr2)
		AssertBool(test.got, test.want, t)
	}
}

type CompareWithDirectionTestCase struct {
	opts   TestSortOpts
	in1    int64
	in2    int64
	inStr1 string
	inStr2 string
	got    bool
	want   bool
}

func TestApplyDataOpts(t *testing.T) {
	var testCases []ApplyDataOptsTestCase = []ApplyDataOptsTestCase{
		{
			opts: TestSortOpts{
				ColumnSortOpt:          0,
				ColumnSeparatorSortOpt: " ",
				CaseSensetiveSortOpt:   false,
			},
			in:   "a B c1 D e F",
			want: "a",
		},
		{
			opts: TestSortOpts{
				ColumnSortOpt:          1,
				ColumnSeparatorSortOpt: " ",
				CaseSensetiveSortOpt:   false,
			},
			in:   "a B c1 D e F",
			want: "b",
		},
		{
			opts: TestSortOpts{
				ColumnSortOpt:          1,
				ColumnSeparatorSortOpt: " ",
				CaseSensetiveSortOpt:   true,
			},
			in:   "a B c1 D e F",
			want: "B",
		},
		{
			opts: TestSortOpts{
				ColumnSortOpt:          -1,
				ColumnSeparatorSortOpt: " ",
				CaseSensetiveSortOpt:   true,
			},
			in:   "a B c1 D e F",
			want: "a B c1 D e F",
		},
		{
			opts: TestSortOpts{
				ColumnSortOpt:          -1,
				ColumnSeparatorSortOpt: " ",
				CaseSensetiveSortOpt:   false,
			},
			in:   "a B c1 D e F",
			want: "a b c1 d e f",
		},
	}
	for _, test := range testCases {
		ApplyTestOpts(test.opts)
		test.got = applyDataOpts(test.in)
		AssertString(test.got, test.want, t)
	}
}
func ApplyTestOpts(o TestSortOpts) {
	UniqueSortOpt = o.UniqueSortOpt
	ReverseSortOpt = o.ReverseSortOpt
	CaseSensetiveSortOpt = o.CaseSensetiveSortOpt
	ColumnSortOpt = o.ColumnSortOpt
	ColumnSeparatorSortOpt = o.ColumnSeparatorSortOpt
	NumberSortOpt = o.NumberSortOpt
	FileNameSortOpt = o.FileNameSortOpt
	DebugSortOpt = o.DebugSortOpt
}

type ApplyDataOptsTestCase struct {
	opts TestSortOpts
	in   string
	got  string
	want string
}

func AssertBool(got, want bool, t *testing.T) {
	if got != want {
		t.Errorf("FAIL VALUE ASSERTION: <%v> == <%v>", got, want)
	} else {
		t.Logf("PASS VALUE ASSERTION: <%v> == <%v>", got, want)
	}
}
func AssertInt64(got, want int64, t *testing.T) {
	if got != want {
		t.Errorf("FAIL VALUE ASSERTION: <%d> == <%d>", got, want)
	} else {
		t.Logf("PASS VALUE ASSERTION: <%d> == <%d>", got, want)
	}
}
func AssertString(got, want string, t *testing.T) {
	if got != want {
		t.Errorf("FAIL VALUE ASSERTION: <%s> == <%s>", got, want)
	} else {
		t.Logf("PASS VALUE ASSERTION: <%s> == <%s>", got, want)
	}
}
func AssertError(got, want error, t *testing.T) {
	if got != want {
		t.Errorf("FAIL ERROR ASSERTION: <%v> == <%v>", got, want)
	} else {
		t.Logf("PASS ERROR ASSERTION: <%v> == <%v>", got, want)
	}
}
