package websocket_chat_go

import "github.com/gorilla/websocket"

type Client struct {
	conn   *websocket.Conn
	server *Server
	send   chan []byte
}
