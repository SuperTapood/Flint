package base

import (
	"strings"

	"github.com/heimdalr/dag"
)

type Unresource struct {
	Name string
	ID   string
}

func (unresource *Unresource) String() string {
	return unresource.ID
}

func (unresource *Unresource) GetID() string {
	return unresource.String()
}

func (unresource *Unresource) GetPrettyName(stack_metadata map[string]any) string {
	// return "Kubernetes::Pod::" + stack_metadata["namespace"].(string) + "::" + pod.GetName()
	return unresource.String()
}

func (unresource *Unresource) Synth(stack_metadata map[string]any) map[string]any {
	return nil
}

func (unresource *Unresource) AddToDag(dag *dag.DAG) {
	panic("what the fuck are you doing")
}

func (unresource *Unresource) Apply(stack_metadata map[string]any, resources map[string]ResourceType, client CloudClient) {
	splitName := strings.Split(unresource.Name, "::")
	namespace := splitName[1]
	kind := splitName[2]
	name := splitName[3]

	client.Delete(map[string]any{
		"kind":      kind,
		"namespace": namespace,
		"name":      name,
	})
}

func (unresource *Unresource) Lookup() map[string]any {
	panic("what the fuck are you doing")
}
