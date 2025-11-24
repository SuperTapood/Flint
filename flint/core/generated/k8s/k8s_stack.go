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
	// } else if out := types.GetLookup(); out != nil {
	// 	return out
	// }
	panic("got bad resource type")
}

func (stack *K8SStack) Synth(name string) (*dag.DAG, map[string]base.ResourceType) {
	objs_map := map[string]base.ResourceType{}
	obj_dag := dag.NewDAG()
	for _, obj := range stack.Objects {
		objs_map[obj.ActualType().GetID()] = obj.ActualType()
		obj.ActualType().AddToDag(obj_dag)
	}

	return obj_dag, objs_map
}

func (stack *K8SStack) GetMetadata() map[string]any {
	return map[string]any{
		"namespace": stack.GetNamespace(),
	}
}

func (stack *K8SStack) FetchObjects() []base.ResourceType {
	objs := stack.GetObjects()
	out := make([]base.ResourceType, len(objs))

	for _, obj := range objs {
		out = append(out, obj.ActualType())
	}

	return out
}
