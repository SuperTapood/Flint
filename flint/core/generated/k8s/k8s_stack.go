package k8s

import (
	"encoding/json"
	"strconv"
	sync "sync"

	"github.com/SuperTapood/Flint/core/base"
	"github.com/heimdalr/dag"
)

func (types *K8STypes) ActualType() base.ResourceType {
	if out := types.GetPod(); out != nil {
		return out
	} else if out := types.GetService_(); out != nil {
		return out
	} else if out := types.GetDeployment(); out != nil {
		return out
	}
	panic("got bad resource type")
}

func (types *K8STypes) Synth(stack_name string, namespace string, dag *dag.DAG, obj_map map[string]map[string]any) {
	types.ActualType().Synth(stack_name, namespace, dag, obj_map)
}

func (stack *K8S_Stack_) GetConnection() base.Connection {
	return &base.K8SConnection{
		Api:   stack.GetApi(),
		Token: stack.GetToken(),
	}
}

func (stack *K8S_Stack_) Synth() (*dag.DAG, map[string]map[string]any) {
	objs_map := map[string]map[string]any{}
	var obj_dag = dag.NewDAG()
	for _, obj := range stack.Objects {
		obj.Synth(stack.GetName(), stack.GetNamespace(), obj_dag, objs_map)
	}

	return obj_dag, objs_map
}

// SimpleVisitor collects nodes in topological order.
type SimpleVisitor struct {
	Order []string
}

func (v *SimpleVisitor) Visit(vertexer dag.Vertexer) {
	id, _ := vertexer.Vertex()
	v.Order = append(v.Order, id)
}

func (stack *K8S_Stack_) process(obj map[string]any) {
	stack.GetConnection().Deploy(obj)
}

// ProcessDAGParallel processes a DAG in parallel, respecting dependencies
func ProcessDAGParallel(d *dag.DAG, workFunc func(id string, vertex interface{}) error) error {
	// Get all vertices
	vertices := d.GetVertices()

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
			descendants, err := d.GetDescendants(id)
			if err != nil {
				mu.Unlock()
				return err
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
			go func(nodeID string) {
				defer wg.Done()

				vertex, err := d.GetVertex(nodeID)
				if err != nil {
					errChan <- err
					return
				}

				if err := workFunc(nodeID, vertex); err != nil {
					errChan <- err
					return
				}

				mu.Lock()
				completed[nodeID] = true
				mu.Unlock()
			}(id)
		}

		wg.Wait()
		close(errChan)

		// Check for errors
		for err := range errChan {
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (stack *K8S_Stack_) Deploy() {
	var dag, obj_map = stack.Synth()

	install_number := stack.GetConnection().GetCurrentRevision(stack.GetName())

	secret := Secret{
		Name: stack.GetName() + "-" + strconv.Itoa(install_number),
		Type: "v1.flint.io",
		Data: make([]*SecretData, 2),
	}

	marshalled, _ := json.Marshal(obj_map)

	data := SecretData{
		Key:   "data",
		Value: string(marshalled),
	}

	status := SecretData{
		Key:   "status",
		Value: "success",
	}

	secret.Data[0] = &data
	secret.Data[1] = &status

	secret.Synth(stack.GetName(), stack.GetNamespace(), dag, obj_map)
	conn := stack.GetConnection()

	err := ProcessDAGParallel(dag, func(id string, vertex interface{}) error {
		conn.Deploy(obj_map[id])
		return nil
	})

	if err != nil {
		panic(err)
	}
}
