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
	var out = stack.Stack_.Stack.GetK8SStack()
	if out != nil {
		return out
	} else {
		panic("got bad stack")
	}
}
