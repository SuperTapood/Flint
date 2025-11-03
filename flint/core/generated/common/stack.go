package common

import (
	"github.com/SuperTapood/Flint/core/base"
	"github.com/heimdalr/dag"
	"google.golang.org/protobuf/proto"
)

// type Stack struct {
// 	Stack_ *common.Stack_
// }

func StackFromBinary(data []byte) *Stack {
	var stack Stack
	if err := proto.Unmarshal(data, &stack); err != nil {
		panic(err)
	}

	return &stack
}

func (stack *Stack) ActualStack() base.StackType {
	if out := stack.Stack[0].GetK8SStack(); out != nil {
		return out
	}
	panic("got bad stack")
}

// func (stack *Stack) String() string {
// 	return stack.ActualStack().String()
// }

func (stack Stack) Synth() (*dag.DAG, map[string]map[string]any) {
	return stack.ActualStack().Synth()
}

func (stack *Stack) Deploy() {
	stack.ActualStack().Deploy()
}
