package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	host := flag.String("host", "localhost", "Server host")
	port := flag.String("port", "8080", "Server port")
	flag.Parse()

	// Подключение к серверу
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", *host, *port))
	if err != nil {
		fmt.Println("Connection error:", err)
		return
	}
	defer conn.Close()

	fmt.Printf("Connected to server %s:%s\n", *host, *port)
	fmt.Println("Enter commands (e.g. INSERT users {...})")

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
