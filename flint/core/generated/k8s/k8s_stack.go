package k8s

import (
	"github.com/SuperTapood/Flint/core/base"
	"github.com/heimdalr/dag"
)

func (types *K8STypes) ActualType() base.ResourceType {
	if out := types.GetPod(); out != nil {
		return out
	} else if out := types.GetService(); out != nil {
		return out
	} else if out := types.GetDeployment(); out != nil {
		return out
	} else if out := types.GetSecret(); out != nil {
		return out
	} else if out := types.GetK8Soutput(); out != nil {
		// fmt.Println(out.GetStrings())
		// fmt.Println(out.GetLookups())
		return out
	}
	panic("got bad resource type")
}

func (types *K8STypes) Synth(stack_metadata map[string]any, dag *dag.DAG, obj_map map[string]map[string]any) {
	types.ActualType().Synth(stack_metadata, dag, obj_map)
}

func (stack *K8S_Stack_) Synth(name string) (*dag.DAG, map[string]map[string]any) {
	objs_map := map[string]map[string]any{}
	obj_dag := dag.NewDAG()
	for _, obj := range stack.Objects {
		obj.Synth(stack.GetMetadata(), obj_dag, objs_map)
	}

	return obj_dag, objs_map
}

func (stack *K8S_Stack_) GetMetadata() map[string]any {
	return map[string]any{
		"namespace": stack.GetNamespace(),
	}
}

func (stack *K8S_Stack_) FetchObjects() []base.ResourceType {
	objs := stack.GetObjects()
	out := make([]base.ResourceType, len(objs))

	for _, obj := range objs {
		out = append(out, obj.ActualType())
	}

	return out
}
