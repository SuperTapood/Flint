package common

import "github.com/heimdalr/dag"

// a generic representation of a stack
type StackType interface {
	// synth this stack and return the resulting dag and object map
	Synth(string) (*dag.DAG, map[string]map[string]any)
	// get the useful metadata of this stack
	GetMetadata() map[string]any
}

// resolve a StackTypes object to a StackType
func (stackType *StackTypes) GetActual() StackType {
	if out := stackType.GetK8SStack(); out != nil {
		return out
	}

	panic("got bad stack type")
}

type ConnectionType interface {
	Deploy(*dag.DAG, []string, map[string]map[string]any, string, map[string]any, int)
	Diff(map[string]map[string]any, string) ([]string, []string, [][]map[string]any)
	ToFileName(map[string]any) string
	Destroy(string, map[string]any)
	GetCurrentRevision(string) int
	Rollback(string, int, map[string]any)
}

func (connType *ConnectionTypes) GetActual() ConnectionType {
	if out := connType.GetK8SConnection(); out != nil {
		return out
	}

	panic("got bad connection type")
}
