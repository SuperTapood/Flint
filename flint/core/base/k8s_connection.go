package base

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
	// Add this for pointer.Bool
)

type K8SConnection struct {
	Api   string
	Token string
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

func (connection *K8SConnection) Deploy(obj map[string]any, name string) {
	location := obj["location"].(string)
	delete(obj, "location")

	data, _ := json.Marshal(obj)

	var resp *http.Response

	_, resp = connection.makeRequest("POST", location, bytes.NewReader(data))

	if resp.StatusCode == http.StatusConflict {
		_, resp = connection.makeRequest("PUT", location+"/"+obj["metadata"].(map[string]any)["name"].(string), bytes.NewReader(data))
		if resp.StatusCode == http.StatusUnprocessableEntity {
			_, resp = connection.makeRequest("DELETE", location+"/"+obj["metadata"].(map[string]any)["name"].(string), bytes.NewReader(make([]byte, 0)))
			if resp.StatusCode == http.StatusOK {
				time.Sleep(2 * time.Second)
				_, resp = connection.makeRequest("POST", location, bytes.NewReader(data))
			}
		}
	}

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		fmt.Println(resp)
	}
}

func (connection *K8SConnection) List() {
	var body, _ = connection.makeRequest("GET", "/api/v1/secrets", bytes.NewReader(make([]byte, 1)))
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		panic(err)
	}

	secrets := make(map[string]any, 0)

	for i, secret := range result["items"].([]any) {
		fmt.Println(i)
		fmt.Println(secret.(map[string]any)["metadata"].(map[string]any)["name"])
		fmt.Println(secret.(map[string]any)["type"])
		if secret.(map[string]any)["type"] == "Opaque" {
			secrets[secret.(map[string]any)["metadata"].(map[string]any)["name"].(string)] = secret
		}
	}
}
