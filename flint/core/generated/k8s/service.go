package k8s

import (
	"strings"

	"github.com/heimdalr/dag"
)

func (service *Service_) GetID() string {
	return service.GetName()
}

func (service *Service_) GetTargetID() string {
	if pod := service.GetTarget().GetPod(); pod != nil {
		return pod.GetID()
	} else if deployment := service.GetTarget().GetDeployment(); deployment != nil {
		return deployment.GetID()
	}
	panic("got bad service target")
}

func (service *Service_) Synth(stack_name string, namespace string, dag *dag.DAG, objs_map map[string]map[string]any) {
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
				"name": service.GetTargetID(),
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
		err := dag.AddVertexByID(service.GetID(), service.GetID())
		if err != nil {
			panic(err)
		}
		err = dag.AddEdge(service.GetID(), service.GetTargetID())
		if err != nil {
			panic(err)
		}
	}

	objs_map[service.GetID()] = obj_map
}
