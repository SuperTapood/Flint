package base

import "github.com/heimdalr/dag"

// a general representation of a synthable resource
type ResourceType interface {
	// return the string representation of this object (usually implemented by protobuf)
	String() string
	// return the id of the current object. The id is very opinionated.
	GetID() string
	/*
		synth this object and add it to both the stack and its directed acyclic graph

		Parameters:
			- map[string]any - stack metadata
			- *dag.DAG - the stack's dag to be modified personally by the object
			- map[string] - the object map to be deployed to the cloud provider

	*/
	AddToDag(dag *dag.DAG)
	// return an opinionated map of this object's important properties from the cloud provider
	Lookup() map[string]any

	Synth(stack_metadata map[string]any) map[string]any
	Apply(stack_metadata map[string]any, resources map[string]ResourceType, client CloudClient)
}
