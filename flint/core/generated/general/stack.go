package general

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"slices"
	"strings"
	"sync"

	"github.com/SuperTapood/Flint/core/base"
	"github.com/SuperTapood/Flint/core/generated/common"
	"github.com/SuperTapood/Flint/core/util"
	"github.com/google/uuid"
	"github.com/heimdalr/dag"
)

// a generic representation of a stack
type StackType interface {
	// synth this stack and return the resulting dag and object map
	Synth(string) (*dag.DAG, map[string]base.ResourceType)
	// get the useful metadata of this stack
	GetMetadata() map[string]any
}

// resolve a StackTypes object to a StackType
func (stackType *StackTypes) GetActual() StackType {
	if out := stackType.GetK8Sstack(); out != nil {
		return out
	}

	fmt.Println("got bad stack type")
	os.Exit(2)

	return nil
}

// a generic representation of a connection
type ConnectionType interface {
	// ToFileName(id string) string
	// List() []gen_base.FlintDeployment
	// GetCurrentRevision(stackName string) int
	PrettyName(resource map[string]any, stackMetadata map[string]any) string
	// Diff(resources map[string]base.ResourceType, stackMetadata map[string]any, stackName string) ([]string, []string, []map[string]map[string]any)
	CleanHistory(stackName string, oldest int, stackMetadata map[string]any)
	Apply(applyMetadata map[string]any, resource map[string]any)
	// MakeRequest(method string, location string, reader io.Reader) ([]byte, *http.Response)
	GetClient() *util.HttpClient
	CreateRevision(stackName string, stackMetadata map[string]any, newDag *dag.DAG, marshalled []byte)
	GetRevisions() map[string]map[string]any
	GetLatestRevision(stackName string) (map[string]any, string, string, int32)
	GetDeployments() []string
	Delete(delete_metadata map[string]any)
	PrintOutputs()
	// Deploy(dag *dag.DAG, resources map[string]base.ResourceType, stackName string, stackMetadata map[string]any, max_revisions int)
	// Destroy(stackName string, stackMetadata map[string]any)
	// Rollback(stackName string, targetRevision int, stackMetadata map[string]any)
}

func (connType *ConnectionTypes) GetActual() ConnectionType {
	if out := connType.GetK8Sconnection(); out != nil {
		return out
	}

	fmt.Println("got bad connection type")
	os.Exit(2)

	return nil
}

func (connType *ConnectionTypes) List() []common.FlintDeployment {
	secrets := connType.GetActual().GetDeployments()
	deployments := []common.FlintDeployment{}
	visited := []string{}
	for _, name := range secrets {
		if slices.Contains(visited, name) {
			continue
		}

		_, status, age, version := connType.GetActual().GetLatestRevision(name)
		deployments = append(deployments, common.FlintDeployment{
			Name:     name,
			Age:      age,
			Status:   status,
			Revision: version,
		})
		visited = append(visited, name)
	}

	return deployments
}

func (connType *ConnectionTypes) GetCurrentRevision(stackName string) int {
	_, _, _, version := connType.GetActual().GetLatestRevision(stackName)
	return int(version)
}

