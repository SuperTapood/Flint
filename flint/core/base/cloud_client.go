package base

import (
	"io"
	"net/http"
)

type CloudClient interface {
	MakeRequest(method string, location string, reader io.Reader) ([]byte, *http.Response)
	Apply(ApplyMetadata map[string]any, resource map[string]any)
	Delete(DeleteMetadata map[string]any)
}
