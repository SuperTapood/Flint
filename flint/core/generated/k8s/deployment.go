package k8s

import (
	"strings"

	"github.com/SuperTapood/Flint/core/base"
	"github.com/heimdalr/dag"
)

func (deployment *Deployment) GetID() string {
	return deployment.GetName()
}

func (deployment *Deployment) Synth(stack_metadata map[string]any) map[string]any {
	if strings.Contains(deployment.GetName(), "::") {
		panic("invalid name " + deployment.Name)
	}
	namespace := stack_metadata["namespace"].(string)
	obj_map := map[string]any{
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

	template := obj_map["spec"].(map[string]any)["template"].(map[string]any)
	pod_map := deployment.GetPod().Synth(stack_metadata)
	template["metadata"] = pod_map["metadata"]
	template["spec"] = pod_map["spec"]

	return obj_map
}

func (deployment *Deployment) AddToDag(dag *dag.DAG) {
	if strings.Contains(deployment.GetName(), "::") {
		panic("invalid name " + deployment.Name)
	}

	if dag != nil {
		dag.AddVertexByID(deployment.GetID(), deployment.GetID())
	}
}

func (deployment *Deployment) Apply(stack_metadata map[string]any, resources map[string]base.ResourceType, client base.CloudClient) {
	apply_metadata := make(map[string]any)
	apply_metadata["name"] = deployment.GetName()
	apply_metadata["location"] = "/apis/apps/v1/namespaces/" + stack_metadata["namespace"].(string) + "/deployments/"

	client.Apply(apply_metadata, deployment.Synth(stack_metadata))
}

func (deployment *Deployment) Lookup() map[string]any {
	panic("can't lookup a deployment")
}
