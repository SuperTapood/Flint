package k8s

import (
	"fmt"
	"os"
	"strings"

	"github.com/SuperTapood/Flint/core/base"
	"github.com/heimdalr/dag"
)

func (deployment *Deployment) GetID() string {
	return deployment.GetName()
}

type STRING string

func (str STRING) enforce() string {
	return string(str)
}

func (deployment *Deployment) Synth(stackMetadata map[string]any) map[string]any {
	if strings.Contains(deployment.GetName(), "::") {
		fmt.Println("invalid name " + STRING(deployment.Name).enforce())
		os.Exit(1)
	}
	namespace := stackMetadata["namespace"].(string)
	objMap := map[string]any{
		"apiVersion": "apps/v1",
		"kind":       "Deployment",
		"metadata": map[string]any{
			"name": deployment.GetName(),
			"labels": map[string]any{
				"name": deployment.GetName(),
			},
			"namespace": namespace,
		},
		"spec": map[string]any{
			"replicas": deployment.GetReplicas(),
			"selector": map[string]any{
				"matchLabels": map[string]any{
					"name": deployment.GetPod().GetName(),
				},
			},
			"template": map[string]any{
				"metadata": map[string]any{},
				"spec":     map[string]any{},
			},
		},
	}

	template := objMap["spec"].(map[string]any)["template"].(map[string]any)
	podMap := deployment.GetPod().Synth(stackMetadata)
	template["metadata"] = podMap["metadata"]
	template["spec"] = podMap["spec"]

	return objMap
}

func (deployment *Deployment) AddToDag(_dag *dag.DAG) {
	if _dag != nil {
		_dag.AddVertexByID(deployment.GetID(), deployment.GetID())
	}
}

func (deployment *Deployment) Apply(stackMetadata map[string]any, resources map[string]base.ResourceType, client base.CloudClient) {
	applyMetadata := make(map[string]any)
	applyMetadata["name"] = deployment.GetName()
	applyMetadata["location"] = "/apis/apps/v1/namespaces/" + stackMetadata["namespace"].(string) + "/deployments/"

	client.Apply(applyMetadata, deployment.Synth(stackMetadata))
}
