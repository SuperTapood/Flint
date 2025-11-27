package k8s

import (
	"strings"

	"github.com/SuperTapood/Flint/core/base"
	"github.com/heimdalr/dag"
)

func (pod *Pod) GetID() string {
	return pod.GetName()
}

func (pod *Pod) Synth(stackMetadata map[string]any) map[string]any {

	namespace := stackMetadata["namespace"].(string)

	objMap := map[string]any{
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
	spec := objMap["spec"].(map[string]any)
	containers := spec["containers"].([]any)
	container := containers[0].(map[string]any)

	// Add ports dynamically
	for _, p := range pod.GetPorts() {
		port_map := map[string]any{
			"containerPort": p,
		}
		container["ports"] = append(container["ports"].([]any), port_map)
	}

	return objMap
}

func (pod *Pod) AddToDag(_dag *dag.DAG) {
	if strings.Contains(pod.GetName(), "::") {
		panic("invalid name " + pod.Name)
	}
	if _dag != nil {
		_dag.AddVertexByID(pod.GetID(), pod.GetID())
	}
}

func (pod *Pod) Apply(stackMetadata map[string]any, resources map[string]base.ResourceType, client base.CloudClient) {
	applyMetadata := make(map[string]any)
	applyMetadata["name"] = pod.GetName()
	applyMetadata["location"] = "/apis/v1/namespaces/" + stackMetadata["namespace"].(string) + "/pods/"

	client.Apply(applyMetadata, pod.Synth(stackMetadata))
}

func (pod *Pod) Lookup() map[string]any {
	panic("can't lookup a pod")
}
