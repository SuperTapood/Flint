package core

import (
	"github.com/SuperTapood/Flint/core/generated/common"
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

func (stack *Stack) GetStack() StackType {
	if out := stack.Stack_.Stack[0].GetK8SStack(); out != nil {
		return out
	}
	panic("got bad stack")
}

func (stack *Stack) String() string {
	return stack.GetStack().String()
}
