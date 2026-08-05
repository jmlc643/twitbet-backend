package websocket

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	redisClient *redis.Client
}

func NewHub(rdb *redis.Client) *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		redisClient: rdb,
	}
}

func (h *Hub) Run() {
	go h.listenRedis()
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

func (h *Hub) listenRedis() {
	ctx := context.Background()
	pubsub := h.redisClient.Subscribe(ctx, "market_events")
	defer pubsub.Close()

	ch := pubsub.Channel()
	for msg := range ch {
		h.broadcast <- []byte(msg.Payload)
	}
	log.Println("Redis subscriber closed")
}
