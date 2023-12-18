package unpack

import "testing"

func TestUnpack(t *testing.T) {
	var testCases []struct {
		in   string
		gerr error
		werr error
		g    string
		w    string
	} = []struct {
		in   string
		gerr error
		werr error
		g    string
		w    string
	}{
		{
			in:   "a\\bs",
			werr: nil,
			w:    "abs",
		},
		{
			in:   "a",
			werr: nil,
			w:    "a",
		},
		{
			in:   "1a",
			werr: nil,
			w:    "a",
		},
		{
			in:   "2a",
			werr: nil,
			w:    "aa",
		},
		{
			in:   "3a",
			werr: nil,
			w:    "aaa",
		},
		{
			in:   "4a",
			werr: nil,
			w:    "aaaa",
		},
		{
			in:   "5a",
			werr: nil,
			w:    "aaaaa",
		},
		{
			in:   "6a",
			werr: nil,
			w:    "aaaaaa",
		},
		{
			in:   "7a",
			werr: nil,
			w:    "aaaaaaa",
		},
		{
			in:   "8a",
			werr: nil,
			w:    "aaaaaaaa",
		},
		{
			in:   "9a",
			werr: nil,
			w:    "aaaaaaaaa",
		},
		{
			in:   "10a",
			werr: nil,
			w:    "aaaaaaaaaa",
		},
		{
			in:   "0a",
			werr: ErrInvalidPattern,
			w:    "",
		},
		{
			in:   "00a",
			werr: ErrInvalidPattern,
			w:    "",
		},
		{
			in:   "5",
			werr: ErrInvalidPattern,
			w:    "",
		},
		{
			in:   "2\\a",
			werr: nil,
			w:    "aa",
		},
		{
			in:   "2\\\\",
			werr: nil,
			w:    "\\\\",
		},
		{
			in:   "3\\2",
			werr: nil,
			w:    "222",
		},
		{
			in:   "3\\A",
			werr: nil,
			w:    "AAA",
		},
		{
			in:   "3\\24\\3",
			werr: nil,
			w:    "2223333",
		},
		{
			in:   "",
			werr: ErrEmptyString,
			w:    "",
		},
		{
			in:   "5",
			werr: ErrInvalidPattern,
			w:    "",
		},
		{
			in:   "5\\",
			werr: ErrInvalidPattern,
			w:    "",
		},
	}
	for _, test := range testCases {
		test.g, test.gerr = Unpack(test.in)
		AssertError(test.gerr, test.werr, t)
		AssertString(test.g, test.w, t)
	}
}
func TestIterate(t *testing.T) {
	in := "abc"
	l := AsLinkedListOfRunes(in)
	got := []rune{}
	_ = l.Iterate(func(s *Simbol) error {
		got = append(got, s.Val)
		return nil
	})
	AssertString(string(got), in, t)
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
