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

	"github.com/SuperTapood/Flint/core/generated/common"
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

func (connection *K8SConnection) Apply(applyMetadata map[string]any, resource map[string]any) {
	name := applyMetadata["name"].(string)
	location := applyMetadata["location"].(string)

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
			secretName := secret.(map[string]any)["metadata"].(map[string]any)["name"].(string)
			output[secretName] = secret.(map[string]any)
			b64Secret := secret.(map[string]any)["data"].(map[string]any)["obj_map"].(string)
			current, err := base64.StdEncoding.DecodeString(b64Secret)

			if err != nil {
				panic(err)
			}

			b64Status := secret.(map[string]any)["data"].(map[string]any)["status"].(string)
			status, err := base64.RawStdEncoding.DecodeString(b64Status)
			if err != nil {
				panic(err)
			}

			var currentMap map[string]any
			err = json.Unmarshal(current, &currentMap)
			if err != nil {
				panic(err)
			}

			output[secretName] = map[string]any{
				"map":       currentMap,
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
	for secretName := range secrets {
		deployment_re := regexp.MustCompile("([-a-z0-9]*[a-z0-9]?)-[0-9]+")
		results := deployment_re.FindStringSubmatch(secretName)
		deploymentName := results[1]
		if slices.Contains(visited, deploymentName) {
			continue
		}
		visited = append(visited, deploymentName)
	}

	return visited
}

func (connection *K8SConnection) List() []common.FlintDeployment {
	secrets := connection.GetDeployments()
	deployments := []common.FlintDeployment{}
	visited := []string{}
	for _, deploymentName := range secrets {
		_, status, age, version := connection.GetLatestRevision(deploymentName)
		deployments = append(deployments, common.FlintDeployment{
			Name:     deploymentName,
			Age:      age,
			Status:   status,
			Revision: version,
		})
		visited = append(visited, deploymentName)
	}

	return deployments
}

func (connection *K8SConnection) GetLatestRevision(stackName string) (map[string]any, string, string, int32) {
	result := connection.GetRevisions()
	latestVersion := 0
	var latestSecret map[string]any

	for secretName, secret := range result {
		versionRe := regexp.MustCompile(`[0-9]+`)
		version, err := strconv.Atoi(versionRe.FindString(secretName))

		if err != nil {
			panic(err)
		}
		if version > latestVersion {
			latestSecret = secret
			latestVersion = version
		}
	}

	date, err := time.Parse(time.RFC3339, latestSecret["timestamp"].(string))
	if err != nil {
		panic(err)
	}

	return latestSecret["map"].(map[string]any), latestSecret["status"].(string), time.Since(date).Round(time.Second).String(), int32(latestVersion)
}

func (connection *K8SConnection) PrettyName(resource map[string]any, stackMetadata map[string]any) string {
	return "Kubernetes::" + stackMetadata["namespace"].(string) + "::" + resource["kind"].(string) + "::" + resource["metadata"].(map[string]any)["name"].(string)
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

func (connection *K8SConnection) Delete(deleteMetadata map[string]any) {
	kind := deleteMetadata["kind"].(string)
	namespace := deleteMetadata["namespace"].(string)
	name := deleteMetadata["name"].(string)
	body, resp := connection.MakeRequest("DELETE", locationMap[kind]["before_namespace"]+namespace+locationMap[kind]["after_namespace"]+name, bytes.NewReader(make([]byte, 0)))
	if resp.StatusCode != http.StatusOK {
		fmt.Println(resp)
		fmt.Println(string(body))
		panic("couldn't delete")
	}
}

func (connection *K8SConnection) CleanHistory(stackName string, oldest int, stackMetadata map[string]any) {
	if oldest < 1 {
		return
	}

	result := connection.GetRevisions()

	namespace := stackMetadata["namespace"].(string)

	for secretName, _ := range result {
		re := regexp.MustCompile(stackName + `-[0-9]+`)
		if re.FindString(secretName) != "" {
			versionRe := regexp.MustCompile(`[0-9]+`)
			version, err := strconv.Atoi(versionRe.FindString(secretName))
			if err != nil {
				panic(err)
			}
			if version < oldest {
				connection.Delete(map[string]any{
					"kind":      "Secret",
					"namespace": namespace,
					"name":      secretName,
				})
			}
		}
	}
}

func (connection *K8SConnection) CreateRevision(stackName string, stackMetadata map[string]any, newDag *dag.DAG, marshalled []byte) {
	_, _, _, version := connection.GetLatestRevision(stackName)
	secret := Secret{
		Name: stackName + "-" + strconv.Itoa(int(version+1)),
		Type: "v1.flint.io",
		Data: make([]*SecretData, 3),
	}

	objMapData := SecretData{
		Key:   "obj_map",
		Value: string(marshalled),
	}

	marshalledDag, _ := json.Marshal(newDag)

	dagData := SecretData{
		Key:   "dag",
		Value: string(marshalledDag),
	}

	status := SecretData{
		Key:   "status",
		Value: "success",
	}

	secret.Data[0] = &objMapData
	secret.Data[1] = &dagData
	secret.Data[2] = &status

	secret.Apply(stackMetadata, nil, connection)
}

func (connection *K8SConnection) PrintOutputs() {
	for i := range len(outputBufferMap) {
		fmt.Println(outputBufferMap[int32(i)].String())
	}
}
