package k8s

import (
	"encoding/base64"
	"strings"

	"github.com/heimdalr/dag"
)

func (secret *Secret) GetID() string {
	return secret.GetName()
}

func (secret *Secret) Synth(stack_name string, namespace string, dag *dag.DAG, objs_map map[string]map[string]any) {
	if strings.Contains(secret.GetName(), "::") {
		panic("invalid name " + secret.Name)
	}
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

	if dag != nil {
		dag.AddVertexByID(secret.GetID(), secret.GetID())
	}

	objs_map[secret.GetID()] = obj_map
}

func (secret *Secret) Lookup() map[string]any {
	panic("can't lookup a secret")
}
