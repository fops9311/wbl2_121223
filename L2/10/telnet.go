package telnet

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

var Host = "rpcx.site"
var Port = "80"

func Con(command string) { // Establishing a connection
	conn, err := net.Dial("tcp", Host+":"+Port)
	if err != nil {
		fmt.Println("dial error:", err)
		return
	}
	defer conn.Close()
	// Send request, http 1.0 protocol
	fmt.Fprintf(conn, command+"\r\n\r\n")
	// Read response
	var sb strings.Builder
	buf := make([]byte, 256)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				fmt.Println("read error:", err)
			}
			break
		}
		sb.Write(buf[:n])
	}
	// Show results
	fmt.Println("response:", sb.String())
	fmt.Println("total response size:", sb.Len())
}

var reader = bufio.NewReader(os.Stdin)

func ReadNextInput(invite string) string {
	fmt.Print(invite)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}
