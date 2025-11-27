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
	"time"

	common "github.com/SuperTapood/Flint/core/generated/common"
	"github.com/heimdalr/dag"
)

func (connection *K8SConnection) MakeRequest(method string, location string, reader io.Reader) ([]byte, *http.Response) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	req, err := http.NewRequest(method, connection.Api+location, reader)

	if err != nil {
		panic(err)
	}

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

func (connection *K8SConnection) Apply(apply_metadata map[string]any, resource map[string]any) {
	name := apply_metadata["name"].(string)
	location := apply_metadata["location"].(string)

	data, err := json.Marshal(resource)
	if err != nil {
		panic(err)
	}

	var resp *http.Response
	var body []byte

	body, resp = connection.MakeRequest("POST", location, bytes.NewReader(data))

	if resp.StatusCode == http.StatusConflict {
		body, resp = connection.MakeRequest("PUT", location+name, bytes.NewReader(data))
		if resp.StatusCode == http.StatusUnprocessableEntity {
			body, resp = connection.MakeRequest("DELETE", location+name, bytes.NewReader(make([]byte, 0)))
			if resp.StatusCode == http.StatusOK {
				time.Sleep(2 * time.Second)
				body, resp = connection.MakeRequest("POST", location, bytes.NewReader(data))
			}
		}
	}

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		fmt.Println("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
		fmt.Println(resp)
		fmt.Println(string(body))
		fmt.Println(location)
		panic("FUCK")
	}

}

func (connection *K8SConnection) GetRevisions() map[string]map[string]any {
	var body, _ = connection.MakeRequest("GET", "/api/v1/secrets", bytes.NewReader(make([]byte, 1)))
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		panic(err)
	}
	output := make(map[string]map[string]any, 0)

	for _, secret := range result["items"].([]any) {
		if secret.(map[string]any)["type"] == "v1.flint.io" {
			secret_name := secret.(map[string]any)["metadata"].(map[string]any)["name"].(string)
			output[secret_name] = secret.(map[string]any)
			b64_secret := secret.(map[string]any)["data"].(map[string]any)["obj_map"].(string)
			current, err := base64.StdEncoding.DecodeString(b64_secret)

			if err != nil {
				panic(err)
			}

			status := "failed"
			if secret.(map[string]any)["data"].(map[string]any)["status"].(string) == "c3VjY2Vzcw==" {
				status = "success"
			}

			var current_map map[string]any
			err = json.Unmarshal(current, &current_map)
			if err != nil {
				panic(err)
			}

			output[secret_name] = map[string]any{
				"map":       current_map,
				"status":    status,
				"timestamp": secret.(map[string]any)["metadata"].(map[string]any)["creationTimestamp"].(string),
			}
		}
	}

	return output
}

func (connection *K8SConnection) GetDeployments() []string {
	visited := []string{}
	secrets := connection.GetRevisions()
	for secret_name, _ := range secrets {
		deployment_re := regexp.MustCompile("([-a-z0-9]*[a-z0-9]?)-[0-9]+")
		results := deployment_re.FindStringSubmatch(secret_name)
		deployment_name := results[1]
		if slices.Contains(visited, deployment_name) {
			continue
		}
		visited = append(visited, deployment_name)
	}

	return visited
}

func (connection *K8SConnection) List() []common.FlintDeployment {
	secrets := connection.GetDeployments()
	deployments := []common.FlintDeployment{}
	visited := []string{}
	for _, deployment_name := range secrets {
		_, status, age, version := connection.GetLatestRevision(deployment_name)
		deployments = append(deployments, common.FlintDeployment{
			Name:     deployment_name,
			Age:      age,
			Status:   status,
			Revision: version,
		})
		visited = append(visited, deployment_name)
	}

	return deployments
}

func (connection *K8SConnection) GetLatestRevision(stack_name string) (map[string]any, string, string, int32) {
	result := connection.GetRevisions()
	latest_version := 0
	var latest_secret map[string]any

	for secret_name, secret := range result {
		version_re := regexp.MustCompile(`[0-9]+`)
		version, err := strconv.Atoi(version_re.FindString(secret_name))

		if err != nil {
			panic(err)
		}
		if version > latest_version {
			latest_secret = secret
			latest_version = version
		}
	}

	date, err := time.Parse(time.RFC3339, latest_secret["timestamp"].(string))
	if err != nil {
		panic(err)
	}

	return latest_secret["map"].(map[string]any), latest_secret["status"].(string), time.Since(date).Round(time.Second).String(), int32(latest_version)
}

func (connection *K8SConnection) PrettyName(resource map[string]any, stack_metadata map[string]any) string {
	return "Kubernetes::" + stack_metadata["namespace"].(string) + "::" + resource["kind"].(string) + "::" + resource["metadata"].(map[string]any)["name"].(string)
}

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

func (connection *K8SConnection) Delete(delete_metadata map[string]any) {
	kind := delete_metadata["kind"].(string)
	namespace := delete_metadata["namespace"].(string)
	name := delete_metadata["name"].(string)
	body, resp := connection.MakeRequest("DELETE", locationMap[kind]["before_namespace"]+namespace+locationMap[kind]["after_namespace"]+name, bytes.NewReader(make([]byte, 0)))
	if resp.StatusCode != http.StatusOK {
		fmt.Println(resp)
		fmt.Println(string(body))
		panic("couldn't delete")
	}
}

func (connection *K8SConnection) CleanHistory(stack_name string, oldest int, stack_metadata map[string]any) {
	if oldest < 1 {
		return
	}

	result := connection.GetRevisions()

	namespace := stack_metadata["namespace"].(string)

	for secret_name, _ := range result {
		re := regexp.MustCompile(stack_name + `-[0-9]+`)
		if re.FindString(secret_name) != "" {
			version_re := regexp.MustCompile(`[0-9]+`)
			version, err := strconv.Atoi(version_re.FindString(secret_name))
			if err != nil {
				panic(err)
			}
			if version < oldest {
				connection.Delete(map[string]any{
					"kind":      "Secret",
					"namespace": namespace,
					"name":      secret_name,
				})
			}
		}
	}
}

func (connection *K8SConnection) CreateRevision(stack_name string, stack_metadata map[string]any, newDag *dag.DAG, marshalled []byte) {
	_, _, _, version := connection.GetLatestRevision(stack_name)
	secret := Secret{
		Name: stack_name + "-" + strconv.Itoa(int(version+1)),
		Type: "v1.flint.io",
		Data: make([]*SecretData, 3),
	}

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

	secret.Apply(stack_metadata, nil, connection)
}

const (
	APIS_APP_V1 = "/apis/apps/v1/namespaces/"
	API_V1      = "/api/v1/namespaces/"
)

func (connection *K8SConnection) PrintOutputs() {
	for i := range len(outputBufferMap) {
		fmt.Println(outputBufferMap[int32(i)].String())
	}
}
