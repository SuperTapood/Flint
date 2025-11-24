package base

import (
	"io"
	"net/http"
)

type CloudClient interface {
	MakeRequest(method string, location string, reader io.Reader) ([]byte, *http.Response)
	Apply(apply_metadata map[string]any, resource map[string]any)
	Delete(delete_metadata map[string]any)
}
