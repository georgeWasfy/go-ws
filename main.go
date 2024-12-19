package main

import (
	"bufio"
	"crypto/sha1"
	"encoding/base64"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

const WEBSOCKETGUID = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"

func getEnvOrDefault(key, defaultValue string) string {
		value := os.Getenv(key)
		if value == "" {
			return defaultValue
		}
		return value
	}
func calculateWebSocketAccept(SecWebSocketKey string) string {
	SecWebSocketKey = strings.Trim(SecWebSocketKey, " ") + WEBSOCKETGUID
	hash := sha1.Sum([]byte(SecWebSocketKey))
	return base64.StdEncoding.EncodeToString(hash[:])
}
func handleWebSocket(conn net.Conn, rwBuffer *bufio.ReadWriter) {
	defer conn.Close()

	for {
		bytesToRead := rwBuffer.Reader.Buffered()
		log.Println(string(bytesToRead))

		if bytesToRead > 0 {
			data, isPrefix, err := rwBuffer.Reader.ReadLine()
			if err != nil {
				log.Println("ssss"+err.Error())
			} else {
				log.Println("cccccc"+string(data), isPrefix)
			}
		}
	}

}
func handleUpgradeRequest(w http.ResponseWriter, req *http.Request) {
	responseKey := calculateWebSocketAccept(req.Header.Get("Sec-Websocket-Key"))
	hj, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "webserver doesn't support hijacking", http.StatusInternalServerError)
		return
	}
	conn, rwBuffer, err := hj.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	// Write the WebSocket handshake response manually
	handshakeResponse := []byte(
		"HTTP/1.1 101 Switching Protocols\r\n" +
			"Upgrade: websocket\r\n" +
			"Connection: Upgrade\r\n" +
			"Sec-WebSocket-Accept: " + responseKey + "\r\n\r\n")
	_, err = conn.Write(handshakeResponse)
	if err != nil {
		log.Printf("Error writing WebSocket handshake response: %v", err)
		return
	}
	go handleWebSocket(conn, rwBuffer)

}

func main() {
	log.Print("Server Started")
	http.HandleFunc("/", handleUpgradeRequest)
	log.Fatal(http.ListenAndServe(":"+getEnvOrDefault("PORT", "5555"), nil))
}
