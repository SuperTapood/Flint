package k8s

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/SuperTapood/Flint/core/base"
	"github.com/SuperTapood/Flint/core/util"
	"github.com/heimdalr/dag"
)

func (daemonSet *DaemonSet) GetID() string {
	if daemonSet.GetName() == "" {
		daemonSet.Name = daemonSet.GetPod().GetID()
	}
	return daemonSet.GetName()
}

func (daemonSet *DaemonSet) Synth(stackMetadata map[string]any) map[string]any {
	if strings.Contains(daemonSet.GetID(), "::") {
		fmt.Println("invalid name " + daemonSet.Name)
		os.Exit(1)
	}
	namespace := stackMetadata["namespace"].(string)
	podMap := daemonSet.GetPod().Synth(stackMetadata)
	objMap := map[string]any{
		"apiVersion": "apps/v1",
		"kind":       "DaemonSet",
		"metadata": map[string]any{
			"name": daemonSet.GetID(),
			"labels": map[string]any{
				"name": daemonSet.GetID(),
			},
			"namespace": namespace,
		},
		"spec": map[string]any{
			"replicas": daemonSet.GetReplicas(),
			"selector": map[string]any{
				"matchLabels": map[string]any{
					"name": daemonSet.GetPod().GetID(),
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

func (daemonSet *DaemonSet) AddToDag(_dag *dag.DAG) {
	if _dag != nil {
		_dag.AddVertexByID(daemonSet.GetID(), daemonSet.GetID())
	}
}

func (daemonSet *DaemonSet) Apply(stackMetadata map[string]any, resources map[string]base.ResourceType, client base.CloudClient) error {
	applyMetadata := make(map[string]any)
	applyMetadata["name"] = daemonSet.GetName()
	applyMetadata["location"] = "/apis/apps/v1/namespaces/" + stackMetadata["namespace"].(string) + "/daemonsets/"

	return client.Apply(applyMetadata, daemonSet.Synth(stackMetadata), daemonSet, stackMetadata)
}

func (daemonSet *DaemonSet) Get(client *util.HttpClient, stackMetadata map[string]any, acceptedStatusCodes []int, autohandleErrors bool) (*util.HttpResponse, error) {
	return client.Get("/apis/apps/v1/namespaces/"+stackMetadata["namespace"].(string)+"/daemonsets/"+daemonSet.GetName(), acceptedStatusCodes, autohandleErrors)
}

func (daemonSet *DaemonSet) ExplainFailure(client *util.HttpClient, stackMetadata map[string]any) string {
	response, _ := daemonSet.Get(client, stackMetadata, []int{200}, true)
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
					cs := containerStatus.(map[string]any)
					state := cs["state"].(map[string]any)
					waiting := state["waiting"].(map[string]any)
					reason := waiting["reason"].(string)
					if reason == "ContainerCreating" {
						// too early
						time.Sleep(50 * time.Millisecond)
						return daemonSet.ExplainFailure(client, stackMetadata)
					}
					return waiting["message"].(string)
				}
			}
		}
	}

	return "Deployment failed to succeed"
}
