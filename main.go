package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/kitproj/sim/internal"

	"github.com/fsnotify/fsnotify"
)

func init() {
	log.SetOutput(os.Stdout)
}

func main() {
	// Parse command-line flags
	if len(os.Args) <= 1 {
		os.Args = []string{"."}
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Find OpenAPI spec files in directory
	sim := internal.NewSim()

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
		stat, err := os.Stat(path)
		if err != nil {
			log.Fatal(err)
		}
		if stat.IsDir() {
			dir, err := os.ReadDir(path)
			if err != nil {
				log.Fatalf("Error reading directory: %s\n", err)
			}
			for _, file := range dir {
				if filepath.Ext(file.Name()) != ".yaml" {
					continue
				}
				if err := sim.Add(filepath.Join(path, file.Name())); err != nil {
					log.Fatalf("Error adding spec: %s\n", err)
				}
			}
		} else {
			if err := sim.Add(path); err != nil {
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
				if err := sim.Add(event.Name); err != nil {
					log.Printf("Error adding spec: %s\n", err)
				}
			}
		}
	}
}
