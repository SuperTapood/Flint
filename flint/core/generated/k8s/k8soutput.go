package k8s

import (
	"bytes"
	"encoding/json"
	"fmt"
	sync "sync"

	"github.com/SuperTapood/Flint/core/base"
	"github.com/heimdalr/dag"
)

var outputMu sync.Mutex
var outputBufferMap map[int32]*bytes.Buffer

func (lookup *K8SLookup) resolve() string {
	obj := lookup.GetObject()
	result := obj.ActualType().Lookup()
	fmt.Println("AAAAAAAAAAA", result)
	return ""
}

func (lookup *K8SLookup) GetID() string {
	return lookup.Object.ActualType().GetID()
}

func (lookup *K8SLookup) Synth(stackMetadata map[string]any) map[string]any {
	panic("WOW")
}

func (lookup *K8SLookup) Lookup() map[string]any {
	panic("can't lookup a lookup what the fuck are you even trying to do?")
}

func (lookup *K8SLookup) AddToDag(dag *dag.DAG) {}
func (lookup *K8SLookup) Apply(stackMetadata map[string]any, resources map[string]base.ResourceType, client base.CloudClient) {
}

func (output *K8SOutput) Synth(stackMetadata map[string]any) map[string]any {
	return nil
}

func (output *K8SOutput) Apply(stackMetadata map[string]any, resources map[string]base.ResourceType, client base.CloudClient) {
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
			body, _ := client.MakeRequest("GET", locationMap[kind]["before_namespace"]+namespace+locationMap[kind]["after_namespace"]+name, bytes.NewReader(make([]byte, 1)))
			var currentMap map[string]any
			err := json.Unmarshal(body, &currentMap)
			if err != nil {
				panic(err)
			}
			var current any = currentMap
			for _, k := range lookup.GetKeys() {
				// must be a map to go deeper
				mmap, ok := current.(map[string]any)
				if !ok {
					panic("badbad")
				}
				v, ok := mmap[k]
				if !ok {
					panic("badbad")
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

func (output *K8SOutput) Lookup() map[string]any {
	panic("can't lookup an output what the fuck are you even trying to do?")
}
