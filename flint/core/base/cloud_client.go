package base

import "github.com/SuperTapood/Flint/core/util"

type CloudClient interface {
	// MakeRequest(method string, location string, reader io.Reader) ([]byte, *http.Response)
	GetClient() *util.HttpClient
	Apply(ApplyMetadata map[string]any, resource map[string]any)
	Delete(DeleteMetadata map[string]any)
}
