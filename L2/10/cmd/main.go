package main

import (
	"os"
	"telnet"
)

func main() {
	if len(os.Args) < 3 {
		os.Exit(1)
	}
	telnet.Host = os.Args[1]
	telnet.Port = os.Args[2]
	telnet.Con(telnet.ReadNextInput(">"))
}
