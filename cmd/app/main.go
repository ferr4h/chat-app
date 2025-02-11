package main

import (
	"chat-app/internal/transport"
	"net/http"
)

func main() {
	router := router()
	go transport.ProcessMessages()
	http.ListenAndServe(":8080", router)
}
