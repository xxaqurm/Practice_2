package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Connection error:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connected to server")
	fmt.Println("Enter commands (e.g. INSERT my_db users {...})")

	input := bufio.NewReader(os.Stdin)
	server := bufio.NewReader(conn)

	for {
		fmt.Print("> ")

		line, err := input.ReadString('\n')
		if err != nil {
			fmt.Println("Input error:", err)
			return
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if line == "exit" {
			return
		}

		_, err = conn.Write([]byte(line + "\n"))
		if err != nil {
			fmt.Println("Send error:", err)
			return
		}

		resp, err := server.ReadString('\n')
		if err != nil {
			fmt.Println("Read error:", err)
			return
		}

		fmt.Println(resp)
	}
}
