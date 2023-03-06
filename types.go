package main

import (
	"fmt"
	"strconv"
)

type Request map[string]any

type Response map[string]any

func (r Response) GetStatus() int {
	text, ok := r["status"].(string)
	if ok {
		status, _ := strconv.Atoi(text)
		if status > 0 {
			return status
		}
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
