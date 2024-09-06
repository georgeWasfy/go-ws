package main

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strings"
)


func calcKey(clientKey string) string {
	clientKey = strings.Trim(clientKey, " ") + "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
	hash := sha1.Sum(([]byte(clientKey)))
	return base64.StdEncoding.EncodeToString(hash[:])
}
func handler(w http.ResponseWriter, req *http.Request) {
	wsKey := getWSKey(req)
	calcKey(wsKey)
	hj, ok := w.(http.Hijacker)
		if !ok {
			http.Error(w, "webserver doesn't support hijacking", http.StatusInternalServerError)
			return
		}
		conn, bufrw, err := hj.Hijack()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Don't forget to close the connection:
		defer conn.Close()
		bufrw.WriteString("Now we're speaking raw TCP. Say hi: ")
		bufrw.Flush()
		s, err := bufrw.ReadString('\n')
		if err != nil {
			log.Printf("error reading string: %v", err)
			return
		}
		fmt.Fprintf(bufrw, "You said: %q\nBye.\n", s)
		bufrw.Flush()

}

func getWSKey(req *http.Request) string {
	return req.Header.Get("Sec-Websocket-Key")
}

func main() {
	log.Print("Server Started")

	http.HandleFunc("/", handler)
	http.ListenAndServe(":5555", nil)
}