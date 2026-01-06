package base

import (
	"github.com/SuperTapood/Flint/core/util"
	"github.com/heimdalr/dag"
)

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
	AddToDag(_dag *dag.DAG)

	Synth(stackMetadata map[string]any) map[string]any
	Apply(stackMetadata map[string]any, resources map[string]ResourceType, client CloudClient) error
	ExplainFailure(client *util.HttpClient, stackMetadata map[string]any) string
}
