package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/granthiggins16/top-shot-terminal-topshot-terminal-service/service"
)

var addr = flag.String("addr", ":8080", "http service address")

func main() {
	http.Handle("/", fs)
	hub := service.NewHub()
	go hub.Run()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Connecting to ws")
		service.RunWs(hub, w, r)
	})
	log.Println("http server started on :8080")
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
