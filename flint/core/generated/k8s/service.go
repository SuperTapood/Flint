package k8s

import (
	"strings"

	"github.com/heimdalr/dag"
)

func (service *Service) GetID() string {
	return service.GetName()
}

func (service *Service) Synth(stack_name string, namespace string, dag *dag.DAG) map[string]any {

	obj_map := map[string]any{
		"location":   "/api/v1/namespaces/" + namespace + "/services",
		"apiVersion": "v1",
		"kind":       "Service",
		"metadata": map[string]any{
			"name":      service.GetName(),
			"namespace": namespace,
		},
		"spec": map[string]any{
			"type": "NodePort",
			"selector": map[string]any{
				"name": service.GetTarget().GetID(),
			},
			"ports": []any{},
		},
	}

	spec := obj_map["spec"].(map[string]any)

	for _, port := range service.GetPorts() {
		port_map := map[string]any{
			"name":       port.Name,
			"protocol":   strings.ToUpper(port.Protocol),
			"port":       port.GetNumber(),
			"targetPort": port.GetNumber(),
		}
		spec["ports"] = append(spec["ports"].([]any), port_map)
	}

	if dag != nil {
		dag.AddVertexByID(service.GetID(), service.GetID())
		dag.AddEdge(service.GetTarget().GetID(), service.GetID())
	}

	return obj_map
}
