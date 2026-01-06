package k8s

import (
	"bytes"
	"fmt"
	"os"
	sync "sync"

	"github.com/SuperTapood/Flint/core/base"
	"github.com/SuperTapood/Flint/core/util"
	"github.com/heimdalr/dag"
)

var outputMu sync.Mutex
var outputBufferMap map[int32]*bytes.Buffer

func (lookup *K8SLookup) GetID() string {
	return lookup.Object.ActualType().GetID()
}

func (lookup *K8SLookup) Synth(stackMetadata map[string]any) map[string]any {
	// fmt.Println("k8slookup is not synthable")
	// fmt.Println("you shouldn't be able to see this error")
	// os.Exit(2)

	return nil
}

func (lookup *K8SLookup) AddToDag(dag *dag.DAG) {}
func (lookup *K8SLookup) Apply(stackMetadata map[string]any, resources map[string]base.ResourceType, client base.CloudClient) error {
	return nil
}
func (lookup *K8SLookup) ExplainFailure(client *util.HttpClient, stackMetadata map[string]any) string {
	fmt.Println("lookup cannot be failed")
	os.Exit(2)
	return ""
}

func (output *K8SOutput) Synth(stackMetadata map[string]any) map[string]any {
	return nil
}

func (output *K8SOutput) Apply(stackMetadata map[string]any, resources map[string]base.ResourceType, client base.CloudClient) error {
	types := output.GetTypes()
	buffer := bytes.Buffer{}

	for _, t := range types {
		if s := t.GetString_(); s != "" {
			fmt.Fprint(&buffer, s)
		}
		if l := t.GetK8Slookup(); l != nil {
			lookup := l
			var lookupId = lookup.GetObject().ActualType().GetID()
			target := resources[lookupId].Synth(stackMetadata)
			kind := target["kind"].(string)
			namespace := target["metadata"].(map[string]any)["namespace"].(string)
			name := target["metadata"].(map[string]any)["name"].(string)
			response, _ := client.GetClient().Get(locationMap[kind]["before_namespace"]+namespace+locationMap[kind]["after_namespace"]+name, nil, true)
			var currentMap = response.Body
			// err := json.Unmarshal(response.Body, &currentMap)
			// if err != nil {
			// 	fmt.Println("could not unmarshal the current secret to a map")
			// 	fmt.Println(err)
			// 	os.Exit(2)
			// }
			var current any = currentMap
			for _, k := range lookup.GetKeys() {
				// must be a map to go deeper
				mmap, ok := current.(map[string]any)
				if !ok {
					fmt.Println("first not ok")
					os.Exit(-1)
				}
				v, ok := mmap[k]
				if !ok {
					fmt.Println("second not ok")
					os.Exit(-1)
				}
				current = v
			}
			fmt.Fprint(&buffer, current)
		}
	}

	outputMu.Lock()
	defer outputMu.Unlock()
	if outputBufferMap == nil {
		outputBufferMap = make(map[int32]*bytes.Buffer)
	}
	outputBufferMap[output.GetIndex()] = &buffer

	os.Exit(-1)
	return nil
}

func (output *K8SOutput) AddToDag(_dag *dag.DAG) {
	_dag.AddVertexByID(output.GetID(), output.GetID())
	for _, lookup := range output.GetTypes() {
		if l := lookup.GetK8Slookup(); l == nil {
			continue
		}
		_dag.AddEdge(output.GetID(), lookup.GetK8Slookup().GetID())
	}
}

func (output *K8SOutput) ExplainFailure(client *util.HttpClient, stackMetadata map[string]any) string {
	fmt.Println("output cannot be failed")
	os.Exit(2)
	return ""
}
