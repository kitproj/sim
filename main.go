package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"

	"github.com/dop251/goja"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/routers"
	"github.com/getkin/kin-openapi/routers/gorillamux"
	"github.com/google/uuid"
	"github.com/kitproj/sim/internal/db"
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
	servers := map[int][]*openapi3.T{}

	log.Printf("Loading OpenAPI specs from %s\n", specsDir)

	err := filepath.Walk(specsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".yaml" {
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
				servers[port] = append(servers[port], spec)

				for path, item := range spec.Paths {
					for method := range item.Operations() {
						log.Printf("%s %s%s", method, server.URL, path)
					}
				}
			}
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Error reading OpenAPI spec directory: %s\n", err)
	}

	for port, specs := range servers {
		startServer(port, specs)
	}

	log.Println("Press Ctrl+C to exit")

	<-ctx.Done()
}

func startServer(port int, specs []*openapi3.T) {
	server := &http.Server{
		Addr: fmt.Sprintf(":%d", port),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("Request: %s %s", r.Method, r.URL.Path)
			var op *openapi3.Operation
			var pathParams map[string]string
			for _, s := range specs {
				router, err := gorillamux.NewRouter(s)
				if err != nil {
					http.Error(w, fmt.Sprintf("failed to create router: %v", err), http.StatusInternalServerError)
					return
				}
				var route *routers.Route
				route, pathParams, err = router.FindRoute(r)
				if err == nil {
					log.Printf("Found route: %v", route)
					op = route.Operation
					break
				}
			}
			if op == nil {
				http.Error(w, "Operation not found in servers", http.StatusNotFound)
				return
			}
			expr, ok := op.Extensions["x-sim-script"]
			if ok {
				log.Printf("Found x-sim-script: %v", expr)

				queryParams := map[string]string{}
				for key := range r.URL.Query() {
					queryParams[key] = r.URL.Query().Get(key)
				}
				headers := map[string]string{}
				for key := range r.Header {
					headers[key] = r.Header.Get(key)
				}
				body := map[string]any{}
				_ = json.NewDecoder(r.Body).Decode(&body)
				request := map[string]any{
					"method":      r.Method,
					"path":        r.URL.Path,
					"queryParams": queryParams,
					"pathParams":  pathParams,
					"headers":     headers,
					"body":        body,
				}
				vm := goja.New()
				vm.SetFieldNameMapper(goja.TagFieldNameMapper("json", true))
				if err := vm.Set("request", request); err != nil {
					http.Error(w, fmt.Sprintf("failed to set request: %v", err), http.StatusInternalServerError)
					return
				}
				var randomUUID = func() string {
					random, err := uuid.NewRandom()
					if err != nil {
						panic(err)
					}
					return random.String()
				}
				if err := vm.Set("randomUUID", randomUUID); err != nil {
					http.Error(w, fmt.Sprintf("failed to set randomUUID: %v", err), http.StatusInternalServerError)
					return
				}
				if err := vm.Set("db", db.Instance); err != nil {
					http.Error(w, fmt.Sprintf("failed to set db: %v", err), http.StatusInternalServerError)
					return
				}
				out, err := vm.RunString(expr.(string))
				if err != nil {
					http.Error(w, fmt.Sprintf("failed to run expression: %v", err), http.StatusInternalServerError)
					return
				}
				response := Response(out.Export().(map[string]any))
				if response == nil {
					http.Error(w, fmt.Sprintf("failed to export response: %T", out.Export()), http.StatusInternalServerError)
					return
				}
				log.Printf("Response: %v", response)
				for key, value := range response.GetHeaders() {
					w.Header().Set(key, value)
				}
				if _, ok := w.Header()["Content-Type"]; !ok {
					w.Header().Set("Content-Type", "application/json")
				}
				w.WriteHeader(response.GetStatus())

				switch body := response.GetBody().(type) {
				case nil:
				case string:
					_, err = w.Write([]byte(body))
				case []byte:
					_, err = w.Write(body)
				default:
					err = json.NewEncoder(w).Encode(response.GetBody())
				}
				if err != nil {
					http.Error(w, fmt.Sprintf("failed to encode response: %v", err), http.StatusInternalServerError)
					return
				}
				return
			}
			// TODO should be ordered
			for status, resp := range op.Responses {
				status, _ := strconv.Atoi(status)
				for mediaType, value := range resp.Value.Content {
					w.Header().Set("Content-Type", mediaType)
					w.WriteHeader(status)
					_ = json.NewEncoder(w).Encode(value.Example)
					return
				}
			}
		}),
	}

	go func() {
		log.Printf("Serving on %s\n", server.Addr)
		if err := server.ListenAndServe(); err != nil {
			log.Printf("error serving simulated API: %v", err)
		}
	}()
}
