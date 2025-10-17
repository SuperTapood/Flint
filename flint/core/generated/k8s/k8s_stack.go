package k8s

import (
	"github.com/SuperTapood/Flint/core/base"
	"github.com/google/uuid"
	"github.com/heimdalr/dag"
)

func (types *K8STypes) ActualType() base.ResourceType {
	if out := types.GetPod(); out != nil {
		return out
	}
	panic("got bad type")
}

func (types *K8STypes) Synth() (uuid.UUID, map[string]any) {
	return types.ActualType().Synth()
}

func (stack *K8S_Stack_) Synth() (*dag.DAG, map[uuid.UUID]map[string]any) {
	objs_map := map[uuid.UUID]map[string]any{}
	var obj_dag = dag.NewDAG()
	for _, obj := range stack.Objects {
		var uuid, obj_map = obj.Synth()
		// log.Printf("%d: %s", i, uuid)
		objs_map[uuid] = obj_map
		obj_dag.AddVertex(uuid)
	}

	return obj_dag, objs_map
}
