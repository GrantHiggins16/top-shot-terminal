package service

import (
	"github.com/onflow/flow-go-sdk/client"
	"context"
	"fmt"
	"net"
	"http"
	"time"
)

const (
	// time between requests to get new blocks
	flowUpdateInterval := 10 * time.Second
	
	// moment listed event id
	listedEventId := "A.c1e4f4f4c4257510.Market.MomentListed"
)

type Hub struct {
	// registered clients
	clients map[*Client]bool

	// register requests from clients
	register chan *Client

	// unregister requests from clients
	unregister chan *Client

	// evnets to send to clients
	flowEvents chan []Event

	// height of last queried block
	lastBlock int
}

func newHub() *Hub {
	return &Hub{
		clients: make(map[*Client]bool),
		register: make(chan *Client),
		unregister: make(chan *Client),
		flowEvents: make(chan Event)
	}
}

func fetchEvents(h *Hub) {
	// get the latest sealed block
	latestBlock, err := client.GetLatestBlock(context.Background(), true)
	if err != nil {
		fmt.Errorf("Unable to fetch latest block %w", err)
	}

	endBlockHeight := latestBlock.block.height
	startBlockHeight := lastBlock.block.height
	if startBlockHeight == nil {
		startBlockHeight = latestBlock.block.height
	}		
	blockEvents, err := flowClient.GetEventsForHeightRange(context.Background(), client.EventRangeQuery{
   		Type:        "A.c1e4f4f4c4257510.Market.MomentPurchased",
		StartHeight: startBlockHeight,
		EndHeight:   endBlockHeight,
	}) 
	if err != nil {
		fmt.Printf(err)
	}
	var newEvents []Event
	for i, eventResponse := range events {		
		event := NewEvent(eventResponse.Fields[0], eventResponse.Fields[2])
		h.flowEvents = append(h.flowEvents, event) 
	}
}

func (h *Hub) run() {
	fetchEventsTicker := time.NewTicker(10 * time.Second)
	for {
		select {
			case client := <-h.register:
				h.clients[client] = true
			case client := <-h.unregister:
				if _, ok := h.clients[client]; ok {
					delete(h.clients, client)
					close(client.send)
				}
			case message := <-h.flowEvents:
				for client := range h.clients {
					select {
						case client.send <- message:
						default:
							close(client.send)
							delete(h.clients, client)
					}
				}
			case <- fetchEventsTicker.C:
				go fetchEvents(h)
		}
	}
}
