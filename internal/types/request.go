package types

import (
	"fmt"
	"net/url"
)

type Request map[string]any

func (r Request) GetMethod() string {
	if v, ok := r["method"].(string); ok {
		return v
	}
	return "GET"
}

func (r Request) GetURL() *url.URL {
	v, ok := r["url"].(string)
	if !ok {
		panic(fmt.Errorf("url absent or not a string"))
	}
	parsed, err := url.Parse(v)
	if err != nil {
		panic(fmt.Errorf("invalid url %q: %w", v, err))
	}
	return parsed
}

func (r Request) GetHeaders() map[string]string {
	out := map[string]string{}
	headers, ok := r["headers"].(map[string]any)
	if !ok {
		return nil
	}
	for k, v := range headers {
		out[k] = fmt.Sprint(v)
	}
	return out
}

func (r Request) GetBody() any {
	return r["body"]
}
