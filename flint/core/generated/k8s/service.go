package k8s

import (
	"fmt"
	"os"
	"strings"

	"github.com/SuperTapood/Flint/core/base"
	"github.com/SuperTapood/Flint/core/util"
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
	fmt.Println("got bad service target")
	os.Exit(2)

	return ""
}

func (service *Service) GetActualTarget() base.ResourceType {
	if pod := service.GetTarget().GetPod(); pod != nil {
		return pod
	} else if deployment := service.GetTarget().GetDeployment(); deployment != nil {
		return deployment
	}
	fmt.Println("got bad service target")
	os.Exit(2)

	return nil
}

func (service *Service) GetTargetID() string {
	return service.GetActualTarget().GetID()
}

func (service *Service) GetPrettyName(stackMetadata map[string]any) string {
	return "Kubernetes::Service::" + stackMetadata["namespace"].(string) + "::" + service.GetName()
}

func (service *Service) Synth(stackMetadata map[string]any) map[string]any {
	if strings.Contains(service.GetName(), "::") {
		fmt.Println("invalid name " + service.Name)
		os.Exit(1)
	}
	namespace := stackMetadata["namespace"].(string)
	objMap := map[string]any{
		"apiVersion": "v1",
		"kind":       "Service",
		"metadata": map[string]any{
			"name":      service.GetName(),
			"namespace": namespace,
		},
		"spec": map[string]any{
			"type": service.GetType(),
			"selector": map[string]any{
				"name": service.GetLabelName(),
			},
			"ports": []any{},
		},
	}

	spec := objMap["spec"].(map[string]any)

	for _, port := range service.GetPorts() {
		portMap := map[string]any{
			"name":       port.Name,
			"protocol":   strings.ToUpper(port.Protocol),
			"port":       port.GetNumber(),
			"targetPort": port.GetNumber(),
		}
		spec["ports"] = append(spec["ports"].([]any), portMap)
	}

	return objMap
}

func (service *Service) AddToDag(_dag *dag.DAG) {
	if _dag != nil {
		_dag.AddVertexByID(service.GetID(), service.GetID())
		err := _dag.AddEdge(service.GetID(), service.GetTargetID())
		if err != nil {
			service.GetActualTarget().AddToDag(_dag)
			err = _dag.AddEdge(service.GetID(), service.GetTargetID())
			if err != nil {
				fmt.Printf("can't add either '%v' (service id) or '%v' (target id) to the DAG\n", service.GetID(), service.GetTargetID())
				fmt.Println(err)
				fmt.Println(_dag)
				os.Exit(2)
			}
		}
	}
}

func (service *Service) Apply(stackMetadata map[string]any, resources map[string]base.ResourceType, client base.CloudClient) error {
	applyMetadata := make(map[string]any)
	applyMetadata["name"] = service.GetName()
	applyMetadata["location"] = "/api/v1/namespaces/" + stackMetadata["namespace"].(string) + "/services/"

	return client.Apply(applyMetadata, service.Synth(stackMetadata))
}

func (service *Service) ExplainFailure(client *util.HttpClient, stackMetadata map[string]any) string {
	// response, _ := client.Get("/api/v1/namespaces/"+stackMetadata["namespace"].(string)+"/services/"+service.GetName(), []int{200}, true)
	// return fmt.Sprintf("%v", response.Body)
	return "Service failed to succeed"
}
