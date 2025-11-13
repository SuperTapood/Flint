package k8s

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"slices"
	"strconv"
	"strings"
	sync "sync"
	"time"

	"github.com/SuperTapood/Flint/core/generated/gen_base"
	"github.com/heimdalr/dag"
)

func (connection *K8S_Connection) getSecrets() map[string]any {
	var body, _ = connection.makeRequest("GET", "/api/v1/secrets", bytes.NewReader(make([]byte, 1)))
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		panic(err)
	}

	return result
}

func (connection *K8S_Connection) getLatestSecret(stack_name string) (map[string]any, int32) {
	result := connection.getSecrets()
	latest_version := 0
	var latest_secret map[string]any

	for _, secret := range result["items"].([]any) {
		if secret.(map[string]any)["type"] == "v1.flint.io" {
			secret_name := secret.(map[string]any)["metadata"].(map[string]any)["name"].(string)
			re := regexp.MustCompile(stack_name + `-[0-9]+`)
			if re.FindString(secret_name) != "" {
				version_re := regexp.MustCompile(`[0-9]+`)
				version, err := strconv.Atoi(version_re.FindString(secret_name))
				if err != nil {
					panic(err)
				}
				if version > latest_version {
					latest_secret = secret.(map[string]any)
					latest_version = version
				}
			}
		}
	}

	return latest_secret, int32(latest_version)
}

func (connection *K8S_Connection) CleanHistory(stack_name string, oldest int, namespace string) {
	if oldest < 1 {
		return
	}

	result := connection.getSecrets()

	for _, secret := range result["items"].([]any) {
		if secret.(map[string]any)["type"] == "v1.flint.io" {
			secret_name := secret.(map[string]any)["metadata"].(map[string]any)["name"].(string)
			re := regexp.MustCompile(stack_name + `-[0-9]+`)
			if re.FindString(secret_name) != "" {
				version_re := regexp.MustCompile(`[0-9]+`)
				version, err := strconv.Atoi(version_re.FindString(secret_name))
				if err != nil {
					panic(err)
				}
				if version < oldest {
					body, resp := connection.makeRequest("DELETE", locationMap["Secret"]["before_namespace"]+namespace+locationMap["Secret"]["after_namespace"]+secret_name, bytes.NewReader(make([]byte, 0)))
					if resp.StatusCode != http.StatusOK {
						fmt.Println(string(body))
						fmt.Println(resp)
						panic("couldn't delete secret")
					}
				}
			}
		}
	}
}

func (connection *K8S_Connection) GetCurrentRevision(stack_name string) int {
	_, latest_version := connection.getLatestSecret(stack_name)
	return int(latest_version)
}

func (connection *K8S_Connection) makeRequest(method string, location string, reader io.Reader) ([]byte, *http.Response) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	req, _ := http.NewRequest(method, connection.Api+location, reader)

	req.Header.Add("Authorization", "Bearer "+connection.Token)
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
	}
	defer resp.Body.Close()

	body, er := io.ReadAll(resp.Body)
	if er != nil {
		log.Println("Error while reading the response bytes:", err)
	}

	// fmt.Println(string(body))

	return body, resp
}

const (
	APIS_APP_V1 = "/apis/apps/v1/namespaces/"
	API_V1      = "/api/v1/namespaces/"
)

var (
	locationMap map[string]map[string]string = map[string]map[string]string{
		"Deployment": {
			"before_namespace": APIS_APP_V1,
			"after_namespace":  "/deployments/",
		},
		"Service": {
			"before_namespace": API_V1,
			"after_namespace":  "/services/",
		},
		"Secret": {
			"before_namespace": API_V1,
			"after_namespace":  "/secrets/",
		},
		"Pod": {
			"before_namespace": API_V1,
			"after_namespace":  "/pods/",
		},
	}
)

