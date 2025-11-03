package base

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"slices"
	"strconv"
	"time"
)

type K8SConnection struct {
	Api   string
	Token string
}

func (connection *K8SConnection) getSecrets() map[string]any {
	var body, _ = connection.makeRequest("GET", "/api/v1/secrets", bytes.NewReader(make([]byte, 1)))
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		panic(err)
	}

	return result
}

func (connection *K8SConnection) getLatestSecret(stack_name string) (map[string]any, int) {
	result := connection.getSecrets()
	latest_version := -1
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

	return latest_secret, latest_version
}

func (connection *K8SConnection) GetCurrentRevision(stack_name string) int {
	_, latest_version := connection.getLatestSecret(stack_name)
	return latest_version + 1
}

func (connection *K8SConnection) makeRequest(method string, location string, reader io.Reader) ([]byte, *http.Response) {
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

func mergeMaps(dst, src map[string]interface{}) {
	for k, v := range src {
		if dv, ok := dst[k]; ok {
			if dm, ok := dv.(map[string]interface{}); ok {
				if sm, ok := v.(map[string]interface{}); ok {
					mergeMaps(dm, sm)
					continue
				}
			}
		}
		dst[k] = v
	}
}

func (connection *K8SConnection) Deploy(obj map[string]any) {
	location := obj["location"].(string)
	delete(obj, "location")

	data, _ := json.Marshal(obj)

	var resp *http.Response
	var body []byte

	body, resp = connection.makeRequest("POST", location, bytes.NewReader(data))

	if resp.StatusCode == http.StatusConflict {
		body, resp = connection.makeRequest("PUT", location+"/"+obj["metadata"].(map[string]any)["name"].(string), bytes.NewReader(data))
		if resp.StatusCode == http.StatusUnprocessableEntity {
			body, resp = connection.makeRequest("DELETE", location+"/"+obj["metadata"].(map[string]any)["name"].(string), bytes.NewReader(make([]byte, 0)))
			if resp.StatusCode == http.StatusOK {
				time.Sleep(2 * time.Second)
				body, resp = connection.makeRequest("POST", location, bytes.NewReader(data))
			}
		}
	}

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		fmt.Println(location)
		fmt.Println(resp)
		fmt.Println(string(body))
		fmt.Println(obj)
	}
}

func (connection *K8SConnection) List() []Deployment {
	secrets := connection.getSecrets()["items"].([]any)
	deployments := []Deployment{}
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
			deployments = append(deployments, Deployment{
				Name:     deployment_name,
				Age:      time.Since(date).Truncate(time.Second),
				Status:   status,
				Revision: version,
			})
			visited = append(visited, deployment_name)
		}

	}
	return deployments
}
