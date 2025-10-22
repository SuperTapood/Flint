package k8s

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

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

func (stack *K8S_Stack_) Deploy() {
	var dag, obj_map = stack.Synth()

	log.Print(dag.String())
	log.Print(obj_map)

	for _, v := range obj_map {
		data, _ := json.Marshal(v)
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

		var location = fmt.Sprint(v["location"])
		req, err := http.NewRequest("POST", stack.GetApi()+location, bytes.NewReader(data))
		req.Header.Add("Authorization", "Bearer "+stack.GetToken())
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Println("Error on response.\n[ERROR] -", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println("Error while reading the response bytes:", err)
		}
		log.Println(string([]byte(body)))
	}
}
