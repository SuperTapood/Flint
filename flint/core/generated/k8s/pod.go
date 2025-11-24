package k8s

import (
	"strings"

	"github.com/SuperTapood/Flint/core/base"
	"github.com/heimdalr/dag"
)

func (pod *Pod) GetID() string {
	return pod.GetName()
}

func (pod *Pod) Synth(stack_metadata map[string]any) map[string]any {

	namespace := stack_metadata["namespace"].(string)

	obj_map := map[string]any{
		"location":   "/api/v1/namespaces/" + namespace + "/pods",
		"apiVersion": "v1",
		"kind":       "Pod",
		"metadata": map[string]any{
			"name": pod.GetName(),
			"labels": map[string]any{
				"name": pod.GetName(),
			},
			"namespace": namespace,
		},
		"spec": map[string]any{
			"containers": []any{
				map[string]any{
					"name":  pod.GetName(),
					"image": pod.GetImage(),
					"ports": []any{}, // start empty
				},
			},
		},
	}

	// Navigate to the container map
	spec := obj_map["spec"].(map[string]any)
	containers := spec["containers"].([]any)
	container := containers[0].(map[string]any)

	// Add ports dynamically
	for _, p := range pod.GetPorts() {
		port_map := map[string]any{
			"containerPort": p,
		}
		container["ports"] = append(container["ports"].([]any), port_map)
	}

	return obj_map
}

func (pod *Pod) AddToDag(dag *dag.DAG) {
	if strings.Contains(pod.GetName(), "::") {
		panic("invalid name " + pod.Name)
	}
	if dag != nil {
		dag.AddVertexByID(pod.GetID(), pod.GetID())
	}
}

func (pod *Pod) Apply(stack_metadata map[string]any, resources map[string]base.ResourceType, client base.CloudClient) {
	apply_metadata := make(map[string]any)
	apply_metadata["name"] = pod.GetName()
	apply_metadata["location"] = "/apis/v1/namespaces/" + stack_metadata["namespace"].(string) + "/pods/"

	client.Apply(apply_metadata, pod.Synth(stack_metadata))
}

func (pod *Pod) Lookup() map[string]any {
	panic("can't lookup a pod")
}
