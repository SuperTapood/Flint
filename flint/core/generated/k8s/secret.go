package k8s

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"github.com/SuperTapood/Flint/core/base"
	"github.com/SuperTapood/Flint/core/util"
	"github.com/heimdalr/dag"
)

func (secret *Secret) GetID() string {
	return secret.GetName()
}

func (secret *Secret) Synth(stackMetadata map[string]any) map[string]any {
	if strings.Contains(secret.GetName(), "::") {
		fmt.Println("invalid name " + secret.Name)
		os.Exit(1)
	}
	namespace := stackMetadata["namespace"].(string)
	data := map[string]string{}

	for _, d := range secret.GetData() {
		data[d.Key] = base64.StdEncoding.EncodeToString([]byte(d.Value))
	}
	objMap := map[string]any{
		"apiVersion": "v1",
		"kind":       "Secret",
		"metadata": map[string]any{
			"name":      secret.GetName(),
			"namespace": namespace,
			"labels": map[string]any{
				"name": secret.GetName(),
			},
		},
		"data": data,
		"type": secret.GetType(),
	}

	return objMap
}

func (secret *Secret) AddToDag(_dag *dag.DAG) {
	if _dag != nil {
		_dag.AddVertexByID(secret.GetID(), secret.GetID())
	}
}

func (secret *Secret) Apply(stackMetadata map[string]any, resources map[string]base.ResourceType, client base.CloudClient) error {
	applyMetadata := make(map[string]any)
	applyMetadata["name"] = secret.GetName()
	applyMetadata["location"] = "/api/v1/namespaces/" + stackMetadata["namespace"].(string) + "/secrets/"

	return client.Apply(applyMetadata, secret.Synth(stackMetadata))
}

func (secret *Secret) ExplainFailure(client *util.HttpClient, stackMetadata map[string]any) string {
	// response, _ := client.Get("/api/v1/namespaces/"+stackMetadata["namespace"].(string)+"/secrets/"+secret.GetName(), []int{200}, true)
	// return fmt.Sprintf("%v", response.Body)
	return "Secret failed to succeed"
}
