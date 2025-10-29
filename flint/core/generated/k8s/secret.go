package k8s

import (
	"encoding/base64"

	"github.com/heimdalr/dag"
)

func (secret *Secret) GetID() string {
	return secret.GetName()
}

func (secret *Secret) Synth(dag *dag.DAG) map[string]any {
	obj_map := map[string]any{
		"location":   "/api/v1/namespaces/default/secrets",
		"apiVersion": "v1",
		"kind":       "Secret",
		"metadata": map[string]any{
			"name": secret.GetName(),
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

	dag.AddVertexByID(secret.GetID(), secret.GetID())

	return obj_map
}
