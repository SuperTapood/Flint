package k8s

import (
	"fmt"

	"github.com/heimdalr/dag"
)

func (lookup *Lookup) resolve() string {
	obj := lookup.GetObject()
	result := obj.ActualType().Lookup()
	fmt.Println(result)
	return ""
}

func (lookup *Lookup) GetID() string {
	return lookup.Object.ActualType().GetID()
}

func (output *K8SOutput) Synth(stack_metadata map[string]any, dag *dag.DAG, objs_map map[string]map[string]any) {
	lookups := output.GetLookups()
	strings := output.GetStrings()
	dag.AddVertexByID(output.GetID(), output.GetID())
	obj_map := map[string]any{
		"lookups": lookups,
		"strings": strings,
		"action":  "lookup",
		"kind":    "",
		"metadata": map[string]any{
			"namespace": "",
			"name":      "",
		},
		"id": output.GetID(),
	}
	objs_map[output.GetID()] = obj_map
}

func (output *K8SOutput) Lookup() map[string]any {
	panic("can't lookup an output what the fuck are you even trying to do?")
}
