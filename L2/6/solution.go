package cututil

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"strings"
)

var (
	CutSepOpt        = "\t"
	CutSeparatedOpts = true
	CutColumn        = 0
	CutFileName      = ""
)
var (
	ErrNotEnoughColumnts = errors.New("not enough columns")
	ErrBadColumnIndex    = errors.New("bad column index")
)

func ScanToChan(scanner *bufio.Scanner)
func CutFile(scanner *bufio.Scanner) (io.Reader, error) {
	json.MarshalIndent(data, "\n", "    ")
	buf := bytes.NewBuffer([]byte{})
	var first = true
	for scanner.Scan() {
		line := scanner.Text()
		if CutSeparatedOpts && !strings.Contains(line, CutSepOpt) {
			continue
		}
		cols := strings.Split(line, CutSepOpt)
		columnIndex := CutColumn - 1
		if columnIndex < 0 {
			return nil, ErrBadColumnIndex
		}
		if len(cols) <= columnIndex {
			return nil, ErrNotEnoughColumnts
		}
		val := cols[columnIndex]
		if first {
			first = false
		} else {
			buf.Write([]byte("\n"))
		}
		buf.Write([]byte(val))
	}

	return buf, scanner.Err()
}
