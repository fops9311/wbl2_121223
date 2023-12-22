package greputil

import (
	"testing"
)

type GrepFileTestCase struct {
	Opts
}
type Opts struct {
	FileNameGrepOpt string
	DebugOpt        bool
}

func (o Opts) Apply() {
	FileNameGrepOpt = o.FileNameGrepOpt
	DebugOpt = o.DebugOpt
}

func TestGrepFile(t *testing.T) {

	var testCases []GrepFileTestCase = []GrepFileTestCase{
		{
			Opts{
				FileNameGrepOpt: "test_data.txt",
				DebugOpt:        true,
			},
		},
	}

	for _, test := range testCases {
		test.Apply()
		GrepFile()
		t.Errorf("NO TEST YET")
	}
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
