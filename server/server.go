package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os/exec"
	"strings"
)

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Count   int         `json:"count,omitempty"`
}

func writeJSON(conn net.Conn, resp Response) {
	b, _ := json.Marshal(resp)
	conn.Write(append(b, '\n'))
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, " ", 4)
		if len(parts) != 4 {
			writeJSON(conn, Response{
				Status:  "error",
				Message: "invalid command format",
			})
			continue
		}

		database, collection, action, jsonArg :=
			parts[0], parts[1], parts[2], parts[3]

		if len(jsonArg) > 1 {
			if (jsonArg[0] == '"' && jsonArg[len(jsonArg)-1] == '"') ||
				(jsonArg[0] == '\'' && jsonArg[len(jsonArg)-1] == '\'') {
				jsonArg = jsonArg[1 : len(jsonArg)-1]
			}
		}

		cmd := exec.Command(
			"../database/main",
			database,
			collection,
			action,
			jsonArg,
		)

		output, err := cmd.CombinedOutput()
		if err != nil {
			writeJSON(conn, Response{
				Status:  "error",
				Message: err.Error(),
			})
			continue
		}

		outputStr := strings.TrimSpace(string(output))

		switch action {

		case "find":
			var data []interface{}
			_ = json.Unmarshal([]byte(outputStr), &data)

			writeJSON(conn, Response{
				Status:  "success",
				Message: fmt.Sprintf("Fetched %d docs from %s", len(data), collection),
				Data:    data,
				Count:  len(data),
			})

		case "insert":
			writeJSON(conn, Response{
				Status:  "success",
				Message: fmt.Sprintf("Document inserted into %s", collection),
			})

		case "delete":
			writeJSON(conn, Response{
				Status:  "success",
				Message: fmt.Sprintf("Documents deleted from %s", collection),
			})

		default:
			writeJSON(conn, Response{
				Status:  "error",
				Message: "unknown action",
			})
		}
	}
}

func main() {
	ln, _ := net.Listen("tcp", ":8080")
	fmt.Println("Server started on :8080")

	for {
		conn, _ := ln.Accept()
		go handleConnection(conn)
	}
}