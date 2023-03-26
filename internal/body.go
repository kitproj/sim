package internal

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/kitproj/sim/internal/types"
)

func getBody(r types.Request) (*bytes.Buffer, error) {
	w := &bytes.Buffer{}
	var err error
	switch value := r.GetBody(); body := value.(type) {
	case nil:
	case string:
		_, err = w.Write([]byte(body))
	case []byte:
		_, err = w.Write(body)
	default:
		err = json.NewEncoder(w).Encode(body)
	}
	return w, err
}

func readBody(r *http.Response) (any, error) {
	cty := r.Header.Get("Content-Type")
	switch cty {
	case "":
		return nil, nil
	case "application/json":
		out := map[string]any{}
		err := json.NewDecoder(r.Body).Decode(&out)
		return out, err
	default:
		out := &bytes.Buffer{}
		_, err := io.Copy(out, r.Body)
		return out.String(), err
	}
}

func writeBody(w io.Writer, value any) error {
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
