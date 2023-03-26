package internal

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/kitproj/sim/internal/types"
)

func httpService(r types.Request) map[string]any {
	w, err := getBody(r)
	if err != nil {
		panic(fmt.Errorf("failed to make HTTP request body: %w", err))
	}
	log.Printf("HTTP %s %s", r.GetMethod(), r.GetURL())
	resp, err := http.DefaultClient.Do(&http.Request{
		Method: r.GetMethod(),
		URL:    r.GetURL(),
		Header: httpHeaders(r.GetHeaders()),
		Body:   io.NopCloser(w),
	})
	log.Printf("HTTP %s %s %d", r.GetMethod(), r.GetURL(), resp.StatusCode)
	if err != nil {
		panic(fmt.Errorf("failed to make HTTP request: %w", err))
	}
	body, err := readBody(resp)
	if err != nil {
		panic(fmt.Errorf("failed to read HTTP response body: %w", err))
	}
	return Response{
		"status":  resp.StatusCode,
		"headers": reverseHttpHeaders(resp.Header),
		"body":    body,
	}
}

func httpHeaders(in map[string]string) http.Header {
	out := http.Header{}
	for k, v := range in {
		out.Set(k, v)
	}
	return out
}

func reverseHttpHeaders(in http.Header) map[string]string {
	out := map[string]string{}
	for k, v := range in {
		out[k] = v[0]
	}
	return out
}
