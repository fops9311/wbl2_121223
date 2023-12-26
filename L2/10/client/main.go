package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

const (
	HOST = "localhost"
	PORT = "5555"
	TYPE = "tcp"
)

func main() {
	tcpServer, err := net.ResolveTCPAddr(TYPE, HOST+":"+PORT)

	if err != nil {
		println("ResolveTCPAddr failed:", err.Error())
		os.Exit(1)
	}

	conn, err := net.DialTCP(TYPE, nil, tcpServer)
	if err != nil {
		println("Dial failed:", err.Error())
		os.Exit(1)
	}
	defer conn.Close()
	for {
		_, err = conn.Write([]byte(ReadNextInput(">")))
		if err != nil {
			println("Write data failed:", err.Error())
			os.Exit(1)
		}
		recieve := make(chan string)
		go func() {
			// buffer to get data
			received := make([]byte, 10)
			_, err = conn.Read(received)
			if err != nil {
				println("Read data failed:", err.Error())
				os.Exit(1)
			}
			recieve <- string(received)
		}()
		select {
		case v := <-recieve:
			println("Received message:", string(v))
		case <-time.NewTimer(time.Second * 1).C:
			println("Timeout!")
		}
	}

}

var reader = bufio.NewReader(os.Stdin)

func ReadNextInput(invite string) string {
	fmt.Print(invite)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}
