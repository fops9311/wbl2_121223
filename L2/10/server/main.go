package main

import (
	"time"

	"github.com/reiver/go-telnet"
)

func main() {

	var handler telnet.Handler = modEchoHandler{}

	err := telnet.ListenAndServe(":5555", handler)
	if nil != err {
		//@TODO: Handle this error better.
		panic(err)
	}
}

type modEchoHandler struct {
}

func (handler modEchoHandler) ServeTELNET(ctx telnet.Context, w telnet.Writer, r telnet.Reader) {

	var buffer [1]byte // Seems like the length of the buffer needs to be small, otherwise will have to wait for buffer to fill up.
	p := buffer[:]

	for {
		n, err := r.Read(p)
		str := string(p[:n])
		time.Sleep(time.Millisecond * 100)
		if n > 0 {
			w.Write([]byte(str))
		}

		if nil != err {
			break
		}
	}
}
