package base

import "github.com/heimdalr/dag"

// a general representation of a synthable resource
type ResourceType interface {
	// return the string representation of this object (usually implemented by protobuf)
	String() string
	// return the id of the current object. The id is very opinionated.
	GetID() string
	// synth this object and add it to both the stack and its directed acyclic graph
	Synth(string, string, *dag.DAG, map[string]map[string]any)
	// return an opinionated map of this object's important properties from the cloud provider
	Lookup() map[string]any
}
