package k8s

import (
	"github.com/heimdalr/dag"
)

func (deployment *Deployment) GetID() string {
	return deployment.GetName()
}

func (deployment *Deployment) Synth(stack_name string, namespace string, dag *dag.DAG) map[string]any {

	obj_map := map[string]any{
		"location":   "/apis/apps/v1/namespaces/" + namespace + "/deployments",
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

	if dag != nil {
		dag.AddVertexByID(deployment.GetID(), deployment.GetID())
	}

	template := obj_map["spec"].(map[string]any)["template"].(map[string]any)
	pod_map := deployment.GetPod().Synth(stack_name, namespace, nil)
	template["metadata"] = pod_map["metadata"]
	template["spec"] = pod_map["spec"]

	return obj_map
}
