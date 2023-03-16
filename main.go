package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"

	"github.com/dop251/goja"
	"github.com/fsnotify/fsnotify"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/routers"
)

func main() {
	// Parse command-line flags
	if len(os.Args) <= 1 {
		os.Args = []string{"."}
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Find OpenAPI spec files in directory
	sim := &Sim{
		servers: make(map[int]*http.Server),
		specs:   make(map[string]*openapi3.T),
		vms:     make(map[*openapi3.T]*goja.Runtime),
		routers: make(map[*openapi3.T]routers.Router),
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("Error creating watcher: %s\n", err)
	}
	defer watcher.Close()

	for _, path := range os.Args[1:] {
		log.Printf("Adding path: %s\n", path)
		if err := watcher.Add(path); err != nil {
			log.Fatalf("Error watching directory: %s\n", err)
		}
		dir, err := os.ReadDir(path)
		if err != nil {
			log.Fatalf("Error reading directory: %s\n", err)
		}
		for _, file := range dir {
			if filepath.Ext(file.Name()) != ".yaml" {
				continue
			}
			if err := sim.add(filepath.Join(path, file.Name())); err != nil {
				log.Fatalf("Error adding spec: %s\n", err)
			}
		}
	}
	log.Println("Press Ctrl+C to exit")

	for {
		select {
		case <-ctx.Done():
			log.Println("Exiting...")
			return
		case event := <-watcher.Events:
			log.Println("event:", event)
			if filepath.Ext(event.Name) == ".yaml" && (event.Has(fsnotify.Write) || event.Has(fsnotify.Create)) {
				if err := sim.add(event.Name); err != nil {
					log.Printf("Error adding spec: %s\n", err)
				}
			}
		}
	}
}

var mu = &sync.Mutex{}
