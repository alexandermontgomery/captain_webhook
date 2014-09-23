package main

import (
	"fmt"
	"net/http"
	"slack_webhooks/webhooks"
	"bytes"
)

func index(w http.ResponseWriter, req *http.Request){
	fmt.Fprintf(w, "Hello World!!")
}

func webhookHandle(w http.ResponseWriter, req *http.Request){
	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)

	msg := webhooks.NewMessage("foo",buf.Bytes())
    go webhooks.Publish(msg)
}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/webhook", webhookHandle)
	addr := ":8080"
	fmt.Printf("Listening on %s\n", addr)
	http.ListenAndServe(addr, nil)	
}