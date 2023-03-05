package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
)

func main() {
	// Parse command-line flags
	specsDir := os.Args[1]

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Find OpenAPI spec files in directory
	specs := map[int][]*openapi3.T{}
	err := filepath.Walk(specsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Base(path) == "openapi.yaml" {
			log.Printf("Loading OpenAPI spec from %s\n", path)
			spec, err := openapi3.NewLoader().LoadFromFile(path)
			if err != nil {
				return err
			}
			log.Printf("Spec has %d servers", len(spec.Servers))
			for _, server := range spec.Servers {
				parse, err := url.Parse(server.URL)
				if err != nil {
					return err
				}
				port, err := strconv.Atoi(parse.Port())
				if err != nil {
					return err
				}
				specs[port] = append(specs[port], spec)

				for path, item := range spec.Paths {
					for method, op := range item.Operations() {
						log.Printf("%s: %s %v%s", op.OperationID, method, parse, path)
					}
				}
			}
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Error reading OpenAPI spec directory: %s\n", err)
	}

	log.Printf("Loaded %d specs", len(specs))

	for port, specs := range specs {

		server := &http.Server{
			Addr: fmt.Sprintf(":%d", port),
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				log.Printf("%s %s", r.Method, r.URL.Path)
				var op *openapi3.Operation
				for _, s := range specs {
					if path, ok := s.Paths[r.URL.Path]; ok {
						op = path.GetOperation(r.Method)
					}
				}
				if op == nil {
					log.Printf("No operation found for %s %s", r.Method, r.URL.Path)
					w.WriteHeader(http.StatusNotFound)
					return
				}
				log.Printf("%s", op.OperationID)
				resp := op.Responses.Get(200)
				if resp == nil {
					log.Printf("No 200 response found for %s %s", r.Method, r.URL.Path)
					w.WriteHeader(http.StatusNotFound)
					return
				}
				example := resp.Value.Content.Get("application/json").Example
				w.WriteHeader(http.StatusOK)
				w.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(w).Encode(example)
			}),
		}
		log.Printf("Serving simulated API on http://localhost%s\n", server.Addr)

		go func() {
			if err := server.ListenAndServe(); err != nil {
				log.Fatalf("error serving simulated API: %v", err)
			}
		}()
	}

	<-ctx.Done()

}
