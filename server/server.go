package main

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
)

func handleClient(conn net.Conn) {
	defer conn.Close()
	buffer := make([]byte, 1024)
	n, _ := conn.Read(buffer)
	requestLine := strings.Split(string(buffer[:n]), "\n")[0]
	fmt.Println("Request:", requestLine)

	parts := strings.Split(requestLine, " ")
	if len(parts) < 2 {
		conn.Write([]byte("HTTP/1.1 400 Bad Request\r\n\r\n"))
		return
	}
	path := parts[1]
	if path == "/" {
		path = "/index.html"
	}

	filePath := "." + path
	data, err := os.ReadFile(filePath)
	if err != nil {
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n404 File Not Found"))
		return
	}

	contentType := getContentType(filePath)
	header := "HTTP/1.1 200 OK\r\nContent-Type: " + contentType + "\r\n\r\n"
	conn.Write([]byte(header))
	conn.Write(data)
}

func getContentType(file string) string {
	ext := filepath.Ext(file)
	switch ext {
	case ".html":
		return "text/html"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".mp3":
		return "audio/mpeg"
	case ".wav":
		return "audio/wav"
	default:
		return "application/octet-stream"
	}
}
