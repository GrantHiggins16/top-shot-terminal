package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/granthiggins16/top-shot-terminal-topshot-terminal-service/service"
)

var addr = flag.String("addr", ":8080", "http service address")

func main() {
	hub := service.NewHub()
	go hub.Run()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		service.RunWs(hub, r)
	})
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
