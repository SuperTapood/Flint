package base

import (
	"github.com/google/uuid"
	"github.com/heimdalr/dag"
)

type StackType interface {
	String() string
	Synth() (*dag.DAG, map[uuid.UUID]map[string]any)
	Deploy()
}
