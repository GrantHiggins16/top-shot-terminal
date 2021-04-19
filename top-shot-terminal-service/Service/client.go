package service

import (
	"fmt"
	"net"
	"net/http"
)

type Client struct {
	hub  *Hub
	send chan []byte
	addr *net.UDPAddr
	conn *net.UDPConn
}

func (c *Client) writePump() {
	defer func() {
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			_, err := c.conn.WriteTo(message, c.addr)
			if err != nil {
				c.hub.unregister <- c
				return
			}
			// todo: add case for ticker of non response?
		}
	}
}

func runWs(hub *Hub, r *http.Request) {
	addr, err := net.ResolveUDPAddr("udp", r.RemoteAddr)
	if err != nil {
		fmt.Printf("error client")
		return
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Printf("error client")
		return
	}
	client := &Client{hub: hub, addr: addr, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client

	go client.writePump()
}
