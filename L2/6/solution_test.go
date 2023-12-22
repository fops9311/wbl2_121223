package cututil

import (
	"bufio"
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestCutFile(t *testing.T) {
	testMessages := []string{
		"table_col1\ttable_col2\ttable_col3\ttable_col4",
		"col1\tcol2\tcol3\tcol4",
		"col1",
	}
	testMessage := strings.Join(testMessages, "\n")

	var testCases []CutFileTestCase = []CutFileTestCase{
		{
			Opts: CutFileTestCaseOpts{
				CutColumn:        1,
				CutSepOpt:        "\t",
				CutSeparatedOpts: false,
			},
			want:    "table_col1\ncol1\ncol1",
			errwant: nil,
		},
		{
			Opts: CutFileTestCaseOpts{
				CutColumn:        1,
				CutSepOpt:        "\t",
				CutSeparatedOpts: true,
			},
			want:    "table_col1\ncol1",
			errwant: nil,
		},
		{
			Opts: CutFileTestCaseOpts{
				CutColumn: -1,
				CutSepOpt: "\t",
			},
			want:    "",
			errwant: ErrBadColumnIndex,
		},
		{
			Opts: CutFileTestCaseOpts{
				CutColumn: 10,
				CutSepOpt: "\t",
			},
			want:    "",
			errwant: ErrNotEnoughColumnts,
		},
	}

	for i, test := range testCases {
		t.Logf("TEST CASE %d", i)
		scanner := bufio.NewScanner(bytes.NewBuffer([]byte(testMessage)))
		test.Apply()
		var r io.Reader
		r, test.errgot = CutFile(scanner)
		AssertError(test.errgot, test.errwant, t)
		if r == nil {
			AssertString(test.got, test.want, t)
			continue
		}
		b, err := io.ReadAll(r)
		AssertError(err, nil, t)
		test.got = string(b)
		AssertString(test.got, test.want, t)
	}
}

type CutFileTestCase struct {
	Opts    CutFileTestCaseOpts
	want    string
	got     string
	errwant error
	errgot  error
}

func (tc CutFileTestCase) Apply() {
	CutColumn = tc.Opts.CutColumn
	CutSepOpt = tc.Opts.CutSepOpt
	CutSeparatedOpts = tc.Opts.CutSeparatedOpts
}

type CutFileTestCaseOpts struct {
	CutSepOpt        string
	CutColumn        int
	CutSeparatedOpts bool
}

func AssertBool(got, want bool, t *testing.T) {
	if got != want {
		t.Errorf("\t|FAIL| VALUE ASSERTION: <%v> == <%v>", got, want)
	} else {
		t.Logf("\t|PASS| VALUE ASSERTION: <%v> == <%v>", got, want)
	}
}
func AssertInt64(got, want int64, t *testing.T) {
	if got != want {
		t.Errorf("\t|FAIL| VALUE ASSERTION: <%d> == <%d>", got, want)
	} else {
		t.Logf("\t|PASS| VALUE ASSERTION: <%d> == <%d>", got, want)
	}
}
func AssertString(got, want string, t *testing.T) {
	if got != want {
		t.Errorf("\t|FAIL| VALUE ASSERTION: <%s> == <%s>", got, want)
	} else {
		t.Logf("\t|PASS| VALUE ASSERTION: <%s> == <%s>", got, want)
	}
}
func AssertError(got, want error, t *testing.T) {
	if got != want {
		t.Errorf("\t|FAIL| ERROR ASSERTION: <%v> == <%v>", got, want)
	} else {
		t.Logf("\t|PASS| ERROR ASSERTION: <%v> == <%v>", got, want)
	}
}
