package k8s

import "github.com/heimdalr/dag"

func (pod *Pod) GetID() string {
	return pod.GetName()
}

func (pod *Pod) Synth(dag *dag.DAG) map[string]any {

	obj_map := map[string]any{
		"location":   "/api/v1/namespaces/default/pods",
		"apiVersion": "v1",
		"kind":       "Pod",
		"metadata": map[string]any{
			"name": pod.GetName(),
			"labels": map[string]any {
				"name": pod.GetName(),
			},
			// "namespace": "default",
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

	dag.AddVertexByID(pod.GetID(), pod.GetID())

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
