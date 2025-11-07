package common

import "github.com/heimdalr/dag"

type StackType interface {
	Synth(string) (*dag.DAG, map[string]map[string]any)
}

func (stackType *StackTypes) GetActual() StackType {
	if out := stackType.GetK8SStack(); out != nil {
		return out
	}

	panic("got bad stack type")
}

type ConnectionType interface {
	Deploy(*dag.DAG, map[string]map[string]any, string)
	Diff(map[string]map[string]any, string) ([]string, []string, [][]map[string]any)
	ToFileName(map[string]any) string
}

func (connType *ConnectionTypes) GetActual() ConnectionType {
	if out := connType.GetK8SConnection(); out != nil {
		return out
	}

	panic("got bad connection type")
}
