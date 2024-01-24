package main

import "log"

type Server struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func newServer() *Server {
	return &Server{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Server) run() {
	for {
		select {
		case message := <-h.broadcast:
			log.Printf("Broadcast message: %s", message)
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		case client := <-h.register:
			log.Printf("New client connect: %s", client.conn.RemoteAddr())
			h.clients[client] = true
		case client := <-h.unregister:
			log.Printf("Client disconnected: %s", client)
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		}
	}
}
