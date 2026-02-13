package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pkg/browser"
)

const version = "0.1.0"

func main() {
	// Define CLI flags
	port := flag.Int("port", 3000, "Port to serve on")
	noBrowser := flag.Bool("no-browser", false, "Don't open browser automatically")
	showVersion := flag.Bool("version", false, "Show version")
	flag.BoolVar(showVersion, "v", false, "Show version (shorthand)")
	helpFlag := flag.Bool("help", false, "Show help")
	flag.BoolVar(helpFlag, "h", false, "Show help (shorthand)")

	flag.Parse()

	// Show version
	if *showVersion {
		fmt.Printf("mdview v%s\n", version)
		os.Exit(0)
	}

	// Show help
	if *helpFlag {
		fmt.Println("mdview - Local markdown file viewer with live reload")
		fmt.Println("\nUsage:")
		fmt.Println("  mdview [options] <file.md>")
		fmt.Println("\nOptions:")
		flag.PrintDefaults()
		fmt.Println("\nExamples:")
		fmt.Println("  mdview README.md")
		fmt.Println("  mdview --port 8080 docs/guide.md")
		fmt.Println("  mdview --no-browser file.md")
		os.Exit(0)
	}

	// Check if file argument is provided
	if flag.NArg() < 1 {
		fmt.Println("Error: No file specified")
		fmt.Println("Usage: mdview [options] <file.md>")
		fmt.Println("Try 'mdview --help' for more information")
		os.Exit(2)
	}

	filepath := flag.Arg(0)

	// Validate file exists and is readable
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		fmt.Printf("Error: File not found: %s\n", filepath)
		fmt.Println("Make sure the file exists and you have read permissions.")
		os.Exit(1)
	}

	// Try to read the file to ensure it's accessible
	if _, err := os.ReadFile(filepath); err != nil {
		fmt.Printf("Error: Cannot read file: %s\n", filepath)
		fmt.Printf("Permission denied or file is not accessible.\n")
		os.Exit(1)
	}

	// Create server
	server, err := NewServer(*port, filepath)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// Start file watcher
	err = WatchFile(filepath, func() {
		log.Println("File changed, reloading...")
		BroadcastReload()
	})
	if err != nil {
		log.Fatalf("Failed to start file watcher: %v", err)
	}

	// Start server in goroutine
	go func() {
		if err := server.Start(); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Open browser after a short delay to ensure server is ready
	if !*noBrowser {
		time.Sleep(500 * time.Millisecond)
		url := fmt.Sprintf("http://localhost:%d", *port)
		if err := browser.OpenURL(url); err != nil {
			log.Printf("Failed to open browser: %v", err)
			log.Printf("Please open %s manually", url)
		}
	}

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("\nShutting down...")
}
