package k8s

import (
	"github.com/SuperTapood/Flint/core/base"
	"github.com/heimdalr/dag"
)

func (statefulSet *StatefulSet) GetID() string {
	if statefulSet.GetName() == "" {
		statefulSet.Name = statefulSet.GetPod().GetName()
	}
	return statefulSet.GetName()
}

/*
synth this object and add it to both the stack and its directed acyclic graph

Parameters:
  - map[string]any - stack metadata
  - *dag.DAG - the stack's dag to be modified personally by the object
  - map[string] - the object map to be deployed to the cloud provider
*/
func (statefulSet *StatefulSet) AddToDag(_dag *dag.DAG) {
	if _dag != nil {
		_dag.AddVertexByID(statefulSet.GetID(), statefulSet.GetID())
	}
}

func (statefulSet *StatefulSet) Synth(stackMetadata map[string]any) map[string]any {
	namespace := stackMetadata["namespace"].(string)

	vct := make([]map[string]any, 0)
	for _, v := range statefulSet.GetVolumeClaimTemplates() {
		accessModes := make([]string, 0)
		for _, accessMode := range v.GetAccessModes() {
			accessModes = append(accessModes, accessMode.String())
		}
		vct = append(vct, map[string]any{
			"metadata": map[string]any{
				"name": v.GetName(),
			},
			"spec": map[string]any{
				"accessModes":      []string{"ReadWriteOnce"},
				"storageClassName": v.GetStorageClassName(),
				"resources": map[string]any{
					"requests": map[string]any{
						"storage": v.GetStorage(),
					},
				},
			},
		})
	}

	objMap := map[string]any{
		"apiVersion": "apps/v1",
		"kind":       "StatefulSet",
		"metadata": map[string]any{
			"name": statefulSet.GetID(),
			"labels": map[string]any{
				"name": statefulSet.GetID(),
			},
			"namespace": namespace,
		},
		"spec": map[string]any{
			"replicas": statefulSet.GetReplicas(),
			"selector": map[string]any{
				"matchLabels": map[string]any{
					"name": statefulSet.GetPod().GetID(),
				},
			},
			"template":             statefulSet.GetPod().Synth(stackMetadata),
			"volumeClaimTemplates": vct,
		},
	}

	return objMap
}

func (statefulSet *StatefulSet) Apply(stackMetadata map[string]any, resources map[string]base.ResourceType, client base.CloudClient) {
	applyMetadata := make(map[string]any)
	applyMetadata["name"] = statefulSet.GetName()
	applyMetadata["location"] = "/apis/apps/v1/namespaces/" + stackMetadata["namespace"].(string) + "/statefulsets/"

	client.Apply(applyMetadata, statefulSet.Synth(stackMetadata))
}
