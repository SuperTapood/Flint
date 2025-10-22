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

	// for _, v := range obj_map {
	// 	data, _ := json.Marshal(v)
	// 	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	// 	var location = fmt.Sprint(v["location"])
	// 	req, err := http.NewRequest("POST", "https://192.168.49.2:8443"+location, bytes.NewReader(data))
	// 	req.Header.Add("Authorization", "Bearer eyJhbGciOiJSUzI1NiIsImtpZCI6ImxXUGU0UEIwZWtaRVlXaHM5TEVENmFzV2FSQTJPRi1ndkVHQ2hOUlNEdWMifQ.eyJhdWQiOlsiaHR0cHM6Ly9rdWJlcm5ldGVzLmRlZmF1bHQuc3ZjLmNsdXN0ZXIubG9jYWwiXSwiZXhwIjoxNzYxMTY0NDE1LCJpYXQiOjE3NjExNjA4MTUsImlzcyI6Imh0dHBzOi8va3ViZXJuZXRlcy5kZWZhdWx0LnN2Yy5jbHVzdGVyLmxvY2FsIiwianRpIjoiNDFkMmJiM2QtN2ZmYy00YmExLTg4NzctMWVhNjgwOTg4MTY0Iiwia3ViZXJuZXRlcy5pbyI6eyJuYW1lc3BhY2UiOiJkZWZhdWx0Iiwic2VydmljZWFjY291bnQiOnsibmFtZSI6ImZsaW50IiwidWlkIjoiNTNlMDM5NzYtOWQxNi00MjgzLTlkYTAtM2QxZTIxMWQ5YmM2In19LCJuYmYiOjE3NjExNjA4MTUsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDpkZWZhdWx0OmZsaW50In0.owxBYV0SPMfSoFto0rmKUt6zP3NdrVyudzrA404N24n-nA2DWobFx-k-6OQkovWkOjvteqCFSj_qqp3TdzuBUXkBIMbMMLzVvcVuIVyFvXl9K39Ru9pY4Cn0hCt3dYgLGvi_M079g97gwLD99sCwLxH1f4TF9DWX4gvmjqcqNGGkJt0zG4W65xVKLxx0dCVgCiSh-oihgE0tP78AthF4CVJ6mudavfQugnhoBh14jKkI3aybsoZJL-GpvmrKSr1KeQxuK38XWxFCbbE9kxkOpsuDYpOWlQUG8f8MLMf1Y3UT5DN7nKXPCENPE7-8ccdlA6SPDMRv_2hWb779fgYlig")
	// 	// resp, err := http.Post(, "application/json", )
	// 	client := &http.Client{}
	// 	resp, err := client.Do(req)
	// 	if err != nil {
	// 		log.Println("Error on response.\n[ERROR] -", err)
	// 	}
	// 	defer resp.Body.Close()

	// 	body, err := io.ReadAll(resp.Body)
	// 	if err != nil {
	// 		log.Println("Error while reading the response bytes:", err)
	// 	}
	// 	log.Println(string([]byte(body)))
	// }
}
