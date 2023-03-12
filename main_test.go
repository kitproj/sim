package main

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSim(t *testing.T) {
	os.RemoveAll(os.Getenv("HOME") + "/.kitproj/sim/db")
	os.Args = []string{"examples"}
	go main()
	time.Sleep(time.Second)
	t.Run("hello", func(t *testing.T) {
		resp, err := http.Get("http://localhost:8080/hello")
		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
	})
	t.Run("script", func(t *testing.T) {
		resp, err := http.Get("http://localhost:4040/teapot")
		assert.NoError(t, err)
		assert.Equal(t, 418, resp.StatusCode)
		assert.Equal(t, "true", resp.Header.Get("Teapot"))
		data, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)
		assert.Equal(t, "{\"message\":\"I'm a teapot\"}\n", string(data))
	})
	t.Run("state", func(t *testing.T) {
		t.Run("listDocuments", func(t *testing.T) {
			resp, err := http.Get("http://localhost:4040/documents")
			assert.NoError(t, err)
			assert.Equal(t, 200, resp.StatusCode)
			data, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)
			assert.Equal(t, "[]\n", string(data))
		})
		var location string
		t.Run("createDocument", func(t *testing.T) {
			resp, err := http.Post("http://localhost:4040/documents", "", bytes.NewBufferString("{\"foo\": \"bar\"}"))
			assert.NoError(t, err)
			assert.Equal(t, 201, resp.StatusCode)
			location = resp.Header.Get("Location")
		})
		t.Run("getDocument", func(t *testing.T) {
			resp, err := http.Get("http://localhost:4040" + location)
			assert.NoError(t, err)
			assert.Equal(t, 200, resp.StatusCode)
			data, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)
			assert.Equal(t, "{\"foo\":\"bar\"}\n", string(data))
		})
		t.Run("deleteDocument", func(t *testing.T) {
			req, err := http.NewRequest(http.MethodDelete, "http://localhost:4040"+location, nil)
			assert.NoError(t, err)
			resp, err := http.DefaultClient.Do(req)
			assert.Equal(t, 204, resp.StatusCode)
		})
	})
}
