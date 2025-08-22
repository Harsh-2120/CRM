package websockets

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow all connections by returning true, adjust this for your CORS policy
		return true
	},
}

// Server represents the WebSocket server
type Server struct {
	clients   map[*websocket.Conn]bool // connected clients
	broadcast chan []byte              // broadcast channel
}

// NewServer initializes a new WebSocket server
func NewServer() *Server {
	return &Server{
		clients:   make(map[*websocket.Conn]bool),
		broadcast: make(chan []byte, 100),
	}
}

// HandleConnections handles incoming WebSocket connections
func (s *Server) HandleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading connection to WebSocket: %v", err)
		return
	}
	defer ws.Close()

	// Register the new client
	s.clients[ws] = true
	log.Printf("Client connected: %v", ws.RemoteAddr())

	// Listen for incoming messages from the client
	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			delete(s.clients, ws)
			break
		}
		// Send the received message to the broadcast channel
		s.broadcast <- message
	}
}

// Start runs the broadcast loop to send messages to all connected clients
func (s *Server) Start() {
	log.Println("WebSocket server started.")
	for {
		log.Println("Waiting for messages on the broadcast channel...")
		// Wait for a message on the broadcast channel
		message := <-s.broadcast
		if len(s.clients) == 0 {
			log.Println("No clients connected, skipping broadcast.")
			continue
		}
		log.Printf("Received message to broadcast: %s", message)
		// Send message to all connected clients
		for client := range s.clients {
			err := client.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Printf("Error writing message to WebSocket: %v", err)
				client.Close()
				delete(s.clients, client)
			}
		}
	}
}

// BroadcastMessage sends a message to all connected clients
func (s *Server) BroadcastMessage(message []byte) {
	log.Printf("Broadcasting message: %s", message)
	select {
	case s.broadcast <- message:
		log.Println("Message sent to broadcast channel.")
	default:
		log.Println("Broadcast channel is full, message dropped.")
	}
}

// Shutdown gracefully closes all WebSocket connections
func (s *Server) Shutdown() {
	for client := range s.clients {
		client.Close()
	}
	close(s.broadcast)
}
