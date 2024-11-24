package main

import (
	"fmt"
	"golang.org/x/net/websocket"
	"io"
	"net/http"
)

type Server struct {
	conns map[*websocket.Conn]bool
}

func NewServer() *Server {
	return &Server{
		conns: make(map[*websocket.Conn]bool),
	}
}

func (s *Server) handleWS(ws *websocket.Conn) {
	fmt.Println("Websocket incoming connection from...", ws.RemoteAddr().String())

	s.conns[ws] = true
	s.readLoop(ws)
}

func (s *Server) readLoop(ws *websocket.Conn) {
	buf := make([]byte, 1024)

	for {
		n, err := ws.Read(buf)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Client disconnected")
				break
			}
			fmt.Println("Error reading from websocket:", err)
			continue
		}
		msg := buf[:n]

		s.podcast(msg)
	}
}

func (s *Server) podcast(b []byte) {
	for ws := range s.conns {
		go func(ws *websocket.Conn) {
			if _, err := ws.Write(b); err != nil {
				fmt.Println("Error writing to websocket:", err)
			}
		}(ws)
	}
}

func main() {
	serv := NewServer()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Security-Policy", "connect-src 'self' ws://localhost:3001")

		websocket.Handler(serv.handleWS).ServeHTTP(w, r)
	})
	fmt.Println("Listening on :3001...")
	http.ListenAndServe(":3001", nil)
}
