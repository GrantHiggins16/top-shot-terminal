package service

import (
	"fmt"
	"time"
	"net/http"
	"net/ResolveTCPAddr"
	"net/ListenUDP"
	"net/ResolveUDPAddr"
	"net/DialUDP"
)

type Client struct {
	hub *Hub
	send chan []byte
	addr *UDPAddr
	conn *UDPConn
}



func writePump(c *Client) {
	defer func() {
		c.conn.Close()
	}()
	for {
		select {
			case message, ok := <-c.send:
				_, err = conn.WriteTo(message, c.addr)
				if err != nil {
					// close client - do later
				}
			// todo: add case for ticker of non response?
		}
	}
}

func runWs(hub *Hub, r *http.Request) {
	conn, err := ResolveUDPAddr("udp", r.RemoteAddr)
	if err != nil {
		fmt.Printf(err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client

	go client.writePump()
}