func (connection *K8S_Connection) Apply(obj map[string]any, objs_map map[string]map[string]any) {
	action, ok := obj["action"]
	if !ok {
		action = "deploy"
	}

	switch action {
	case "remove":
		split_name := strings.Split(obj["name"].(string), "::")
		kind := split_name[1]
		namespace := split_name[2]
		name := split_name[3]
		connection.makeRequest("DELETE", locationMap[kind]["before_namespace"]+namespace+locationMap[kind]["after_namespace"]+name, bytes.NewReader(make([]byte, 0)))
		// fmt.Println(string(body))
		return
	case "lookup":
		lookups := obj["lookups"].([]*Lookup)
		strings := obj["strings"].([]string)
		length := max(len(lookups), len(strings))
		for i := range length {
			if i < len(strings) {
				fmt.Print(strings[i])
			}
			if i < len(lookups) {
				lookup := lookups[i]
				var lookup_id = lookup.GetObject().ActualType().GetID()
				target := objs_map[lookup_id]
				kind := target["kind"].(string)
				namespace := target["metadata"].(map[string]any)["namespace"].(string)
				name := target["metadata"].(map[string]any)["name"].(string)
				body, _ := connection.makeRequest("GET", locationMap[kind]["before_namespace"]+namespace+locationMap[kind]["after_namespace"]+name, bytes.NewReader(make([]byte, 1)))
				var currentMap map[string]any
				err := json.Unmarshal(body, &currentMap)
				if err != nil {
					panic(err)
				}
				var current any = currentMap
				for _, k := range lookup.GetKeys() {
					// must be a map to go deeper
					mmap, ok := current.(map[string]any)
					if !ok {
						panic("badbad")
					}
					v, ok := mmap[k]
					if !ok {
						panic("badbad")
					}
					current = v
				}

				fmt.Print(current)
			}
		}
		fmt.Println()
		return
	}

	data, _ := json.Marshal(obj)
	kind := obj["kind"].(string)
	namespace := obj["metadata"].(map[string]any)["namespace"].(string)
	name := obj["metadata"].(map[string]any)["name"].(string)

	var resp *http.Response
	var body []byte

	body, resp = connection.makeRequest("POST", locationMap[kind]["before_namespace"]+namespace+locationMap[kind]["after_namespace"], bytes.NewReader(data))

	if resp.StatusCode == http.StatusConflict {
		body, resp = connection.makeRequest("PUT", locationMap[kind]["before_namespace"]+namespace+locationMap[kind]["after_namespace"]+name, bytes.NewReader(data))
		if resp.StatusCode == http.StatusUnprocessableEntity {
			body, resp = connection.makeRequest("DELETE", locationMap[kind]["before_namespace"]+namespace+locationMap[kind]["after_namespace"]+name, bytes.NewReader(make([]byte, 0)))
			if resp.StatusCode == http.StatusOK {
				time.Sleep(2 * time.Second)
				body, resp = connection.makeRequest("POST", locationMap[kind]["before_namespace"]+namespace+locationMap[kind]["after_namespace"], bytes.NewReader(data))
			}
		}
	}

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		fmt.Println("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
		fmt.Println(locationMap[kind])
		fmt.Println(resp)
		fmt.Println(string(body))
		fmt.Println(obj)
		panic("FUCK")
	}
}

func (connection *K8S_Connection) List() []gen_base.FlintDeployment {
	secrets := connection.getSecrets()["items"].([]any)
	deployments := []gen_base.FlintDeployment{}
	visited := []string{}
	for _, secret := range secrets {
		if secret.(map[string]any)["type"] == "v1.flint.io" {
			secret_name := secret.(map[string]any)["metadata"].(map[string]any)["name"].(string)
			deployment_re := regexp.MustCompile("([-a-z0-9]*[a-z0-9]?)-[0-9]+")
			results := deployment_re.FindStringSubmatch(secret_name)
			deployment_name := results[1]
			if slices.Contains(visited, deployment_name) {
				continue
			}
			secret, version := connection.getLatestSecret(deployment_name)
			status := "failed"
			if secret["data"].(map[string]any)["status"].(string) == "c3VjY2Vzcw==" {
				status = "success"
			}
			date, err := time.Parse(time.RFC3339, secret["metadata"].(map[string]any)["creationTimestamp"].(string))
			if err != nil {
				panic(err)
			}
			deployments = append(deployments, gen_base.FlintDeployment{
				Name:     deployment_name,
				Age:      time.Since(date).Round(time.Second).String(),
				Status:   status,
				Revision: version,
			})
			visited = append(visited, deployment_name)
		}

	}
	return deployments
}

func (connection *K8S_Connection) ToFileName(obj map[string]any) string {
	if obj["kind"] == "" {
		return obj["id"].(string)
	}
	return "Kubernetes::" + obj["kind"].(string) + "::" + obj["metadata"].(map[string]any)["namespace"].(string) + "::" + obj["metadata"].(map[string]any)["name"].(string)
}

