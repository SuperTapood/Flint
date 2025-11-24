package k8s

import (
	"encoding/base64"
	"strings"

	"github.com/SuperTapood/Flint/core/base"
	"github.com/heimdalr/dag"
)

func (secret *Secret) GetID() string {
	return secret.GetName()
}

func (secret *Secret) Synth(stack_metadata map[string]any) map[string]any {
	if strings.Contains(secret.GetName(), "::") {
		panic("invalid name " + secret.Name)
	}
	namespace := stack_metadata["namespace"].(string)
	obj_map := map[string]any{
		"apiVersion": "v1",
		"kind":       "Secret",
		"metadata": map[string]any{
			"name":      secret.GetName(),
			"namespace": namespace,
			"labels": map[string]any{
				"name": secret.GetName(),
			},
		},
		"data": map[string]string{},
		"type": secret.GetType(),
	}

	data := obj_map["data"].(map[string]string)

	for _, d := range secret.GetData() {
		data[d.Key] = base64.StdEncoding.EncodeToString([]byte(d.Value))
	}

	return obj_map
}

func (secret *Secret) AddToDag(dag *dag.DAG) {
	if dag != nil {
		dag.AddVertexByID(secret.GetID(), secret.GetID())
	}
}

func (secret *Secret) Apply(stack_metadata map[string]any, resources map[string]base.ResourceType, client base.CloudClient) {
	apply_metadata := make(map[string]any)
	apply_metadata["name"] = secret.GetName()
	apply_metadata["location"] = "/api/v1/namespaces/" + stack_metadata["namespace"].(string) + "/secrets/"

	client.Apply(apply_metadata, secret.Synth(stack_metadata))
}

func (secret *Secret) Lookup() map[string]any {
	panic("can't lookup a secret")
}
