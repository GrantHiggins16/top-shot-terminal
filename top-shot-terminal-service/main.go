package main

import (
	"github.com/onflow/flow-go-sdk/client"
	"context"
	"fmt"
	"net/http"
	"time"
	"github.com/granthiggins16/top-shot-terminal-topshot-terminal-service/Service"
	"github.com/granthiggins16/top-shot-terminal-topshot-terminal-service/Service"
)

func main() {
	hub := newHub()
	go hub.run()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		client.runWs(hub, w, r)
	})
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
