package k8s

import (
	"log"

	"github.com/SuperTapood/Flint/core/base"
	"github.com/google/uuid"
	"github.com/heimdalr/dag"
)

func (stack *K8S_Stack_) Synth() (dag.DAG, map[uuid.UUID]base.ResourceType) {
	for i, obj := range stack.Objects {
		log.Printf("%d: %s", i, obj.String())
	}

	return *dag.NewDAG(), nil
}
