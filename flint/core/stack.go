package core

import (
	"github.com/SuperTapood/Flint/core/base"
	"github.com/SuperTapood/Flint/core/generated/common"
	"github.com/google/uuid"
	"github.com/heimdalr/dag"
	"google.golang.org/protobuf/proto"
)

type Stack struct {
	Stack_ *common.Stack_
}

func StackFromBinary(data []byte) Stack {
	var stack common.Stack_
	if err := proto.Unmarshal(data, &stack); err != nil {
		panic(err)
	}

	return Stack{
		Stack_: &stack,
	}
}

func (stack *Stack) ActualStack() base.StackType {
	if out := stack.Stack_.Stack[0].GetK8SStack(); out != nil {
		return out
	}
	panic("got bad stack")
}

func (stack *Stack) String() string {
	return stack.ActualStack().String()
}

func (stack *Stack) Synth() (*dag.DAG, map[uuid.UUID]map[string]any) {
	return stack.ActualStack().Synth()
}

func (stack *Stack) Deploy() {
	stack.ActualStack().Deploy()
}
