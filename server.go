package main

import (
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"path/filepath"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all connections
		},
	}
	clients   = make(map[*websocket.Conn]bool)
	clientsMu sync.Mutex
	hadClient bool // tracks whether we've ever had a WebSocket client
)

// Server holds the server state
type Server struct {
	port     int
	filepath string
	tmpl     *template.Template
	listener net.Listener
	Done     chan struct{} // closed when all browser clients disconnect
}

// GetPort returns the port the server is configured to use
func (s *Server) GetPort() int {
	return s.port
}

// NewServer creates a new server instance
// If the requested port is not available, it will try up to 10 alternative ports
func NewServer(port int, filepath string) (*Server, error) {
	tmpl, err := GetTemplate()
	if err != nil {
		return nil, err
	}

	// Try to find an available port by binding a listener directly
	// This avoids the race condition of check-then-listen
	listener, availablePort, err := findAvailablePort(port, 10)
	if err != nil {
		return nil, fmt.Errorf("port %d and next 9 ports are all in use: %v", port, err)
	}

	// Warn if we had to use a different port
	if availablePort != port {
		log.Printf("Port %d is in use, using port %d instead", port, availablePort)
	}

	return &Server{
		port:     availablePort,
		filepath: filepath,
		tmpl:     tmpl,
		listener: listener,
		Done:     make(chan struct{}),
	}, nil
}

// Start starts the HTTP server using the pre-bound listener
func (s *Server) Start() error {
	http.HandleFunc("/", s.handleIndex)
	http.HandleFunc("/ws", s.handleWebSocket)

	log.Printf("Serving %s at http://localhost:%d", filepath.Base(s.filepath), s.port)
	log.Println("Press Ctrl+C to stop")

	return http.Serve(s.listener, nil)
}

// handleIndex serves the rendered markdown HTML
func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	// Render the markdown file
	content, err := RenderMarkdown(s.filepath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error rendering markdown: %v", err), http.StatusInternalServerError)
		return
	}

	// Prepare template data
	data := TemplateData{
		Filename: filepath.Base(s.filepath),
		Content:  template.HTML(content),
		Port:     s.port,
	}

	// Execute template
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := s.tmpl.Execute(w, data); err != nil {
		http.Error(w, fmt.Sprintf("Error executing template: %v", err), http.StatusInternalServerError)
	}
}

// handleWebSocket handles WebSocket connections for live reload
func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}

	clientsMu.Lock()
	clients[conn] = true
	hadClient = true
	clientsMu.Unlock()

	// Keep connection alive and handle disconnection
	defer func() {
		clientsMu.Lock()
		delete(clients, conn)
		remaining := len(clients)
		clientsMu.Unlock()
		conn.Close()

		// If no clients remain, wait briefly (for refreshes) then signal shutdown
		if remaining == 0 {
			go func() {
				time.Sleep(2 * time.Second)
				clientsMu.Lock()
				count := len(clients)
				clientsMu.Unlock()
				if count == 0 {
					log.Println("All browser clients disconnected, shutting down...")
					select {
					case <-s.Done:
						// already closed
					default:
						close(s.Done)
					}
				}
			}()
		}
	}()

	// Read messages (ping/pong to keep connection alive)
	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			break
		}
	}
}

// BroadcastReload sends reload message to all connected clients
func BroadcastReload() {
	clientsMu.Lock()
	defer clientsMu.Unlock()

	for client := range clients {
		err := client.WriteMessage(websocket.TextMessage, []byte("reload"))
		if err != nil {
			log.Println("Error broadcasting to client:", err)
			client.Close()
			delete(clients, client)
		}
	}
}

// findAvailablePort tries to bind a listener starting from the given port.
// It returns the held listener and the port it bound to, avoiding race conditions.
func findAvailablePort(startPort int, maxAttempts int) (net.Listener, int, error) {
	for i := 0; i < maxAttempts; i++ {
		port := startPort + i
		addr := fmt.Sprintf(":%d", port)
		listener, err := net.Listen("tcp", addr)
		if err == nil {
			return listener, port, nil
		}
	}
	return nil, 0, fmt.Errorf("no available port found between %d and %d", startPort, startPort+maxAttempts-1)
}
