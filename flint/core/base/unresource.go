package base

import (
	"fmt"
	"os"
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

func (unresource *Unresource) Synth(stackMetadata map[string]any) map[string]any {
	splitName := strings.Split(unresource.Name, "::")
	namespace := splitName[1]
	kind := splitName[2]
	name := splitName[3]
	return map[string]any{
		"action":    "delete",
		"kind":      kind,
		"namespace": namespace,
		"metadata": map[string]any{
			"name": name,
		},
	}
}

func (unresource *Unresource) AddToDag(dag *dag.DAG) {
	fmt.Println("Unresource cannot be added to a dag like this")
	os.Exit(2)
}

func (unresource *Unresource) Apply(stackMetadata map[string]any, resources map[string]ResourceType, client CloudClient) {
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
	fmt.Println("Unresource cannot be looked up")
	os.Exit(2)
	return nil
}
