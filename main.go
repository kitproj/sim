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

	"github.com/getkin/kin-openapi/routers"

	"github.com/dop251/goja"
	"github.com/getkin/kin-openapi/openapi3"
)

func main() {
	// Parse command-line flags
	specsDir := "."
	if len(os.Args) > 1 {
		specsDir = os.Args[1]
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

	log.Printf("Loading OpenAPI specs from %s\n", specsDir)

	err := filepath.Walk(specsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".yaml" {
			log.Printf("Loading OpenAPI spec from %s\n", path)

			if err = sim.add(path); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Error reading OpenAPI spec directory: %s\n", err)
	}

	log.Println("Press Ctrl+C to exit")

	<-ctx.Done()
}

var mu = &sync.Mutex{}
