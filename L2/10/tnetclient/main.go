package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

func main() {
	if len(os.Args) < 3 {
		os.Exit(1)
	}
	const network = "tcp"
	var addr = os.Args[1] //"localhost"
	var port = os.Args[2] //"5555"
	conn, err := net.Dial(network, addr+":"+port)
	if nil != err {
		panic(err)
	}
	ClientRun(os.Stdin, os.Stdout, conn)
}

func ClientRun(stdin io.Reader, stdout io.Writer, rw io.ReadWriteCloser) {
	setTimer := make(chan struct{})
	stopTimer := make(chan struct{})
	go func() {
		timer := time.NewTicker(time.Second)
		timer.Stop()
		for {
			select {
			case <-setTimer:
				timer.Reset(time.Second * 10)
			case <-stopTimer:
				timer.Stop()
			case <-timer.C:
				fmt.Println("-TIMEOUT-")
				rw.Close()
				os.Exit(1)
			}
		}
	}()
	go func(writer io.Writer, reader io.Reader) {
		var buffer [1]byte // Seems like the length of the buffer needs to be small, otherwise will have to wait for buffer to fill up.
		p := buffer[:]

		for {
			// Read 1 byte.
			n, err := reader.Read(p)
			if n <= 0 && nil == err {
				continue
			} else if n <= 0 && nil != err {
				break
			}
			stopTimer <- struct{}{}
			writer.Write(p)
		}
	}(stdout, rw)

	var buffer bytes.Buffer
	var p []byte

	crlf := []byte{'\r', '\n'}

	scanner := bufio.NewScanner(stdin)

	for scanner.Scan() {
		nb := scanner.Bytes()
		if len(nb) == 1 {
			if []rune(string(nb))[0] == 4 {
				fmt.Println("-END OF TRANSMISSION-")
				rw.Close()
				break
			}
		}

		buffer.Write(nb)
		buffer.Write(crlf)

		p = buffer.Bytes()
		_, err := rw.Write(p)
		if nil != err {
			break
		}
		setTimer <- struct{}{}
		buffer.Reset()
	}
}
