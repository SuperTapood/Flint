package k8s

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"slices"
	"strconv"
	"time"

	"github.com/SuperTapood/Flint/core/generated/common"
	"github.com/SuperTapood/Flint/core/util"

	"github.com/heimdalr/dag"
)

type KubeError struct{}

func (kubeError *KubeError) Error() string {
	return "kube failed :("
}

func (connection *K8SConnection) GetClient() *util.HttpClient {
	return util.NewHttpClient(map[string]string{
		"Authorization": "Bearer " + connection.Token,
		"Content-Type":  "application/json",
	}, connection.Api)
}

func (connection *K8SConnection) Apply(applyMetadata map[string]any, resource map[string]any) error {
	name := applyMetadata["name"].(string)
	location := applyMetadata["location"].(string)

	data, err := json.Marshal(resource)
	if err != nil {
		fmt.Println("couldn't unmarshal the resource to deploy")
		fmt.Println(err)
		os.Exit(2)
	}

	client := connection.GetClient()

	response, err := client.Post(location, bytes.NewReader(data), []int{http.StatusOK, http.StatusCreated, http.StatusConflict}, false)

	if response.StatusCode == http.StatusConflict {
		response, err = client.Put(location+name, bytes.NewReader(data), []int{http.StatusOK, http.StatusCreated, http.StatusUnprocessableEntity}, false)
		if response.StatusCode == http.StatusUnprocessableEntity {
			response, err = client.Delete(location+name, nil, false)
			if response.StatusCode == http.StatusOK {
				time.Sleep(2 * time.Second)
				response, err = client.Post(location, bytes.NewReader(data), nil, false)
			}
		}
	}

	if err != nil {
		return err
	}

	start := time.Now()

	var status map[string]any

	for {
		response, err = client.Get(location+name, []int{http.StatusOK}, false) // make sure resource exists
		var ok bool
		status, ok = response.Body["status"].(map[string]any)

		if !ok {
			return nil
		}
		if len(status) != 0 {
			available, ok := status["availableReplicas"]
			replicas, rok := status["replicas"]

			if !rok {
				break
			}

			if ok && available == replicas {
				break
			}
		}

		if time.Since(start) > 5*time.Second {
			fmt.Println(response)
			return &KubeError{}
		}
		time.Sleep(100 * time.Millisecond)
	}
	return err
}

func (connection *K8SConnection) GetRevisions() map[string]map[string]any {
	client := connection.GetClient()
	var response, _ = client.Get("/api/v1/secrets", nil, true)
	var result = response.Body
	// if err := json.Unmarshal(response.Body, &result); err != nil {
	// 	fmt.Println("couldn't unmarshal response from server")
	// 	fmt.Println(err)
	// 	os.Exit(2)
	// }
	output := make(map[string]map[string]any, 0)

	for _, secret := range result["items"].([]any) {
		if secret.(map[string]any)["type"] == "v1.flint.io" {
			secretName := secret.(map[string]any)["metadata"].(map[string]any)["name"].(string)
			output[secretName] = secret.(map[string]any)
			b64Secret := secret.(map[string]any)["data"].(map[string]any)["obj_map"].(string)
			current, err := base64.StdEncoding.DecodeString(b64Secret)

			if err != nil {
				fmt.Println("failed to decode current secret")
				fmt.Println(err)
				os.Exit(-1)
			}

			b64Status := secret.(map[string]any)["data"].(map[string]any)["status"].(string)
			status, err := base64.StdEncoding.DecodeString(b64Status)
			if err != nil {
				fmt.Println("failed to decode current secret status")
				fmt.Println(err)
				os.Exit(-1)
			}

			var currentMap map[string]any
			err = json.Unmarshal(current, &currentMap)
			if err != nil {
				fmt.Println("couldn't unmarshal the current map to get revisions")
				fmt.Println(err)
				os.Exit(2)
			}

			output[secretName] = map[string]any{
				"map":       currentMap,
				"status":    string(status),
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
		version, _ := strconv.Atoi(versionRe.FindString(secretName))

		if version > latestVersion {
			latestSecret = secret
			latestVersion = version
		}
	}

	if latestSecret == nil {
		return nil, "", "", 0
	}

	date, err := time.Parse(time.RFC3339, latestSecret["timestamp"].(string))
	if err != nil {
		fmt.Println("could not parse the timestamp")
		fmt.Println(err)
		os.Exit(2)
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
	client := connection.GetClient()
	client.Delete(locationMap[kind]["before_namespace"]+namespace+locationMap[kind]["after_namespace"]+name, nil, true)
}

func (connection *K8SConnection) CleanHistory(stackName string, oldest int, stackMetadata map[string]any) {
	if oldest < 1 {
		return
	}

	result := connection.GetRevisions()

	namespace := stackMetadata["namespace"].(string)

	for secretName := range result {
		re := regexp.MustCompile(stackName + `-[0-9]+`)
		if re.FindString(secretName) != "" {
			versionRe := regexp.MustCompile(`[0-9]+`)
			version, _ := strconv.Atoi(versionRe.FindString(secretName))
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
