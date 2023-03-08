package main

import (
	"fmt"
)

type Request map[string]any

type Response map[string]any

func (r Response) GetStatus() int {
	v, ok := r["status"].(int64)
	if ok {
		return int(v)
	}
	if r.GetBody() != nil {
		return 200
	}
	return 204
}

func (r Response) GetHeaders() map[string]string {
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

func (r Response) GetBody() any {
	return r["body"]
}
