package timeprinter

import (
	"fmt"
	"os"
	tm "time"

	ntp "github.com/beevik/ntp"
)

func PrintTime() {
	time, err := ntp.Time("0.beevik-ntp.pool.ntp.org")
	var now = tm.Now()
	out := fmt.Sprintf("%v / %v", now, time)
	if err != nil {
		_, err = os.Stderr.Write([]byte(err.Error()))
		if err != nil {
			os.Exit(1)
		}
		os.Exit(0)
	}
	fmt.Println(out)
}
