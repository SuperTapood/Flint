package k8s

import (
	"fmt"
	"os"
	"strings"

	"github.com/SuperTapood/Flint/core/base"
	"github.com/SuperTapood/Flint/core/util"
	"github.com/heimdalr/dag"
)

func (deployment *Deployment) GetID() string {
	if deployment.GetName() == "" {
		deployment.Name = deployment.GetPod().GetID()
	}
	return deployment.GetName()
}

type STRING string

func (str STRING) enforce() string {
	return string(str)
}

func (deployment *Deployment) Synth(stackMetadata map[string]any) map[string]any {
	if strings.Contains(deployment.GetID(), "::") {
		fmt.Println("invalid name " + STRING(deployment.Name).enforce())
		os.Exit(1)
	}
	namespace := stackMetadata["namespace"].(string)
	podMap := deployment.GetPod().Synth(stackMetadata)
	objMap := map[string]any{
		"apiVersion": "apps/v1",
		"kind":       "Deployment",
		"metadata": map[string]any{
			"name": deployment.GetID(),
			"labels": map[string]any{
				"name": deployment.GetID(),
			},
			"namespace": namespace,
		},
		"spec": map[string]any{
			"replicas": deployment.GetReplicas(),
			"selector": map[string]any{
				"matchLabels": map[string]any{
					"name": deployment.GetPod().GetID(),
				},
			},
			"template": map[string]any{
				"metadata": podMap["metadata"],
				"spec":     podMap["spec"],
			},
		},
	}

	return objMap
}

func (deployment *Deployment) AddToDag(_dag *dag.DAG) {
	if _dag != nil {
		_dag.AddVertexByID(deployment.GetID(), deployment.GetID())
	}
}

func (deployment *Deployment) Apply(stackMetadata map[string]any, resources map[string]base.ResourceType, client base.CloudClient) error {
	applyMetadata := make(map[string]any)
	applyMetadata["name"] = deployment.GetName()
	applyMetadata["location"] = "/apis/apps/v1/namespaces/" + stackMetadata["namespace"].(string) + "/deployments/"

	return client.Apply(applyMetadata, deployment.Synth(stackMetadata))
}

func (deployment *Deployment) ExplainFailure(client *util.HttpClient, stackMetadata map[string]any) string {
	response, _ := client.Get("/apis/apps/v1/namespaces/"+stackMetadata["namespace"].(string)+"/deployments/"+deployment.GetName(), []int{200}, true)
	uid := response.Body["metadata"].(map[string]any)["uid"].(string)
	response, _ = client.Get("/apis/apps/v1/namespaces/"+stackMetadata["namespace"].(string)+"/replicasets/", []int{200}, true)
	replicaUid := ""
	for _, item := range response.Body["items"].([]any) {
		refs := item.(map[string]any)["metadata"].(map[string]any)["ownerReferences"].([]any)
		for _, ref := range refs {
			if ref.(map[string]any)["uid"] == uid {
				replicaUid = item.(map[string]any)["metadata"].(map[string]any)["uid"].(string)
				break
			}
		}
		if replicaUid != "" {
			break
		}
	}

	response, _ = client.Get("/api/v1/namespaces/"+stackMetadata["namespace"].(string)+"/pods/", []int{200}, true)
	for _, item := range response.Body["items"].([]any) {
		refs := item.(map[string]any)["metadata"].(map[string]any)["ownerReferences"].([]any)
		for _, ref := range refs {
			if ref.(map[string]any)["uid"] == replicaUid {
				containerStatuses := item.(map[string]any)["status"].(map[string]any)["containerStatuses"].([]any)
				for _, containerStatus := range containerStatuses {
					return containerStatus.(map[string]any)["state"].(map[string]any)["waiting"].(map[string]any)["message"].(string)
				}
			}
		}
	}

	return "Deployment failed to succeed"
}
