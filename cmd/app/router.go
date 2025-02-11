package main

import (
	"chat-app/internal/transport"
	"github.com/bmizerany/pat"
	"net/http"
)

func router() http.Handler {
	mux := pat.New()
	mux.Get("/", http.HandlerFunc(transport.DisplayIndexPage))
	mux.Get("/ws", http.HandlerFunc(transport.Connect))
	return mux
}