func (connection *K8S_Connection) Diff(stack map[string]map[string]any, name string) ([]string, []string, [][]map[string]any) {
	secret, version := connection.getLatestSecret(name)

	if version == 0 {
		return make([]string, 1), make([]string, 1), make([][]map[string]any, 1)
	}
	b64_secret := secret["data"].(map[string]any)["obj_map"].(string)
	current, err := base64.StdEncoding.DecodeString(b64_secret)

	if err != nil {
		panic(err)
	}

	added := make([]string, 0)
	removed := make([]string, 0)
	changed := make([][]map[string]any, 0)

	var obj_map map[string]map[string]any
	err = json.Unmarshal(current, &obj_map)
	if err != nil {
		panic(err)
	}

	for newName, newObj := range stack {
		if newObj == nil {
			continue
		}
		found := false
		var foundObjKey string
		for name, _ := range obj_map {
			if name == newName {
				found = true
				foundObjKey = name
				break
			}
		}

		if found {
			bytesNew, err := json.Marshal(newObj)
			if err != nil {
				panic(err)
			}
			bytesOld, err := json.Marshal(obj_map[foundObjKey])
			if err != nil {
				panic(err)
			}
			if !strings.EqualFold(string(bytesNew), string(bytesOld)) {
				objects := make([]map[string]any, 2)
				delete(newObj, "location")
				delete(obj_map[foundObjKey], "location")
				objects[0] = newObj
				objects[1] = obj_map[foundObjKey]
				changed = append(changed, objects)
			}
			delete(obj_map, foundObjKey)
		} else {
			added = append(added, connection.ToFileName(newObj))
		}
	}

	for _, obj := range obj_map {
		removed = append(removed, connection.ToFileName(obj))
	}

	return added, removed, changed
}

func (self *K8S_Connection) deployObjects(dag *dag.DAG, obj_map map[string]map[string]any) {
	// Get all vertices
	vertices := dag.GetVertices()

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
			descendants, err := dag.GetDescendants(id)
			if err != nil {
				mu.Unlock()
				panic(err)
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

				vertex, err := dag.GetVertex(nodeID)
				if err != nil {
					errChan <- err
					return
				}

				self.Apply(obj_map[vertex.(string)], obj_map)

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
				panic(err)
			}
		}
	}
}

func (conn *K8S_Connection) Deploy(dag_ *dag.DAG, to_remove []string, obj_map map[string]map[string]any, name string, stack_metadata map[string]any, max_revisions int) {
	install_number := conn.GetCurrentRevision(name) + 1
	namespace := stack_metadata["namespace"].(string)
	lowest_revision := max(install_number-max_revisions, 0) + 1

	conn.CleanHistory(name, lowest_revision, namespace)

	for _, rem := range to_remove {
		if rem == "" {
			continue
		}
		obj_map[rem] = map[string]any{
			"action": "remove",
			"name":   rem,
		}
		err := dag_.AddVertexByID(rem, rem)
		if err != nil {
			panic(err)
		}
		delete(obj_map, rem)
	}

	conn.deployObjects(dag_, obj_map)

	remove_from_dag := make([]string, 0)
	for name, obj := range obj_map {
		if obj["kind"] == "" {
			delete(obj_map, name)
			remove_from_dag = append(remove_from_dag, name)
		}
	}

	b, err := json.Marshal(dag_)
	if err != nil {
		panic(err)
	}

	var dag_map map[string]any
	err = json.Unmarshal(b, &dag_map)
	if err != nil {
		panic(err)
	}

	newDag := dag.NewDAG()

	for _, value := range dag_map["vs"].([]any) {
		i := value.(map[string]any)["i"].(string)
		v := value.(map[string]any)["v"].(string)
		if slices.Contains(remove_from_dag, i) || slices.Contains(remove_from_dag, v) {
			continue
		}
		newDag.AddVertexByID(v, i)
	}

	for _, value := range dag_map["es"].([]any) {
		d := value.(map[string]any)["d"].(string)
		s := value.(map[string]any)["s"].(string)
		if slices.Contains(remove_from_dag, d) || slices.Contains(remove_from_dag, s) {
			continue
		}
		newDag.AddEdge(s, d)
	}

	secret := Secret{
		Name: name + "-" + strconv.Itoa(install_number),
		Type: "v1.flint.io",
		Data: make([]*SecretData, 3),
	}

	marshalled, _ := json.Marshal(obj_map)

	obj_map_data := SecretData{
		Key:   "obj_map",
		Value: string(marshalled),
	}

	marshalled_dag, _ := json.Marshal(newDag)

	dag_data := SecretData{
		Key:   "dag",
		Value: string(marshalled_dag),
	}

	status := SecretData{
		Key:   "status",
		Value: "success",
	}

	secret.Data[0] = &obj_map_data
	secret.Data[1] = &dag_data
	secret.Data[2] = &status

	secret.Synth(name, namespace, newDag, obj_map)

	conn.Apply(obj_map[secret.GetID()], nil)
}

