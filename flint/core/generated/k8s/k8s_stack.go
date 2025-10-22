package k8s

import (
	"log"

	"github.com/SuperTapood/Flint/core/base"
	"github.com/heimdalr/dag"
)

func (types *K8STypes) ActualType() base.ResourceType {
	if out := types.GetPod(); out != nil {
		return out
	}
	panic("got bad type")
}

func (types *K8STypes) Synth(dag *dag.DAG) map[string]any {
	return types.ActualType().Synth(dag)
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
		var obj_map = obj.Synth(obj_dag)
		// log.Printf("%d: %s", i, uuid)
		objs_map[obj.ActualType().GetID()] = obj_map
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

func process(id string) {

}

func (stack *K8S_Stack_) Deploy() {
	var dag, obj_map = stack.Synth()

	log.Print(dag.String())
	log.Print(obj_map)

	// for _, v := range obj_map {
	// 	stack.GetConnection().Deploy(v)
	// }

	visitor := &SimpleVisitor{}
	dag.OrderedWalk(visitor)

	// Iterate over sorted nodes
	var connection = stack.GetConnection()
	for _, node := range visitor.Order {
		var obj = obj_map[node]
		log.Print(obj)
		connection.Deploy(obj)
	}

	// // Prepare for parallel processing: compute indegrees atomically
	// indegree := make(map[string]*int32)
	// vertices := dag.GetVertices()
	// for id := range vertices {
	// 	parents, _ := dag.GetParents(id) // Ignore error for example
	// 	deg := int32(len(parents))
	// 	indegree[id] = &deg
	// }

	// // Define the process function (runs your node logic, then notifies children)
	// var wg sync.WaitGroup
	// process := func(id string) {
	// 	// Your node processing here (e.g., simulate work)
	// 	fmt.Printf("Processing %s\n", id)

	// 	// Get children (thread-safe)
	// 	children, _ := dag.GetChildren(id) // Ignore error for example

	// 	// Notify children: decrement their indegree
	// 	for childID := range children {
	// 		if atomic.AddInt32(indegree[childID], -1) == 0 {
	// 			// Child is ready (all parents done), spawn its goroutine
	// 			wg.Add(1)
	// 			go process(childID)
	// 		}
	// 	}

	// 	wg.Done()
	// }

	// // Start from roots (indegree 0)
	// roots := dag.GetRoots()
	// for rootID := range roots {
	// 	wg.Add(1)
	// 	go process(rootID)
	// }

	// // Wait for all nodes to finish
	// wg.Wait()

	// fmt.Println("All processing complete")
	// // Possible output order: A, then B and C (parallel), then D
}
