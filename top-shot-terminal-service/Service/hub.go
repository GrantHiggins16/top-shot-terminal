package service

import (
	"github.com/onflow/flow-go-sdk/client"
	"context"
	"fmt"
	"time"
	"google.golang.org/grpc"
)

const (
	// time between requests to get new blocks
	flowUpdateInterval = 10 * time.Second
	
	// moment listed event id
	listedEventId = "A.c1e4f4f4c4257510.Market.MomentListed"
)

type Hub struct {
	// registered clients
	clients map[*Client]bool

	// register requests from clients
	register chan *Client

	// unregister requests from clients
	unregister chan *Client

	// evnets to send to clients
	flowEvents chan *Event

	// height of last queried block
	lastBlock uint64
}

func newHub() *Hub {
	return &Hub{
		clients: make(map[*Client]bool),
		register: make(chan *Client),
		unregister: make(chan *Client),
		flowEvents: make(chan *Event),
	}
}

func (h *Hub) fetchEvents(c *client.Client) {
	// get the latest sealed block
	latestBlock, err := c.GetLatestBlock(context.Background(), true)
	if err != nil {
		fmt.Errorf("Unable to fetch latest block")
	}

	endBlockHeight := latestBlock.Height
	startBlockHeight := h.lastBlock
	if startBlockHeight == 0 {
		startBlockHeight = latestBlock.Height
	}		
	blockEvents, err := c.GetEventsForHeightRange(context.Background(), client.EventRangeQuery{
   		Type:        listedEventId,
		StartHeight: startBlockHeight,
		EndHeight:   endBlockHeight,
	}) 
	if err != nil {
		fmt.Printf("Unable to fetch events from latest blocks")
	}
	var newEvents []Event
	for i, eventResponse := range blockEvents {		
		event := NewEvent(eventResponse.Fields[0], eventResponse.Fields[2], endBlockHeight, c)
		h.flowEvents <- event 
	}
	h.lastBlock = endBlockHeight 
}

func (h *Hub) run() {
	fetchEventsTicker := time.NewTicker(10 * time.Second)
	client, err := client.New("access.mainnet.nodes.onflow.org:9000", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
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
				go h.fetchEvents(client)
		}
	}
}
