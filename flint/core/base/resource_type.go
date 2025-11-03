package base

import "github.com/heimdalr/dag"

type ResourceType interface {
	String() string
	GetID() string
	Synth(string, string, *dag.DAG, map[string]map[string]any)
}
