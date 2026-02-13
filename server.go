package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"

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
)

// Server holds the server state
type Server struct {
	port     int
	filepath string
	tmpl     *template.Template
}

// NewServer creates a new server instance
func NewServer(port int, filepath string) (*Server, error) {
	tmpl, err := GetTemplate()
	if err != nil {
		return nil, err
	}

	return &Server{
		port:     port,
		filepath: filepath,
		tmpl:     tmpl,
	}, nil
}

// Start starts the HTTP server
func (s *Server) Start() error {
	http.HandleFunc("/", s.handleIndex)
	http.HandleFunc("/ws", s.handleWebSocket)

	addr := fmt.Sprintf(":%d", s.port)
	log.Printf("Serving %s at http://localhost%s", filepath.Base(s.filepath), addr)
	log.Println("Press Ctrl+C to stop")

	return http.ListenAndServe(addr, nil)
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
	clientsMu.Unlock()

	// Keep connection alive and handle disconnection
	defer func() {
		clientsMu.Lock()
		delete(clients, conn)
		clientsMu.Unlock()
		conn.Close()
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