func (conn *K8S_Connection) Destroy(stack_name string, stack_metadata map[string]any) {
	secret, _ := conn.getLatestSecret(stack_name)
	b64_obj_data := secret["data"].(map[string]any)["obj_map"].(string)
	obj_data_string, err := base64.StdEncoding.DecodeString(b64_obj_data)

	if err != nil {
		panic(err)
	}

	var obj_map map[string]map[string]any
	err = json.Unmarshal(obj_data_string, &obj_map)
	if err != nil {
		panic(err)
	}

	b64_dag_data := secret["data"].(map[string]any)["dag"].(string)
	dag_data_string, err := base64.StdEncoding.DecodeString(b64_dag_data)

	if err != nil {
		panic(err)
	}

	var dag_map map[string]any

	err = json.Unmarshal(dag_data_string, &dag_map)
	if err != nil {
		panic(err)
	}

	dag := dag.NewDAG()

	for _, value := range dag_map["vs"].([]any) {
		i := value.(map[string]any)["i"].(string)
		v := value.(map[string]any)["v"].(string)
		dag.AddVertexByID(i, v)
	}

	for _, value := range dag_map["es"].([]any) {
		d := value.(map[string]any)["d"].(string)
		s := value.(map[string]any)["s"].(string)
		dag.AddEdge(d, s)
	}

	for _, value := range obj_map {
		value["action"] = "remove"
		value["name"] = conn.ToFileName(value)
	}

	conn.deployObjects(dag, obj_map)
}

func (self *K8S_Connection) Rollback(stack_name string, targetRevision int, stack_metadata map[string]any) {
	result := self.getSecrets()
	var target_secret map[string]any

	for _, secret := range result["items"].([]any) {
		if secret.(map[string]any)["type"] == "v1.flint.io" {
			secret_name := secret.(map[string]any)["metadata"].(map[string]any)["name"].(string)
			re := regexp.MustCompile(stack_name + `-[0-9]+`)
			if re.FindString(secret_name) != "" {
				version_re := regexp.MustCompile(`[0-9]+`)
				version, err := strconv.Atoi(version_re.FindString(secret_name))
				if err != nil {
					panic(err)
				}
				if version == targetRevision {
					target_secret = secret.(map[string]any)
					break
				}
			}
		}
	}

	data := target_secret["data"]

	if data == nil {
		panic("revision " + strconv.Itoa(targetRevision) + " doesn't exist or isn't remembered")
	}

	b64_obj_data := data.(map[string]any)["obj_map"].(string)
	obj_data_string, err := base64.StdEncoding.DecodeString(b64_obj_data)

	if err != nil {
		panic(err)
	}

	var obj_map map[string]map[string]any
	err = json.Unmarshal(obj_data_string, &obj_map)
	if err != nil {
		panic(err)
	}

	b64_dag_data := data.(map[string]any)["dag"].(string)
	dag_data_string, err := base64.StdEncoding.DecodeString(b64_dag_data)

	if err != nil {
		panic(err)
	}

	var dag_map map[string]any

	err = json.Unmarshal(dag_data_string, &dag_map)
	if err != nil {
		panic(err)
	}

	dag := dag.NewDAG()

	for _, value := range dag_map["vs"].([]any) {
		i := value.(map[string]any)["i"].(string)
		v := value.(map[string]any)["v"].(string)
		dag.AddVertexByID(i, v)
	}

	for _, value := range dag_map["es"].([]any) {
		d := value.(map[string]any)["d"].(string)
		s := value.(map[string]any)["s"].(string)
		dag.AddEdge(s, d)
	}

	self.deployObjects(dag, obj_map)
}