func (connType *ConnectionTypes) Deploy(_dag *dag.DAG, resources map[string]base.ResourceType, stackName string, stackMetadata map[string]any, max_revisions int, createRevision bool) {
	install_number := connType.GetCurrentRevision(stackName) + 1
	lowest_revision := max(install_number-max_revisions, 0) + 1

	connType.GetActual().CleanHistory(stackName, lowest_revision, stackMetadata)

	// Get all vertices
	vertices := _dag.GetVertices()

	current := 1
	total := len(_dag.GetVertices())

	prettys := make([]string, 0)

	for id, res := range resources {
		if resources[id].Synth(stackMetadata) == nil {
			total -= 1
			continue
		}
		prettys = append(prettys, connType.GetActual().PrettyName(res.Synth(stackMetadata), stackMetadata))
	}

	deployPrint := util.CreateDeployPrint(stackName, prettys, stackMetadata)

	completed := make(map[string]bool)
	var mu sync.Mutex

	for len(completed) < len(vertices) {
		// Find nodes ready to process (all dependencies completed)
		var ready []string
		for id := range vertices {
			mu.Lock()
			if completed[id] {
				mu.Unlock()
				continue
			}

			// Get descendants (dependencies)
			descendants, err := _dag.GetDescendants(id)
			if err != nil {
				mu.Unlock()
				fmt.Println("failed to get descendants from dag")
				fmt.Println(err)
				os.Exit(-1)
			}

			// Check if all dependencies are completed
			allDepsComplete := true
			for depID := range descendants {
				if !completed[depID] {
					allDepsComplete = false
					break
				}
			}
			mu.Unlock()

			if allDepsComplete {
				ready = append(ready, id)
			}
		}

		if len(ready) == 0 {
			break // No more nodes to process
		}

		// Process ready nodes in parallel
		var wg sync.WaitGroup
		errChan := make(chan error, len(ready))

		for _, id := range ready {
			wg.Add(1)
			go func(nodeID string, idx int) {
				defer wg.Done()

				vertex, err := _dag.GetVertex(nodeID)
				if err != nil {
					errChan <- err
					return
				}

				res := resources[vertex.(string)]
				synthed := res.Synth(stackMetadata)

				if synthed != nil {
					action := "CREATING"
					if synthed["action"] == "delete" {
						action = "DELETING"
					}
					deployPrint.PrettyPrint(stackName, idx, total, action, connType.GetActual().PrettyName(synthed, stackMetadata))
				}

				res.Apply(stackMetadata, resources, connType.GetActual())
				if synthed != nil {
					action := "CREATED"
					if synthed["action"] == "delete" {
						action = "DELETED"
					}
					deployPrint.PrettyPrint(stackName, idx, total, action, connType.GetActual().PrettyName(synthed, stackMetadata))
				}

				mu.Lock()
				completed[nodeID] = true
				mu.Unlock()
			}(id, current)
			if resources[id].Synth(stackMetadata) != nil {
				current += 1
			}
		}

		wg.Wait()
		close(errChan)

		// Check for errors
		for err := range errChan {
			if err != nil {
				fmt.Println("deploy failed with error:")
				fmt.Println(err)
				os.Exit(-1)
			}
		}
	}

	connType.GetActual().PrintOutputs()

	if !createRevision {
		return
	}

	newDag := dag.NewDAG()

	b, err := json.Marshal(_dag)
	if err != nil {
		fmt.Println("failed to marshal dag")
		fmt.Println(err)
		os.Exit(-1)
	}

	var dagMap map[string]any
	err = json.Unmarshal(b, &dagMap)
	if err != nil {
		fmt.Println("failed to marshal dag map")
		fmt.Println(err)
		os.Exit(-1)
	}

	removeFromDag := make([]string, 0)
	for name, obj := range resources {
		if obj.Synth(stackMetadata)["kind"] == "" {
			// delete(obj_map, name)
			removeFromDag = append(removeFromDag, name)
		}
	}

	for _, value := range dagMap["vs"].([]any) {
		i := value.(map[string]any)["i"].(string)
		v := value.(map[string]any)["v"].(string)
		if slices.Contains(removeFromDag, i) || slices.Contains(removeFromDag, v) {
			continue
		}
		newDag.AddVertexByID(v, i)
	}

	for _, value := range dagMap["es"].([]any) {
		d := value.(map[string]any)["d"].(string)
		s := value.(map[string]any)["s"].(string)
		if slices.Contains(removeFromDag, d) || slices.Contains(removeFromDag, s) {
			continue
		}
		newDag.AddEdge(s, d)
	}
	objs_map := make(map[string]map[string]any, 0)

	for _, resource := range resources {
		synthed := resource.Synth(stackMetadata)
		if synthed == nil && synthed["action"] == nil {
			continue
		}
		objs_map[resource.GetID()] = synthed
	}

	marshalled, _ := json.Marshal(objs_map)

	connType.GetActual().CreateRevision(stackName, stackMetadata, newDag, marshalled)

}

func (connType *ConnectionTypes) Diff(resources map[string]base.ResourceType, stackMetadata map[string]any, stackName string) ([]string, []string, []map[string]map[string]any) {
	obj_map, _, _, version := connType.GetActual().GetLatestRevision(stackName)

	added := make([]string, 0)

	if version == 0 {
		for _, res := range resources {
			synthed := res.Synth(stackMetadata)
			if len(synthed) == 0 {
				continue
			}
			added = append(added, connType.GetActual().PrettyName(synthed, stackMetadata))
		}
		return added, make([]string, 0), make([]map[string]map[string]any, 0)
	}

	removed := make([]string, 0)
	changed := make([]map[string]map[string]any, 0)

	for newName, newObj := range resources {
		found := false
		if newObj.Synth(stackMetadata) == nil {
			continue
		}
		var foundObjKey string
		for name := range obj_map {
			if name == newName {
				found = true
				foundObjKey = name
				break
			}
		}

		if found {
			bytesOld, err := json.Marshal(obj_map[foundObjKey])
			if err != nil {
				fmt.Println("failed to marshal found object")
				fmt.Println(err)
				os.Exit(-1)
			}

			bytesNew, err := json.Marshal(newObj.Synth(stackMetadata))
			if err != nil {
				fmt.Println("failed to marshal synthed new objet")
				fmt.Println(err)
				os.Exit(-1)
			}
			if !strings.EqualFold(string(bytesOld), string(bytesNew)) {
				objects := make(map[string]map[string]any, 2)
				objects["new"] = newObj.Synth(stackMetadata)
				objects["old"] = obj_map[foundObjKey].(map[string]any)
				changed = append(changed, objects)
			}

			delete(obj_map, foundObjKey)
		} else {
			added = append(added, connType.GetActual().PrettyName(newObj.Synth(stackMetadata), stackMetadata))
		}
	}

	for _, obj := range obj_map {
		removed = append(removed, connType.GetActual().PrettyName(obj.(map[string]any), stackMetadata))
	}

	return added, removed, changed
}

func (connType *ConnectionTypes) Destroy(stackName string, stackMetadata map[string]any) {
	obj_map, _, _, version := connType.GetActual().GetLatestRevision(stackName)

	if version == 0 {
		return
	}

	removes := make(map[string]base.ResourceType, 0)
	obj_dag := dag.NewDAG()

	for _, resource := range obj_map {
		unresource := base.Unresource{
			Name: connType.GetActual().PrettyName(resource.(map[string]any), stackMetadata),
			ID:   uuid.New().String(),
		}
		removes[unresource.GetID()] = &unresource
		obj_dag.AddVertexByID(unresource.GetID(), unresource.GetID())
	}

	connType.Deploy(obj_dag, removes, stackName, stackMetadata, math.MaxInt, false)
}
