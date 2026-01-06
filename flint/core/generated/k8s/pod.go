package k8s

import (
	"fmt"
	"os"
	"strings"

	"github.com/SuperTapood/Flint/core/base"
	"github.com/SuperTapood/Flint/core/util"
	"github.com/heimdalr/dag"
)

func (pod *Pod) GetID() string {
	if pod.GetName() == "" {
		pod.Name = pod.GetContainers()[0].GetName()
	}
	return pod.GetName()
}

func (container *Container) Synth(stackMetadata map[string]any) map[string]any {
	objMap := map[string]any{
		"name":  container.GetName(),
		"image": container.GetImage(),
		"ports": []any{}, // start empty
	}

	for _, p := range container.GetPorts() {
		portMap := map[string]any{
			"containerPort": p,
		}
		objMap["ports"] = append(objMap["ports"].([]any), portMap)
	}

	return objMap
}

func (pod *Pod) Synth(stackMetadata map[string]any) map[string]any {
	namespace := stackMetadata["namespace"].(string)

	if pod.GetRestartPolicy() == "" {
		pod.RestartPolicy = "Always"
	}

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
			"containers":    []any{},
			"restartPolicy": pod.GetRestartPolicy(),
		},
	}

	// Navigate to the container map
	spec := objMap["spec"].(map[string]any)

	for _, c := range pod.GetContainers() {
		spec["containers"] = append(spec["containers"].([]any), c.Synth(stackMetadata))
	}

	return objMap
}

func (pod *Pod) AddToDag(_dag *dag.DAG) {
	if strings.Contains(pod.GetName(), "::") {
		fmt.Println("invalid name " + pod.Name)
		os.Exit(1)
	}
	if _dag != nil {
		_dag.AddVertexByID(pod.GetID(), pod.GetID())
	}
}

func (pod *Pod) Apply(stackMetadata map[string]any, resources map[string]base.ResourceType, client base.CloudClient) error {
	applyMetadata := make(map[string]any)
	applyMetadata["name"] = pod.GetName()
	applyMetadata["location"] = "/api/v1/namespaces/" + stackMetadata["namespace"].(string) + "/pods/"

	return client.Apply(applyMetadata, pod.Synth(stackMetadata))
}

func (pod *Pod) ExplainFailure(client *util.HttpClient, stackMetadata map[string]any) string {
	response, _ := client.Get("/api/v1/namespaces/"+stackMetadata["namespace"].(string)+"/pods/"+pod.GetName(), []int{200}, true)
	containerStatuses := response.Body["status"].(map[string]any)["containerStatuses"].([]any)
	for _, containerStatus := range containerStatuses {
		return containerStatus.(map[string]any)["state"].(map[string]any)["waiting"].(map[string]any)["message"].(string)
	}

	return "Pod failed to succeed"
}
