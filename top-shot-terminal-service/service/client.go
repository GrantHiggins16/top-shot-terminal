package service

import (
	"log"
	"net"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,

	// For now, do no checking and allow any connection
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Client struct {
	hub  *Hub
	send chan []byte
	conn *websocket.Conn
}

func (c *Client) writePump() {
	defer c.conn.Close()
	for {
		select {
		case message := <-c.send:
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				c.hub.unregister <- c
				return
			}
			w.Write(message)
			// todo: add case for ticker of non response?
		}
	}
}

func RunWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		conn.Close()
		log.Fatal("Failed to upgrade to websocket connection")
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client

	go client.writePump()
}
