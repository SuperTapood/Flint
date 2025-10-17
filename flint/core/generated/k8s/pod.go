package k8s

import (
	"github.com/google/uuid"
)

func (pod *Pod) Synth() (uuid.UUID, map[string]any) {
	var uuid = uuid.New()

	obj_map := map[string]any{
		"apiVersion": "v1",
		"kind":       "Pod",
		"metadata": map[string]any{
			"name":      pod.GetName(),
			"namespace": "default",
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

	return uuid, obj_map
}
