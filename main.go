package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func socketServer(hub *Server, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{server: hub, conn: conn, send: make(chan []byte, 256)}
	client.server.register <- client

	go client.writePump()
	go client.readPump()
}

func main() {
	server := newServer()
	go server.run()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		socketServer(server, w, r)
	})

	port := ":9595"
	httpServer := &http.Server{
		Addr:              port,
		ReadHeaderTimeout: 3 * time.Second,
	}
	log.Printf("Websocket server runing at: %s", port)

	if err := httpServer.ListenAndServe(); err != nil {
		log.Fatal("Socket server run error: ", err)
	}
}
