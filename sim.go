package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/dop251/goja"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/routers"
	"github.com/getkin/kin-openapi/routers/gorillamux"
	"github.com/google/uuid"
	"github.com/kitproj/sim/internal/db"
)

type Sim struct {
	servers map[int]*http.Server
	specs   map[string]*openapi3.T
	vms     map[*openapi3.T]*goja.Runtime
	routers map[*openapi3.T]routers.Router
}

func (s *Sim) add(path string) error {
	spec, err := openapi3.NewLoader().LoadFromFile(path)
	if err != nil {
		return err
	}
	log.Printf("%s: Spec has %d servers", path, len(spec.Servers))
	for _, server := range spec.Servers {
		parse, err := url.Parse(server.URL)
		if err != nil {
			return err
		}
		if parse.Hostname() != "localhost" {
			log.Printf("Skipping server %s", parse.Hostname())
			continue
		}
		port, err := strconv.Atoi(parse.Port())
		if err != nil {
			return err
		}
		for path, item := range spec.Paths {
			for method := range item.Operations() {
				log.Printf("%s %s%s", method, server.URL, path)
			}
		}

		router, err := gorillamux.NewRouter(spec)
		if err != nil {
			return err
		}
		s.routers[spec] = router

		vm := goja.New()
		vm.SetFieldNameMapper(goja.TagFieldNameMapper("json", true))
		var randomUUID = func() string {
			random, err := uuid.NewRandom()
			if err != nil {
				panic(err)
			}
			return random.String()
		}
		if err := vm.Set("randomUUID", randomUUID); err != nil {
			return err
		}
		if err := vm.Set("db", db.Instance); err != nil {
			return err
		}
		script, ok := spec.Extensions["x-sim-script"]
		if ok {
			log.Printf("Found x-sim-script: %v", script)
			if _, err := vm.RunString(script.(string)); err != nil {
				return err
			}
		}
		log.Printf("globals: %v", vm.GlobalObject().Keys())

		s.vms[spec] = vm
		s.specs[path] = spec

		if _, ok := s.servers[port]; !ok {
			server := s.servers[port]
			server = &http.Server{
				Addr:    fmt.Sprintf(":%d", port),
				Handler: http.HandlerFunc(s.Handle),
			}

			go func() {
				log.Printf("Serving on %s\n", server.Addr)
				if err := server.ListenAndServe(); err != nil {
					log.Printf("error serving simulated API: %v", err)
				}
			}()
			s.servers[port] = server
		}
	}
	return nil
}

func (s *Sim) Handle(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	log.Printf("Request: %s %s", r.Method, r.URL.Path)
	log.Printf("Request URL: %v", r.Host)
	spec, route, pathParams, err := s.find(r)
	if err != nil {
		http.Error(w, "Operation not found in servers", http.StatusNotFound)
		return
	}
	op := route.Operation
	log.Printf("Found operation: %v", op.OperationID)
	var writeBody = func(value any) error {
		switch body := value.(type) {
		case nil:
			return nil
		case string:
			_, err := w.Write([]byte(body))
			return err
		case []byte:
			_, err := w.Write(body)
			return err
		default:
			return json.NewEncoder(w).Encode(body)
		}
	}
	script, ok := op.Extensions["x-sim-script"]
	if ok {
		log.Printf("Found x-sim-script: %v", script)
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
		vm := s.vms[spec]
		log.Printf("globals: %v", vm.GlobalObject().Keys())

		if err := vm.Set("request", request); err != nil {
			http.Error(w, fmt.Sprintf("failed to set request: %v", err), http.StatusInternalServerError)
			return
		}
		out, err := vm.RunString(script.(string))
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
		if err := writeBody(response.GetBody()); err != nil {
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
			if err := writeBody(value.Example); err != nil {
				http.Error(w, fmt.Sprintf("failed to encode response: %v", err), http.StatusInternalServerError)
			}
			return
		}
	}

}

func (s *Sim) find(r *http.Request) (*openapi3.T, *routers.Route, map[string]string, error) {
	reqUrl, err := url.Parse("http://" + r.Host)
	if err != nil {
		return nil, nil, nil, err
	}
	for _, spec := range s.specs {
		for _, server := range spec.Servers {
			serverURL, err := url.Parse(server.URL)
			if err != nil {
				return nil, nil, nil, err
			}
			if serverURL.Port() == reqUrl.Port() {
				router := s.routers[spec]
				var route *routers.Route
				route, pathParams, err := router.FindRoute(r)
				if err == nil {
					return spec, route, pathParams, nil
				}
			}
		}
	}
	return nil, nil, nil, fmt.Errorf("no matching server found")
}
