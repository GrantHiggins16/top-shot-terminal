package service

import (
	"fmt"
	"net"
	"http"
	"time"
)

const (
	// time between requests to get new blocks
	flowUpdateInterval = 10 * time.Second
)

type Hub struct {
	// registered clients
	clients map[*Client]bool

	// register requests from clients
	register chan *Client

	// unregister requests from clients
	unregister chan *Clienti

	// evnets to send to clients
	flowEvents chan []Eventi

	// height of last queried block
	lastBlock int
}

func createHub() *Hub {
	return &Hub{
		clients: make(map[*Client]bool),
		register: make(chan *Client),
		unregister: make(chan *Client),
		flowEvents: make(chan Event)
	}
}

func fetchEvents() {
	
}

func (h *Hub) run() {
	for {
		select {
			case client := <-h.register:
				h.clients[client] = true
			case client := <-h.unregister:
				if _, ok := h.clients[client]; ok {
					delete(h.clients, client)
					close(client.send)
				}
			case message := <-h.broadcast:
				for client := range h.clients {
					select {
						case client.send <- message:
						default:
							close(client.send)
							delete(h.clients, client)
					}
				}
		}
	}
}
