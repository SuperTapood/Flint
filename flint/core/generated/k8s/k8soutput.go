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
var buffer bytes.Buffer

func (lookup *K8SLookup) resolve() string {
	obj := lookup.GetObject()
	result := obj.ActualType().Lookup()
	fmt.Println("AAAAAAAAAAA", result)
	return ""
}

func (lookup *K8SLookup) GetID() string {
	return lookup.Object.ActualType().GetID()
}

func (lookup *K8SLookup) Synth(stack_metadata map[string]any, dag *dag.DAG, objs_map map[string]map[string]any) {
	panic("WOW")
}

func (lookup *K8SLookup) Lookup() map[string]any {
	panic("can't lookup a lookup what the fuck are you even trying to do?")
}

func (output *K8SOutput) Synth(stack_metadata map[string]any) map[string]any {
	// lookups := output.GetLookups()
	// strings := output.GetStrings()

	// obj_map := map[string]any{
	// 	"lookups": lookups,
	// 	"strings": strings,
	// 	"action":  "lookup",
	// 	"kind":    "",
	// 	"metadata": map[string]any{
	// 		"namespace": "",
	// 		"name":      "",
	// 	},
	// 	"id": output.GetID(),
	// }
	return nil
}

func (output *K8SOutput) Apply(stack_metadata map[string]any, resources map[string]base.ResourceType, client base.CloudClient) {
	outputMu.Lock()
	defer outputMu.Unlock()
	lookups := output.GetLookups()
	strings := output.GetStrings()
	length := max(len(lookups), len(strings))
	for i := range length {
		if i < len(strings) {
			fmt.Fprint(&buffer, strings[i])
		}
		if i < len(lookups) {
			lookup := lookups[i]
			var lookup_id = lookup.GetObject().ActualType().GetID()
			target := resources[lookup_id].Synth(stack_metadata)
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

	fmt.Fprintln(&buffer)
}

func (output *K8SOutput) AddToDag(dag *dag.DAG) {
	dag.AddVertexByID(output.GetID(), output.GetID())
	for _, lookup := range output.GetLookups() {
		dag.AddEdge(output.GetID(), lookup.GetObject().ActualType().GetID())
	}
}

func (output *K8SOutput) Lookup() map[string]any {
	panic("can't lookup an output what the fuck are you even trying to do?")
}
