package main

import (
	"log"
	"time"

	"github.com/fsnotify/fsnotify"
)

// WatchFile monitors a file for changes and calls onChange when it's modified
func WatchFile(filepath string, onChange func()) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	err = watcher.Add(filepath)
	if err != nil {
		watcher.Close()
		return err
	}

	// Debouncing variables
	var timer *time.Timer
	const debounceDelay = 300 * time.Millisecond

	go func() {
		defer watcher.Close()

		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				// Only react to write events
				if event.Op&fsnotify.Write == fsnotify.Write {
					// Reset the timer for debouncing
					if timer != nil {
						timer.Stop()
					}
					timer = time.AfterFunc(debounceDelay, onChange)
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("Watcher error:", err)
			}
		}
	}()

	return nil
}
