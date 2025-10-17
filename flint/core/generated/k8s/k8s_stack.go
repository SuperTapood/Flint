package k8s

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"io"

	"github.com/SuperTapood/Flint/core/base"
	"github.com/google/uuid"
	"github.com/heimdalr/dag"
)

func (types *K8STypes) ActualType() base.ResourceType {
	if out := types.GetPod(); out != nil {
		return out
	}
	panic("got bad type")
}

func (types *K8STypes) Synth() (uuid.UUID, map[string]any) {
	return types.ActualType().Synth()
}

func (stack *K8S_Stack_) Synth() (*dag.DAG, map[uuid.UUID]map[string]any) {
	objs_map := map[uuid.UUID]map[string]any{}
	var obj_dag = dag.NewDAG()
	for _, obj := range stack.Objects {
		var uuid, obj_map = obj.Synth()
		// log.Printf("%d: %s", i, uuid)
		objs_map[uuid] = obj_map
		obj_dag.AddVertex(uuid)
	}

	return obj_dag, objs_map
}

func (stack *K8S_Stack_) Deploy() {
	var _, obj_map = stack.Synth()

	for _, v := range obj_map {
		data, _ := json.Marshal(v)
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

		var location = fmt.Sprint(v["location"])
		req, err := http.NewRequest("POST", "https://192.168.49.2:8443"+location, bytes.NewReader(data))
		req.Header.Add("Authorization", "Bearer eyJhbGciOiJSUzI1NiIsImtpZCI6ImxXUGU0UEIwZWtaRVlXaHM5TEVENmFzV2FSQTJPRi1ndkVHQ2hOUlNEdWMifQ.eyJhdWQiOlsiaHR0cHM6Ly9rdWJlcm5ldGVzLmRlZmF1bHQuc3ZjLmNsdXN0ZXIubG9jYWwiXSwiZXhwIjoxNzYwNzEyODg4LCJpYXQiOjE3NjA3MDkyODgsImlzcyI6Imh0dHBzOi8va3ViZXJuZXRlcy5kZWZhdWx0LnN2Yy5jbHVzdGVyLmxvY2FsIiwianRpIjoiYzVkOGU3NjQtZTY2Zi00ZmQ0LWJjZDQtMGNjNDJlYWM4ZjNlIiwia3ViZXJuZXRlcy5pbyI6eyJuYW1lc3BhY2UiOiJkZWZhdWx0Iiwic2VydmljZWFjY291bnQiOnsibmFtZSI6ImZsaW50IiwidWlkIjoiNTNlMDM5NzYtOWQxNi00MjgzLTlkYTAtM2QxZTIxMWQ5YmM2In19LCJuYmYiOjE3NjA3MDkyODgsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDpkZWZhdWx0OmZsaW50In0.wyvJ1OJodJ1BEFGAzpWQ6dI9Beg3qBpExpfCu5aCDOwJHc4twopWfUSxjnnLYd3dBked8xGG1fCfPKmySUYR12zW6HoNsvtpl9XEeXyqKckB9lap3f6xr8I4eEqXFK7o2d85mDLqLlFN_iye1T2fVyyZbpcdxKlPnJexyXUnWXouZvFLuU8nvng31kMBAU0Y5RR2lCQDoVBqlQYhP3mzDAgUbxb23QlPzr0-jS82JVoMScQ2V9Rh4bEV-3D-RhfrBlCiDQ2KhOYUyS0Z7aqWxMiob1jbp4nElgi0qkwXdIxwYb25bFGNqaj9vdMd3Vs4mjDWD-zZZLG5tiMLVWbaLw")
		// resp, err := http.Post(, "application/json", )
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
