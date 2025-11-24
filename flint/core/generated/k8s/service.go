package k8s

import (
	"strings"

	"github.com/SuperTapood/Flint/core/base"
	"github.com/heimdalr/dag"
)

func (service *Service) GetID() string {
	return service.GetName()
}

func (service *Service) GetLabelName() string {
	if pod := service.GetTarget().GetPod(); pod != nil {
		return pod.GetID()
	} else if pod := service.GetTarget().GetDeployment().GetPod(); pod != nil {
		return pod.GetID()
	}
	panic("got bad service target for label")
}

func (service *Service) GetTargetID() string {
	if pod := service.GetTarget().GetPod(); pod != nil {
		return pod.GetID()
	} else if deployment := service.GetTarget().GetDeployment(); deployment != nil {
		return deployment.GetID()
	}
	panic("got bad service target")
}

func (service *Service) GetPrettyName(stack_metadata map[string]any) string {
	return "Kubernetes::Service::" + stack_metadata["namespace"].(string) + "::" + service.GetName()
}

func (service *Service) Synth(stack_metadata map[string]any) map[string]any {
	if strings.Contains(service.GetName(), "::") {
		panic("invalid name " + service.Name)
	}
	namespace := stack_metadata["namespace"].(string)
	obj_map := map[string]any{
		"apiVersion": "v1",
		"kind":       "Service",
		"metadata": map[string]any{
			"name":      service.GetName(),
			"namespace": namespace,
		},
		"spec": map[string]any{
			"type": "NodePort",
			"selector": map[string]any{
				"name": service.GetLabelName(),
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

	return obj_map
}

func (service *Service) AddToDag(dag *dag.DAG) {
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
}

func (service *Service) Apply(stack_metadata map[string]any, resources map[string]base.ResourceType, client base.CloudClient) {
	apply_metadata := make(map[string]any)
	apply_metadata["name"] = service.GetName()
	apply_metadata["location"] = "/api/v1/namespaces/" + stack_metadata["namespace"].(string) + "/services/"

	client.Apply(apply_metadata, service.Synth(stack_metadata))
}

func (service *Service) Lookup() map[string]any {
	panic("fuck")
}
