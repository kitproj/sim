package main

import (
	"fmt"
	"strconv"
)

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
	for k, v := range r["headers"].(map[string]any) {
		out[k] = fmt.Sprint(v)
	}
	return out
}

func (r Response) GetBody() any {
	return r["body"]
}
