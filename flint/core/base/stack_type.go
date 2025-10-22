package base

import (
	"github.com/heimdalr/dag"
)

type StackType interface {
	String() string
	Synth() (*dag.DAG, map[string]map[string]any)
	GetConnection() Connection
	Deploy()
}
