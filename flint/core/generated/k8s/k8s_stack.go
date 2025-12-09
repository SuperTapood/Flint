package k8s

import (
	"fmt"
	"os"

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
		return out
	} else if out := types.GetK8Slookup(); out != nil {
		return out
	}
	fmt.Println("got bad k8s resource type")
	os.Exit(2)

	return nil
}

func (stack *K8SStack) Synth(name string) (*dag.DAG, map[string]base.ResourceType) {
	objsMap := map[string]base.ResourceType{}
	objDag := dag.NewDAG()
	for _, obj := range stack.Objects {
		objsMap[obj.ActualType().GetID()] = obj.ActualType()
		obj.ActualType().AddToDag(objDag)
	}

	return objDag, objsMap
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
