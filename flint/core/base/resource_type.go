package base

import "github.com/heimdalr/dag"

type ResourceType interface {
	String() string
	GetID() string
	Synth(*dag.DAG) map[string]any
}
